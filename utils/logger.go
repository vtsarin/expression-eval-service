package utils

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Logger wraps zap.Logger to provide a consistent logging interface
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger() (*Logger, error) {
	config := zap.NewProductionConfig()
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	correlationID := GetCorrelationID(ctx)
	if correlationID == "" {
		correlationID = uuid.New().String()
		ctx = WithCorrelationID(ctx, correlationID)
	}
	return l.With(zap.String("correlation_id", correlationID))
}

// Info logs an info message with context
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Error logs an error message with context
func (l *Logger) Error(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	l.Logger.Error(msg, fields...)
}

// Warn logs a warning message with context
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Debug logs a debug message with context
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// WithRequestID adds request ID to the logger context
func (l *Logger) WithRequestID(c *gin.Context) *zap.Logger {
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
		c.Set("request_id", requestID)
	}
	return l.With(zap.String("request_id", requestID))
}

// WithGinContext adds Gin context fields to the logger
func (l *Logger) WithGinContext(c *gin.Context) *zap.Logger {
	if c == nil {
		return l.Logger
	}

	fields := []zap.Field{
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()),
	}

	if startTime, exists := c.Get("start_time"); exists {
		if t, ok := startTime.(time.Time); ok {
			fields = append(fields, zap.Duration("latency", time.Since(t)))
		}
	}

	return l.With(fields...)
}
