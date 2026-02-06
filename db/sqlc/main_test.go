package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"simple-bank/util"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const dbDriver = "postgres"

func TestMain(m *testing.M) {
	// from db/sqlc → project root
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(dbDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
