package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

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
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}
