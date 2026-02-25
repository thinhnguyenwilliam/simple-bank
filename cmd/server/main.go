// simple-bank/cmd/server/main.go
package main

import (
	"log"
	"net"
	"os"

	"github.com/thinhcompany/simple-bank/interceptor"
	"github.com/thinhcompany/simple-bank/model"
	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const grpcAddress = ":50051"

func main() {
	log.Println("🚀 server starting...")

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

	// 4️⃣ Init store
	laptopStore := service.NewDBLaptopStore(db)

	// 5️⃣ Init gRPC server
	laptopServer := service.NewLaptopServer(laptopStore)

	grpcServer := grpc.NewServer(
		interceptor.Unary(),
		interceptor.Stream(),
	)
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
