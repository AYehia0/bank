package api

import (
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	// methods
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)

	server.router = router

	return server
}

func (server *Server) StartServer(addressUrl string) error {
	return server.router.Run(addressUrl)
}
