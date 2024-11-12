package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AllowContentType enforces a whitelist of request Content-Types otherwise responds
// with a 415 Unsupported Media Type status.
func AllowContentType(contentTypes ...string) gin.HandlerFunc {
	allowedContentTypes := make(map[string]struct{}, len(contentTypes))
	for _, c := range contentTypes {
		allowedContentTypes[strings.TrimSpace(strings.ToLower(c))] = struct{}{}
	}

	return func(c *gin.Context) {
		s := strings.ToLower(strings.TrimSpace(c.ContentType()))
		if _, ok := allowedContentTypes[s]; ok {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusUnsupportedMediaType)
	}
}
