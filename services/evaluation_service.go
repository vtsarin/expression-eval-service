package services

import (
	"context"
	"sync"
	"time"

	"expression-eval-service/evaluator"
	"expression-eval-service/models"

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
	history []models.Evaluation // In-memory storage for evaluation history
	mu      sync.RWMutex        // Mutex for thread-safe access to history
	logger  *zap.Logger         // Logger for tracking operations
	parser  *evaluator.Parser   // Parser for parsing expressions
}

// NewEvaluationService creates a new instance of EvaluationService
// It initializes an empty history and sets up the logger
func NewEvaluationService(logger *zap.Logger) *EvaluationService {
	return &EvaluationService{
		history: make([]models.Evaluation, 0),
		logger:  logger,
		parser:  evaluator.NewParser(),
	}
}

// Evaluate evaluates an expression and stores the result in history
// It handles both successful evaluations and errors, storing both in history
func (s *EvaluationService) Evaluate(ctx context.Context, expression string) (models.Evaluation, error) {
	s.logger.Info("Starting evaluation of expression",
		zap.String("expression", expression),
	)

	// Create evaluation record
	eval := models.NewEvaluation(expression)

	// Parse and evaluate expression
	expr, err := s.parser.Parse(expression)
	if err != nil {
		eval.Error = err.Error()
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
		eval.Error = err.Error()
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

// EvaluateBatch evaluates multiple expressions concurrently
func (s *EvaluationService) EvaluateBatch(ctx context.Context, expressions []string) models.BatchEvaluationResponse {
	s.logger.Info("Starting batch evaluation",
		zap.Int("expression_count", len(expressions)))

	results := make([]models.Evaluation, 0, len(expressions))
	var wg sync.WaitGroup
	resultChan := make(chan models.Evaluation, len(expressions))

	// Launch a goroutine for each expression
	for i, expr := range expressions {
		wg.Add(1)
		go func(index int, expression string) {
			defer wg.Done()

			eval := models.NewEvaluation(expression)
			ast, err := s.parser.Parse(expression)
			if err != nil {
				eval.Error = err.Error()
				resultChan <- eval
				return
			}

			result, err := ast.Evaluate()
			if err != nil {
				eval.Error = err.Error()
			} else {
				eval.Result = result
			}
			resultChan <- eval
		}(i, expr)
	}

	// Wait for all evaluations to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		results = append(results, result)
		// Store in history
		s.mu.Lock()
		s.history = append(s.history, result)
		s.mu.Unlock()
	}

	s.logger.Info("Completed batch evaluation",
		zap.Int("result_count", len(results)))

	return models.BatchEvaluationResponse{
		Results: results,
	}
}

// GetHistory retrieves the evaluation history with pagination
func (s *EvaluationService) GetHistory(ctx context.Context, page, pageSize int) ([]models.Evaluation, int, error) {
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
		return []models.Evaluation{}, total, nil
	}

	// Get paginated slice
	history := make([]models.Evaluation, end-start)
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
func (s *EvaluationService) addToHistory(ctx context.Context, eval models.Evaluation) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.history = append(s.history, eval)
	s.logger.Info("Added evaluation to history",
		zap.String("id", eval.ID),
		zap.Int("total", len(s.history)),
	)
}
