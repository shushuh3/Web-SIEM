package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuth(user, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, p, hasAuth := c.Request.BasicAuth()

		if !hasAuth || u != user || p != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
			})
			return
		}

		c.Next()
	}
}
