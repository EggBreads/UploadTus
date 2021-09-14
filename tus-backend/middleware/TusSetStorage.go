package middleware

import (
	"github.com/gin-gonic/gin"
)

func ValidateCustomMiddleware(err error) func(*gin.Context) {
	return func(c *gin.Context) {
		if err != nil {
			c.Abort()
		}
		c.Next()
	}
}
