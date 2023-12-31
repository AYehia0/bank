package main

import (
	"database/sql"
	"log"

	"github.com/AYehia0/go-bk-mst/api"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/utils"

	// important for database init
	_ "github.com/lib/pq"
)

func main() {
	// connect to the database
	config, err := utils.ConfigStore(".", "config", "env")

	if err != nil {
		log.Fatalf("Couldn't load configs, error: %s", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatalf("Failed to connect to the database : %v", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatalf("Failed to start the server : %v", err)
	}

	server.StartServer(config.ServerAddr)

}
