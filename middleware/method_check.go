package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MethodCheck() gin.HandlerFunc{
	return func(c*gin.Context){
		if c.Request.Method!=http.MethodPost && c.Request.Method!=http.MethodPut{
			c.JSON(http.StatusMethodNotAllowed,gin.H{"error":"method not allowed"})
			c.Abort()
			return
		}
		c.Next()
	}
}