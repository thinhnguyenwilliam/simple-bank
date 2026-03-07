// simple-bank/cmd/worker/main.go
// go run cmd/worker/main.go
package main

import (
	"database/sql"
	"log"

	"github.com/hibiken/asynq"
	db "github.com/thinhcompany/simple-bank/db/sqlc"
	"github.com/thinhcompany/simple-bank/mail"
	"github.com/thinhcompany/simple-bank/util"
	"github.com/thinhcompany/simple-bank/worker"

	_ "github.com/lib/pq"
)

func main() {
	util.SetupLogger("worker")
	log.Println("🚀 starting worker service...")

	// Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// DB
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	store := db.NewStore(conn)

	// Mailer
	mailer := mail.NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)

	// Asynq
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	processor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Println("📨 worker is running and waiting for tasks...")
	if err := processor.Start(); err != nil {
		log.Fatal("worker stopped:", err)
	}
}
