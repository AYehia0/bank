package main

import (
	"database/sql"
	"log"

	"github.com/AYehia0/go-bk-mst/api"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"

	// important for database init
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"

	// TODO: use env_vars
	dbSourceUrl   = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// connect to the database
	var err error
	conn, err := sql.Open(dbDriver, dbSourceUrl)

	if err != nil {
		log.Fatalf("Failed to connect to the database : %v", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err != nil {
		log.Fatalf("Failed to start the server : %v", err)
	}

	server.StartServer(serverAddress)

}
