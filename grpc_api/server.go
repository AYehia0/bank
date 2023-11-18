package grpcapi

import (
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/pb"
	"github.com/AYehia0/go-bk-mst/token"
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	pb.UnimplementedSimpleBankServer // multiple rpcs
	store                            db.Store
	tokenCreator                     token.TokenCreator
	config                           utils.Config
	router                           *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	// create a token factory
	creator, err := token.NewPasteoToken(config.TokenKey)

	if err != nil {
		return nil, err
	}
	server := &Server{
		store:        store,
		tokenCreator: creator,
		config:       config,
	}

	return server, nil
}
