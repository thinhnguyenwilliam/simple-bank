// simple-bank/cmd/client/main.go
package main

import (
	"context"
	"log"
	"time"

	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddress = "localhost:50051"

func main() {
	log.Println("🧑‍💻 client starting...")

	// 1️⃣ Connect to gRPC server
	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
