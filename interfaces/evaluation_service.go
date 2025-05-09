package interfaces

import (
	"context"
	"time"
)

// Evaluation represents a single expression evaluation
type Evaluation struct {
	ID         string
	Expression string
	Result     float64
	Error      error
	Timestamp  time.Time
}

// EvaluationService defines the interface for expression evaluation
type EvaluationService interface {
	// Evaluate evaluates an expression and stores the result
	Evaluate(ctx context.Context, expression string) (Evaluation, error)

	// GetHistory retrieves the evaluation history with pagination
	GetHistory(ctx context.Context, page, pageSize int) ([]Evaluation, int, error)
}
