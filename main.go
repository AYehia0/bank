package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/AYehia0/go-bk-mst/api"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	grpcapi "github.com/AYehia0/go-bk-mst/grpc_api"
	"github.com/AYehia0/go-bk-mst/pb"
	"github.com/AYehia0/go-bk-mst/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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
	runGrpcServer(config, store)

}

func runGrpcServer(config utils.Config, store db.Store) {

	srv, err := grpcapi.NewServer(config, store)

	if err != nil {
		log.Fatalf("Couldn't create the gRPC server : %s", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, srv)

	// DANDER: enabling reflection allows client to explore all RPCs on the server
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddr)

	if err != nil {
		log.Fatalf("Couldn't create the listener : %s", err)
	}
	log.Printf("Server [gRPC] is running at : %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Couldn't start gRPC server : %s", err)
	}

}
func runGinServer(config utils.Config, store db.Store) {

	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatalf("Failed to start the server : %v", err)
	}

	server.StartServer(config.GinServerAddr)

}
