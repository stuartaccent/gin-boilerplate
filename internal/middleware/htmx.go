package middleware

import (
	"gin.go.dev/internal/htmx"
	"github.com/gin-gonic/gin"
)

// HTMX middleware func to set the HTMX helper.
func HTMX() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("htmx", htmx.New(c.Request, c.Writer))
		c.Next()
	}
}
