package api

import (
	"fmt"

	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	// registering validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		fmt.Println("Registering Validation Functions")
		v.RegisterValidation("currency", validCurrency)
	}

	// middlewares
	router.Use(utils.LogRequestBodyMiddleware)

	// methods
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)

	router.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)

	server.router = router

	return server
}

func (server *Server) StartServer(addressUrl string) error {
	return server.router.Run(addressUrl)
}
