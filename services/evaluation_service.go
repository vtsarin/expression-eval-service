package services

import (
	"fmt"
	"sync"
	"time"

	"go-service/evaluator"
	"go-service/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Evaluation represents a single expression evaluation with its metadata
type Evaluation struct {
	ID         string    `json:"id"`               // Unique identifier for the evaluation
	Expression string    `json:"expression"`       // The mathematical expression to evaluate
	Result     float64   `json:"result,omitempty"` // The computed result (if successful)
	Error      string    `json:"error,omitempty"`  // Error message (if evaluation failed)
	Timestamp  time.Time `json:"timestamp"`        // When the evaluation was performed
}

// EvaluationService manages expression evaluations and their history
// It provides thread-safe operations for evaluating expressions and retrieving history
type EvaluationService struct {
	history []Evaluation             // In-memory storage for evaluation history
	mu      sync.RWMutex             // Mutex for thread-safe access to history
	logger  interfaces.LoggerService // Logger for tracking operations
}

// NewEvaluationService creates a new instance of EvaluationService
// It initializes an empty history and sets up the logger
func NewEvaluationService(logger interfaces.LoggerService) *EvaluationService {
	return &EvaluationService{
		history: make([]Evaluation, 0),
		logger:  logger,
	}
}

// Evaluate evaluates an expression and stores the result in history
// It handles both successful evaluations and errors, storing both in history
func (s *EvaluationService) Evaluate(c *gin.Context, expression string) (Evaluation, error) {
	// Generate a new evaluation with unique ID and timestamp
	eval := Evaluation{
		ID:         uuid.New().String(),
		Expression: expression,
		Timestamp:  time.Now(),
	}

	s.logger.Info(c, fmt.Sprintf("Starting evaluation of expression: %s (ID: %s)", expression, eval.ID))

	// Parse the expression
	parser := evaluator.NewParser(expression)
	expr, err := parser.Parse()
	if err != nil {
		s.logger.Error(c, fmt.Sprintf("Failed to parse expression: %s (ID: %s)", expression, eval.ID), err)
		eval.Error = err.Error()
		s.addToHistory(c, eval)
		return eval, err
	}

	// Evaluate the parsed expression
	result, err := expr.Evaluate()
	if err != nil {
		s.logger.Error(c, fmt.Sprintf("Failed to evaluate expression: %s (ID: %s)", expression, eval.ID), err)
		eval.Error = err.Error()
		s.addToHistory(c, eval)
		return eval, err
	}

	// Store successful evaluation
	eval.Result = result
	s.logger.Info(c, fmt.Sprintf("Successfully evaluated expression: %s = %f (ID: %s)",
		expression, result, eval.ID))
	s.addToHistory(c, eval)
	return eval, nil
}

// GetHistory returns all evaluations, most recent first
// It creates a copy of the history to avoid race conditions
func (s *EvaluationService) GetHistory(c *gin.Context) []Evaluation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.logger.Info(c, fmt.Sprintf("Retrieving evaluation history (count: %d)", len(s.history)))

	// Create a copy of the history to avoid race conditions
	history := make([]Evaluation, len(s.history))
	copy(history, s.history)

	// Reverse the slice to get most recent first
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history
}

// addToHistory adds an evaluation to the history
// It uses a mutex to ensure thread-safe access
func (s *EvaluationService) addToHistory(c *gin.Context, eval Evaluation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, eval)
	s.logger.Info(c, fmt.Sprintf("Added evaluation to history (ID: %s, Total: %d)",
		eval.ID, len(s.history)))
}
