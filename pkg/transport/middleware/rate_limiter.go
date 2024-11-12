package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

// RateLimiter provides a Gin middleware that limits request rates by IP address.
func RateLimiter(r rate.Limit, burst int) gin.HandlerFunc {
	var sm sync.Map

	return func(c *gin.Context) {
		ip := c.ClientIP()

		lim, _ := sm.LoadOrStore(ip, rate.NewLimiter(r, burst))
		limiter, ok := lim.(*rate.Limiter)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limiter.Allow() {
			c.Next()
			return
		}

		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}
