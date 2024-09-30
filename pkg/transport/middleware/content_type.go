package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ContentTypes blocks requests with invalid Content-Type headers.
func ContentTypes(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		valid := false
		for _, v := range allowedTypes {
			if c.GetHeader("Content-Type") == v {
				valid = true
				break
			}
		}
		if !valid {
			c.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}

		c.Next()
	}
}
