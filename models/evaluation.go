package models

import (
	"time"

	"github.com/google/uuid"
)

// BatchEvaluationRequest represents a request to evaluate multiple expressions
type BatchEvaluationRequest struct {
	Expressions []string `json:"expressions" binding:"required,min=1"`
}

// Evaluation represents a single expression evaluation result
type Evaluation struct {
	ID         string    `json:"id"`
	Expression string    `json:"expression"`
	Result     float64   `json:"result,omitempty"`
	Error      string    `json:"error,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// BatchEvaluationResponse represents the response for batch evaluation
type BatchEvaluationResponse struct {
	Results []Evaluation `json:"results"`
}

// NewEvaluation creates a new evaluation with a generated UUID and current timestamp
func NewEvaluation(expression string) Evaluation {
	return Evaluation{
		ID:         uuid.New().String(),
		Expression: expression,
		Timestamp:  time.Now().UTC(),
	}
}
