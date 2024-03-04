package webx

import (
	"net/http"
	"sync"

	"gin.go.dev/internal/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/time/rate"
)

// Authenticated middleware func to ensure logged in, redirects to log-in if not.
// This calls CurrentUser first, so you don't need to chain both.
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user"); !exists {
			CurrentUser()(c)
		}

		if _, exists := c.Get("user"); !exists {
			htmx := &HTMXHelper{Request: c.Request, Response: c.Writer}
			if htmx.IsHTMXRequest() {
				htmx.SetRedirect("/auth/login")
				c.Status(http.StatusNoContent)
			} else {
				c.Redirect(http.StatusFound, "/auth/login")
			}
			c.Abort()
		}
	}
}

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

// CurrentUser middleware func to set the current active user.
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := c.MustGet("session").(sessions.Session)
		queries := c.MustGet("queries").(*db.Queries)

		userID, ok := session.Get("user_id").([16]byte)
		if !ok {
			return
		}

		userUUID := pgtype.UUID{Bytes: userID, Valid: true}
		user, err := queries.GetUserByID(c.Request.Context(), userUUID)
		if err != nil || !user.IsActive {
			return
		}

		c.Set("user", user)
	}
}

// ipLimiter holds the rate limiter for each IP
type ipLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

// newIPLimiter creates a new IP-based rate limiter
func newIPLimiter() *ipLimiter {
	return &ipLimiter{
		ips: make(map[string]*rate.Limiter),
	}
}

// getLimiter returns the rate limiter for the given IP, creating a new one if necessary
func (l *ipLimiter) getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		l.ips[ip] = limiter
	}

	return limiter
}

// RateLimiter returns a Gin middleware for rate limiting with the specified rate and burst size
func RateLimiter(r rate.Limit, burst int) gin.HandlerFunc {
	limiter := newIPLimiter()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		lim := limiter.getLimiter(ip, r, burst)

		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}
