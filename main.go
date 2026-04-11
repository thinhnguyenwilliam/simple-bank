// simple-bank\main.go
// sudo lsof -i :9091
// kill -9 492819
// http://192.168.1.8:8084/swagger-ui/
// http://192.168.1.8:8084/swagger/simple_bank.swagger.json

package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/cors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/thinhcompany/simple-bank/api"
	db "github.com/thinhcompany/simple-bank/db/sqlc"
	"github.com/thinhcompany/simple-bank/gapi"
	"github.com/thinhcompany/simple-bank/worker"

	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"
)

func main() {
	util.SetupLogger("api")
	log.Println("🚀 Starting simple-bank API")

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
	// go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runGatewayServer(
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {
	_, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("cannot create gRPC server:", err)
	}

	// ===== gRPC-Gateway mux =====
	jsonOption := runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	)
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = pb.RegisterSimpleBankServiceHandlerFromEndpoint(
		ctx,
		grpcMux,
		config.GrpcServerAddress, // 👈 connect tới gRPC server
		opts,
	)

	// ===== HTTP mux (QUAN TRỌNG) =====
	mux := http.NewServeMux()

	// =========================
	// 🧾 Serve Swagger JSON (STATIC FILE)
	// =========================
	swaggerFS := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/",
		http.StripPrefix("/swagger/", swaggerFS),
	)

	// =========================
	// 🎨 Swagger UI
	// =========================
	swaggerUI := httpSwagger.Handler(
		httpSwagger.URL("/swagger/pb/v1/service_simple_bank.swagger.json"),
	)
	mux.Handle("/swagger-ui/", swaggerUI)

	// =========================
	// gRPC Gateway (CUỐI CÙNG)
	// =========================
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("🌐 HTTP Gateway listening at %s", listener.Addr().String())

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://192.168.1.8:5175",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		ExposedHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 giờ
	})

	handler := c.Handler(mux)

	if err := http.Serve(listener, handler); err != nil {
		log.Fatal("cannot start HTTP gateway:", err)
	}
}

// func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
// 	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
// 	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

// 	log.Println("starting task processor...")
// 	if err := taskProcessor.Start(); err != nil {
// 		log.Fatal("failed to start task processor:", err)
// 	}
// }

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("cannot create grpc server:", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// server.AuthInterceptor(), // ✅ auth trước
			gapi.GrpcLogger, // ✅ logger sau
		),
	)
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
func runMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(
		migrationURL,
		dbSource,
	)
	if err != nil {
		log.Fatal("cannot create migration instance:", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migration:", err)
	}

	log.Println("✅ database migrated successfully")
}
