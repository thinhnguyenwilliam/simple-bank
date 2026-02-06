// simple-bank\db\sqlc\main_test.go

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

// Start using ci on github action
// package db

// import (
// 	"database/sql"
// 	"log"
// 	"os"
// 	"testing"

// 	_ "github.com/lib/pq"
// )

// var testQueries *Queries
// var testDB *sql.DB

// func TestMain(m *testing.M) {
// 	dbSource := os.Getenv("DB_SOURCE")
// 	if dbSource == "" {
// 		log.Fatal("DB_SOURCE is not set")
// 	}

// 	var err error
// 	testDB, err = sql.Open("postgres", dbSource)
// 	if err != nil {
// 		log.Fatal("cannot connect to db:", err)
// 	}

// 	testQueries = New(testDB)

// 	code := m.Run()

// 	testDB.Close()
// 	os.Exit(code)
// }
