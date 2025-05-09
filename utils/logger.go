package utils

import (
	"time"

	"go-service/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

// NewLogger creates a new instance of Logger with Zap
func NewLogger() (interfaces.LoggerService, error) {
	zapConfig := zap.NewProductionConfig()
	log, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{log}, nil
}

// WithRequestId adds request ID to the logger context
func (l *Logger) WithRequestId(c *gin.Context) interfaces.LoggerService {
	requestId := c.MustGet("requestId").(string)
	return &Logger{l.With(zap.String("requestId", requestId))}
}

// Info logs an informational message with context
func (l *Logger) Info(c *gin.Context, msg string) {
	if c == nil {
		l.Logger.Info(msg)
		return
	}
	logFields := logFieldsFromContext(c)
	latency := elapsedSince(c.Value("startTime").(time.Time))

	l.With(zap.Any("context", logFields), zap.Duration("latency", latency)).Info(msg)
}

// Warn logs a warning message with context
func (l *Logger) Warn(c *gin.Context, msg string) {
	if c == nil {
		l.Logger.Warn(msg)
		return
	}
	logFields := logFieldsFromContext(c)
	latency := elapsedSince(c.Value("startTime").(time.Time))

	l.With(zap.Any("context", logFields), zap.Duration("latency", latency)).Warn(msg)
}

// Error logs an error message with context
func (l *Logger) Error(c *gin.Context, msg string, err error) {
	if c == nil {
		l.Logger.Error(msg, zap.Error(err))
		return
	}
	logFields := logFieldsFromContext(c)
	latency := elapsedSince(c.Value("startTime").(time.Time))

	l.With(zap.Any("context", logFields), zap.Duration("latency", latency)).Error(msg, zap.Error(err))
}

// Debug logs a debug message with context
func (l *Logger) Debug(c *gin.Context, msg string) {
	if c == nil {
		l.Logger.Debug(msg)
		return
	}
	logFields := logFieldsFromContext(c)
	latency := elapsedSince(c.Value("startTime").(time.Time))

	l.With(zap.Any("context", logFields), zap.Duration("latency", latency)).Debug(msg)
}

// logFieldsFromContext extracts relevant fields from the Gin context
func logFieldsFromContext(c *gin.Context) map[string]interface{} {
	requestID := c.Request.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
		c.Request.Header.Set("X-Request-ID", requestID)
	}

	startTime := c.GetTime("startTime")
	if startTime == (time.Time{}) {
		startTime = time.Now()
		c.Set("startTime", startTime)
	}

	return map[string]interface{}{
		"requestId":     requestID,
		"correlationId": c.GetString("correlationId"),
		"path":          c.Request.URL.Path,
		"controller":    c.HandlerName(),
	}
}

// GenerateRequestID generates and sets a new request ID in the context
func GenerateRequestID(c *gin.Context) {
	c.Set("requestId", uuid.New().String())
	c.Set("startTime", time.Now())
}

// elapsedSince calculates the time elapsed since the given start time
func elapsedSince(startTime time.Time) time.Duration {
	return time.Since(startTime)
}
