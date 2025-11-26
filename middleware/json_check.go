package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			if c.GetHeader("Content-Type") != "application/json" {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Content-Type must be application/json"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
