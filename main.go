// simple-bank\main.go

package main

import (
	"database/sql"
	"log"

	"simple-bank/api"
	db "simple-bank/db/sqlc"
	"simple-bank/util"

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
	log.Println("DB is:", config.DBSource)
	log.Println("Server is:", config.ServerAddress)

	// Verify DB connection
	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	// Create store
	store := db.NewStore(conn)

	// Create server
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// Start server
	if err := server.Run(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
