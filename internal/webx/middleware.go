package webx

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/time/rate"
	"net/http"
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

// CurrentUser middleware func to get the current active user, redirects to login if not.
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		cc := c.MustGet("custom").(*CustomContext)

		getHTMX := func() *HTMXHelper {
			return &HTMXHelper{Request: c.Request, Response: c.Writer}
		}

		redirect := func() {
			if htmx := getHTMX(); htmx.IsHTMXRequest() {
				htmx.SetRedirect("/auth/login")
				c.Status(http.StatusNoContent)
			} else {
				c.Redirect(http.StatusFound, "/auth/login")
			}
			c.Abort()
		}

		userID, ok := cc.Session.Get("user_id").([16]byte)
		if !ok {
			redirect()
			return
		}

		userUUID := pgtype.UUID{Bytes: userID, Valid: true}
		user, err := cc.Queries.GetUserByID(c.Request.Context(), userUUID)
		if err != nil || !user.IsActive {
			redirect()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// RateLimiter creates a new rate limiter for each client based on IP address.
// It allows up to maxBurst requests instantly and refills at a rate of r tokens per second.
func RateLimiter(r rate.Limit, maxBurst int) gin.HandlerFunc {
	var limiter = rate.NewLimiter(r, maxBurst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
