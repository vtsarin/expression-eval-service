package controllers

import (
	"fmt"

	"go-service/errors"
	"go-service/interfaces"
	"go-service/services"
	"go-service/utils"

	"github.com/gin-gonic/gin"
)

// EvaluateRequest represents the request body for the evaluate endpoint
// It contains the mathematical expression to be evaluated
type EvaluateRequest struct {
	Expression string `json:"expression" binding:"required"` // The mathematical expression to evaluate
}

// EvaluateResponse represents the response body for the evaluate endpoint
// It includes the evaluation ID and result (if successful) or error (if failed)
type EvaluateResponse struct {
	ID     string  `json:"id"`               // Unique identifier for the evaluation
	Result float64 `json:"result,omitempty"` // The computed result (if successful)
	Error  string  `json:"error,omitempty"`  // Error message (if evaluation failed)
}

// EvaluateController handles expression evaluation endpoints
// It provides endpoints for evaluating expressions and retrieving history
type EvaluateController struct {
	evalService *services.EvaluationService // Service for expression evaluation
	logger      interfaces.LoggerService    // Logger for tracking operations
}

// NewEvaluateController creates a new instance of EvaluateController
// It initializes the controller with the evaluation service and logger
func NewEvaluateController(evalService *services.EvaluationService, logger interfaces.LoggerService) *EvaluateController {
	return &EvaluateController{
		evalService: evalService,
		logger:      logger,
	}
}

// Evaluate handles the expression evaluation endpoint
// It validates the request, evaluates the expression, and returns the result
func (e *EvaluateController) Evaluate(c *gin.Context) {
	e.logger.Info(c, "Received evaluation request")

	// Parse and validate request body
	var req EvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.logger.Error(c, "Invalid request body", err)
		utils.SendError(c, errors.ErrInvalidRequestBody)
		return
	}

	e.logger.Info(c, fmt.Sprintf("Processing expression: %s", req.Expression))

	// Evaluate expression using service
	eval, err := e.evalService.Evaluate(c, req.Expression)
	if err != nil {
		e.logger.Error(c, fmt.Sprintf("Evaluation failed for expression: %s", req.Expression), err)
		utils.SendError(c, errors.NewCustomError(400, "E4007701", "Invalid expression", err.Error(), err))
		return
	}

	// Send success response
	e.logger.Info(c, fmt.Sprintf("Successfully evaluated expression: %s (ID: %s)",
		req.Expression, eval.ID))
	utils.SendSuccess(c, "Expression evaluated successfully", EvaluateResponse{
		ID:     eval.ID,
		Result: eval.Result,
	})
}

// GetHistory handles the history endpoint
// It retrieves and returns the evaluation history
func (e *EvaluateController) GetHistory(c *gin.Context) {
	e.logger.Info(c, "Received history request")

	// Get history from service
	history := e.evalService.GetHistory(c)

	e.logger.Info(c, fmt.Sprintf("Retrieved %d evaluations from history", len(history)))
	utils.SendSuccess(c, "History retrieved successfully", history)
}
