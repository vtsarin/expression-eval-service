package errors

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// ErrorMetrics tracks error rates and types
type ErrorMetrics struct {
	mu            sync.RWMutex
	errorCounts   map[string]int64
	lastResetTime time.Time
	logger        *zap.Logger
}

// NewErrorMetrics creates a new ErrorMetrics instance
func NewErrorMetrics(logger *zap.Logger) *ErrorMetrics {
	return &ErrorMetrics{
		errorCounts:   make(map[string]int64),
		lastResetTime: time.Now(),
		logger:        logger,
	}
}

// RecordError records an error occurrence
func (m *ErrorMetrics) RecordError(errorType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errorCounts[errorType]++
	m.logger.Info("Error recorded",
		zap.String("type", errorType),
		zap.Int64("count", m.errorCounts[errorType]),
	)
}

// GetErrorCounts returns the current error counts
func (m *ErrorMetrics) GetErrorCounts() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	counts := make(map[string]int64)
	for k, v := range m.errorCounts {
		counts[k] = v
	}
	return counts
}

// Reset resets the error counts
func (m *ErrorMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errorCounts = make(map[string]int64)
	m.lastResetTime = time.Now()
	m.logger.Info("Error metrics reset")
}

// GetLastResetTime returns the last reset time
func (m *ErrorMetrics) GetLastResetTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastResetTime
}
