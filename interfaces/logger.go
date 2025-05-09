package interfaces

import "github.com/gin-gonic/gin"

// LoggerService defines the interface for logging operations
type LoggerService interface {
	// Info logs an informational message
	Info(c *gin.Context, message string)
	// Error logs an error message
	Error(c *gin.Context, message string, err error)
	// Debug logs a debug message
	Debug(c *gin.Context, message string)
	// Warn logs a warning message
	Warn(c *gin.Context, message string)
	// WithRequestId adds request ID to the logger context
	WithRequestId(c *gin.Context) LoggerService
}
