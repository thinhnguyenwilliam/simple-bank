// simple-bank/cmd/server/main.go
// Có TLS / mTLS:  go run cmd/server/main.go -port 50051 -tls
// Không TLS (dev nhanh): go run cmd/server/main.go
// go run cmd/server/main.go -port 50051 -tls
// go run cmd/server/main.go -port 50052 -tls
// go run cmd/server/main.go -port 50053
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/thinhcompany/simple-bank/interceptor"
	"github.com/thinhcompany/simple-bank/model"
	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// 1️⃣ Load server cert & key
	serverCert, err := tls.LoadX509KeyPair(
		"cert/server-cert.pem",
		"cert/server-key.pem",
	)
	if err != nil {
		return nil, fmt.Errorf("cannot load server cert/key: %w", err)
	}

	// 2️⃣ Load CA cert (để verify CLIENT cert)
	pemCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot read CA cert: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCA) {
		return nil, fmt.Errorf("failed to add CA cert to pool")
	}

	// 3️⃣ TLS config (mTLS)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool, // dùng CA để verify client
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}

func main() {
	// 0️⃣ Flags
	port := flag.String("port", "50051", "gRPC server port")
	enableTLS := flag.Bool("tls", false, "Enable SSL/TLS")
	flag.Parse()

	grpcAddress := ":" + *port
	log.Printf("🚀 server starting... TLS=%v", *enableTLS)

	// 1️⃣ Load env
	if err := godotenv.Load("app.env"); err != nil {
		log.Println("⚠️  cannot load app.env, using system env")
	}

	dsn := os.Getenv("DB_SOURCE")
	if dsn == "" {
		log.Fatal("DB_SOURCE is not set")
	}

	// 2️⃣ Connect DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect db: %v", err)
	}

	// 3️⃣ Auto migrate
	if err := db.AutoMigrate(&model.LaptopModel{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// 4️⃣ Init store + service
	laptopStore := service.NewDBLaptopStore(db)
	laptopServer := service.NewLaptopServer(laptopStore)

	// 5️⃣ Init gRPC server
	var grpcServer *grpc.Server

	if *enableTLS {
		tlsCreds, err := loadTLSCredentials()
		if err != nil {
			log.Fatalf("cannot load TLS credentials: %v", err)
		}

		grpcServer = grpc.NewServer(
			grpc.Creds(tlsCreds),
			interceptor.Unary(),
			interceptor.Stream(),
		)

		log.Println("🔐 gRPC running with TLS (mTLS enabled)")
	} else {
		grpcServer = grpc.NewServer(
			interceptor.Unary(),
			interceptor.Stream(),
		)

		log.Println("⚠️  gRPC running WITHOUT TLS")
	}

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	// 6️⃣ Listen
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}

	log.Printf("📡 gRPC server listening on %s", grpcAddress)
	reflection.Register(grpcServer)

	// 7️⃣ Serve
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("cannot serve gRPC: %v", err)
	}
}
