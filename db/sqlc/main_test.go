// connect to the database
// you will need the actual database driver installed : github.com/lib/pq

package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/AYehia0/go-bk-mst/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.ConfigStore("../..", "config", "env")
	testDb, err = sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatalf("Failed to connect to the database : %v", err)
	}

	testQueries = New(testDb)

	// terminate the connection if success
	os.Exit(m.Run())

}
