// simple-bank\main.go
package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/hibiken/asynq"
	"github.com/thinhcompany/simple-bank/api"
	db "github.com/thinhcompany/simple-bank/db/sqlc"
	"github.com/thinhcompany/simple-bank/gapi"
	"github.com/thinhcompany/simple-bank/worker"

	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	// Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Open database connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	// Verify DB connection
	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	// Create store
	store := db.NewStore(conn)

	//
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	//
	runGrpcServer(config, store, taskDistributor)
	// runGinServer(config, store)
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)

	log.Println("starting task processor...")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal("failed to start task processor:", err)
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("cannot create grpc server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServiceServer(grpcServer, server)

	// ✅ Enable reflection
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("gRPC server running at %s", config.GrpcServerAddress)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start grpc server:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	if err := server.Run(config.HttpServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
