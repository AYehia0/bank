package api

import (
	"fmt"

	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/token"
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store        db.Store
	tokenCreator token.TokenCreator
	config       utils.Config
	router       *gin.Engine
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

	// registering validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		fmt.Println("Registering Validation Functions")
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupServer()

	return server, nil
}

func (server *Server) StartServer(addressUrl string) error {
	return server.router.Run(addressUrl)
}

// define the routes
func (server *Server) setupServer() {
	// methods
	router := gin.Default()

	// middlewares
	router.Use(utils.LogRequestBodyMiddleware)

	// requires no login
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_token", server.requestNewAccessToken)

	// create a group for them
	authRequired := router.Group("/").Use(authMiddleware(server.tokenCreator))

	authRequired.POST("/accounts", server.createAccount)
	authRequired.GET("/accounts/:id", server.getAccount)
	authRequired.GET("/accounts", server.getAccounts)

	authRequired.POST("/transfers", server.createTransfer)

	server.router = router
}
