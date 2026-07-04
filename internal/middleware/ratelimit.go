package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	tokens int
	last   time.Time
}

var (
	limiters = make(map[string]*rateLimiter)
	mu       sync.Mutex
)

// RateLimit is a simple token bucket rate limiter per IP.
// Max requests per minute.
func RateLimit(maxRequests int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		defer mu.Unlock()

		limiter, exists := limiters[ip]
		if !exists {
			limiter = &rateLimiter{
				tokens: maxRequests,
				last:   time.Now(),
			}
			limiters[ip] = limiter
		}

		now := time.Now()
		elapsed := now.Sub(limiter.last)

		// Replenish tokens based on time passed
		limiter.tokens += int(elapsed.Minutes()) * maxRequests
		if limiter.tokens > maxRequests {
			limiter.tokens = maxRequests
		}
		limiter.last = now

		if limiter.tokens <= 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		limiter.tokens--
		c.Next()
	}
}
