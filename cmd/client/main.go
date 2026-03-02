// simple-bank/cmd/client/main.go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// 1️⃣ Load CA cert (để verify SERVER)
	pemCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot read CA cert: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCA) {
		return nil, fmt.Errorf("failed to add CA cert")
	}

	// 2️⃣ Load CLIENT cert & key (để server verify)
	clientCert, err := tls.LoadX509KeyPair(
		"cert/client-cert.pem",
		"cert/client-key.pem",
	)
	if err != nil {
		return nil, fmt.Errorf("cannot load client cert/key: %w", err)
	}

	// 3️⃣ TLS config cho client (mTLS)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert}, // 👈 gửi client cert
		RootCAs:      certPool,                      // 👈 tin CA ký server
		ServerName:   "localhost",                   // 👈 khớp SAN
	}

	return credentials.NewTLS(tlsConfig), nil
}

const serverAddress = "localhost:50051"

func main() {
	log.Println("🧑‍💻 client starting...")
	tlsCreds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("cannot load TLS credentials: %v", err)
	}

	// 1️⃣ Connect to gRPC server
	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(tlsCreds),
	)
	if err != nil {
		log.Fatalf("cannot connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewLaptopServiceClient(conn)

	// 2️⃣ Prepare request
	req := &pb.CreateLaptopRequest{
		Laptop: util.NewLaptop(),
	}

	// 3️⃣ Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 4️⃣ Call RPC
	res, err := client.CreateLaptop(ctx, req)
	if err != nil {
		log.Fatalf("CreateLaptop failed: %v", err)
	}

	log.Printf("✅ laptop created with ID: %s", res.Id)
}
