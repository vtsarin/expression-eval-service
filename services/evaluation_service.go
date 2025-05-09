package services

import (
	"context"
	"sync"
	"time"

	"go-service/evaluator"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Evaluation represents a single expression evaluation with its metadata
type Evaluation struct {
	ID         string    `json:"id"`               // Unique identifier for the evaluation
	Expression string    `json:"expression"`       // The mathematical expression to evaluate
	Result     float64   `json:"result,omitempty"` // The computed result (if successful)
	Error      error     `json:"error,omitempty"`  // Error message (if evaluation failed)
	Timestamp  time.Time `json:"timestamp"`        // When the evaluation was performed
}

// EvaluationService manages expression evaluations and their history
// It provides thread-safe operations for evaluating expressions and retrieving history
type EvaluationService struct {
	history []Evaluation      // In-memory storage for evaluation history
	mu      sync.RWMutex      // Mutex for thread-safe access to history
	logger  *zap.Logger       // Logger for tracking operations
	parser  *evaluator.Parser // Parser for parsing expressions
}

// NewEvaluationService creates a new instance of EvaluationService
// It initializes an empty history and sets up the logger
func NewEvaluationService(logger *zap.Logger) *EvaluationService {
	return &EvaluationService{
		history: make([]Evaluation, 0),
		logger:  logger,
		parser:  evaluator.NewParser(),
	}
}

// Evaluate evaluates an expression and stores the result in history
// It handles both successful evaluations and errors, storing both in history
func (s *EvaluationService) Evaluate(ctx context.Context, expression string) (Evaluation, error) {
	s.logger.Info("Starting evaluation of expression",
		zap.String("expression", expression),
	)

	// Create evaluation record
	eval := Evaluation{
		ID:         uuid.New().String(),
		Expression: expression,
		Timestamp:  time.Now(),
	}

	// Parse and evaluate expression
	expr, err := s.parser.Parse(expression)
	if err != nil {
		eval.Error = err
		s.addToHistory(ctx, eval)
		s.logger.Error("Failed to parse expression",
			zap.String("expression", expression),
			zap.Error(err),
		)
		return eval, err
	}

	// Evaluate the expression
	result, err := expr.Evaluate()
	if err != nil {
		eval.Error = err
		s.addToHistory(ctx, eval)
		s.logger.Error("Failed to evaluate expression",
			zap.String("expression", expression),
			zap.Error(err),
		)
		return eval, err
	}

	eval.Result = result
	s.addToHistory(ctx, eval)

	s.logger.Info("Successfully evaluated expression",
		zap.String("expression", expression),
		zap.Float64("result", result),
		zap.String("id", eval.ID),
	)

	return eval, nil
}

// GetHistory retrieves the evaluation history with pagination
func (s *EvaluationService) GetHistory(ctx context.Context, page, pageSize int) ([]Evaluation, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Calculate pagination
	total := len(s.history)
	start := (page - 1) * pageSize
	end := start + pageSize

	// Adjust end index if it exceeds the slice length
	if end > total {
		end = total
	}

	// Return empty slice if start index is out of bounds
	if start >= total {
		return []Evaluation{}, total, nil
	}

	// Get paginated slice
	history := make([]Evaluation, end-start)
	copy(history, s.history[start:end])

	s.logger.Info("Retrieved paginated history",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Int("total", total),
	)

	return history, total, nil
}

// addToHistory adds an evaluation to the history
// It uses a mutex to ensure thread-safe access
func (s *EvaluationService) addToHistory(ctx context.Context, eval Evaluation) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.history = append(s.history, eval)
	s.logger.Info("Added evaluation to history",
		zap.String("id", eval.ID),
		zap.Int("total", len(s.history)),
	)
}
