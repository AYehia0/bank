package api

import (
	"os"
	"testing"
	"time"

	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenKey:            utils.RandomString(32),
		TokenExpireDuration: time.Minute,
	}
	server, err := NewServer(config, store)

	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	// ignore gin debug mode
	gin.SetMode(gin.TestMode)

	// terminate the connection if success
	os.Exit(m.Run())
}
