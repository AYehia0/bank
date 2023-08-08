package helpers

import "github.com/gin-gonic/gin"

func ErrorResp(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
