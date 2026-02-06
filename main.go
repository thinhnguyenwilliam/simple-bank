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
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	log.Println("DB is:", config.DBSource)
	log.Println("Server is:", config.ServerAddress)

	// always ping to verify connection
	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	if err := server.Run(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
