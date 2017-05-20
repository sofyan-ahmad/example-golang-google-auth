package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthorizeRequest is used to authorize a request for a certain end-point group.
func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-id")

		if v == nil {
			c.HTML(http.StatusUnauthorized, "error", gin.H{"message": "Unauthorized, please login"})
			c.Abort()
		}

		c.Next()
	}
}
