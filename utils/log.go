package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LogRequestBodyMiddleware(c *gin.Context) {
	// Read the request body
	if c.Request.Body == nil {
		fmt.Printf("Request body is empty!")
		c.Next()
		return
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		c.Abort()
		return
	}

	// Restore the request body for downstream handlers
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Log the request body
	fmt.Printf("[%s] Request Body: %s\n", time.Now().Format(time.RFC3339), body)

	// Continue to the next middleware/handler
	c.Next()
}
