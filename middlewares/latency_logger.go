package middlewares

import (
	"fmt"
	"go-service/interfaces"
	"time"

	"github.com/gin-gonic/gin"
)

/*
 * LatencyLogger returns a middleware that logs detailed request information including latency
 * It includes information such as:
 * - API endpoint path
 * - HTTP method
 * - Request latency
 * - Status code
 * - Correlation ID
 * - Client IP
 */
func LatencyLogger(logger interfaces.LoggerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get request details
		path := c.Request.URL.Path
		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Get correlation ID from context
		correlationID, exists := c.Get("correlationId")
		if !exists {
			correlationID = "N/A"
		}

		// Log detailed request information
		logger.Info(c, fmt.Sprintf(
			"API Request - Method: %s, Path: %s, Status: %d, Latency: %v, ClientIP: %s, CorrelationID: %s",
			method,
			path,
			statusCode,
			latency,
			clientIP,
			correlationID,
		))
	}
}
