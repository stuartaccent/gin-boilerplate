package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContentTypes blocks requests with invalid Content-Type headers.
func ContentTypes(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isValidType := false
		for _, v := range allowedTypes {
			if c.GetHeader("Content-Type") == v {
				isValidType = true
				break
			}
		}
		if !isValidType {
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
				"error": "Invalid content type",
			})
			return
		}

		c.Next()
	}
}
