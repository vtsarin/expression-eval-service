package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	rate       float64
	bucketSize float64
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
	logger     *zap.Logger
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate float64, bucketSize float64, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     bucketSize,
		lastRefill: time.Now(),
		logger:     logger,
	}
}

// refill adds tokens to the bucket based on elapsed time
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens = min(rl.bucketSize, rl.tokens+elapsed*rl.rate)
	rl.lastRefill = now
}

// allow checks if a request can be processed
func (rl *RateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// RateLimitMiddleware creates a middleware that limits requests based on rate
func RateLimitMiddleware(rate float64, bucketSize float64, logger *zap.Logger) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, bucketSize, logger)

	return func(c *gin.Context) {
		if !limiter.allow() {
			logger.Warn("Rate limit exceeded",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Rate limit exceeded",
				"error": gin.H{
					"code":    "E4297701",
					"message": "Too many requests",
					"details": "Please try again later",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
