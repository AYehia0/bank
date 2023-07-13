// connect to the database
// you will need the actual database driver installed : github.com/lib/pq

package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

const (
	dbDriver = "postgres"

	// TODO: use env_vars
	dbSourceUrl = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbDriver, dbSourceUrl)

	if err != nil {
		log.Fatalf("Failed to connect to the database : %v", err)
	}

	testQueries = New(testDb)

	// terminate the connection if success
	os.Exit(m.Run())

}
