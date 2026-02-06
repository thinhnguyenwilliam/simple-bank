// simple-bank\main.go

package main

import (
	"database/sql"
	"log"

	"simple-bank/api"
	db "simple-bank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://admin:admin123@192.168.1.8:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}

	// always ping to verify connection
	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	if err := server.Run(serverAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
