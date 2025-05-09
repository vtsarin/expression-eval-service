package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware creates a middleware that logs request details
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		logger.Info("Request completed",
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.Duration("latency", latency),
		)
	}
}
