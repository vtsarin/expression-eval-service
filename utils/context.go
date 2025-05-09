package utils

import (
	"context"
)

const (
	// CorrelationIDKey is the key used to store correlation ID in context
	CorrelationIDKey = "correlation_id"
)

// GetCorrelationID gets the correlation ID from context
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}

// WithCorrelationID adds correlation ID to context
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, id)
}
