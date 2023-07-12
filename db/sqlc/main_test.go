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

const (
	dbDriver = "postgres"

	// TODO: use env_vars
	dbSourceUrl = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSourceUrl)

	if err != nil {
		log.Fatalf("Failed to connect to the database : %v", err)
	}

	testQueries = New(conn)

	// terminate the connection if success
	os.Exit(m.Run())

}
