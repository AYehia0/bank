package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// ignore gin debug mode
	gin.SetMode(gin.TestMode)

	// terminate the connection if success
	os.Exit(m.Run())
}
