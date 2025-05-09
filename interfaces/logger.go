package interfaces

import (
	"context"

	"go.uber.org/zap"
)

// LoggerService defines the interface for logging
type LoggerService interface {
	// WithContext adds context fields to the logger
	WithContext(ctx context.Context) *zap.Logger

	// Info logs an info message with context
	Info(ctx context.Context, msg string, fields ...zap.Field)

	// Error logs an error message with context
	Error(ctx context.Context, msg string, err error, fields ...zap.Field)

	// Warn logs a warning message with context
	Warn(ctx context.Context, msg string, fields ...zap.Field)

	// Debug logs a debug message with context
	Debug(ctx context.Context, msg string, fields ...zap.Field)
}
