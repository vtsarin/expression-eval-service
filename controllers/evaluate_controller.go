package controllers

import (
	"net/http"
	"strconv"

	"go-service/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EvaluateRequest represents the request body for expression evaluation
type EvaluateRequest struct {
	Expression string `json:"expression" binding:"required"` // The mathematical expression to evaluate
}

// EvaluateResponse represents the response for expression evaluation
type EvaluateResponse struct {
	ID         string  `json:"id"`               // Unique identifier for the evaluation
	Expression string  `json:"expression"`       // The evaluated expression
	Result     float64 `json:"result,omitempty"` // The computed result (if successful)
	Error      string  `json:"error,omitempty"`  // Error message (if evaluation failed)
	Timestamp  string  `json:"timestamp"`        // When the evaluation was performed
}

// EvaluateController handles HTTP requests for expression evaluation
// It provides endpoints for evaluating expressions and retrieving history
type EvaluateController struct {
	evalService *services.EvaluationService // Service for evaluating expressions
	logger      *zap.Logger                 // Logger for tracking operations
}

// NewEvaluateController creates a new instance of EvaluateController
// It initializes the controller with both the evaluation service and logger
func NewEvaluateController(evalService *services.EvaluationService, logger *zap.Logger) *EvaluateController {
	return &EvaluateController{
		evalService: evalService,
		logger:      logger,
	}
}

// Evaluate handles POST requests to evaluate expressions
// It validates the request, evaluates the expression, and returns the result
func (c *EvaluateController) Evaluate(ctx *gin.Context) {
	var req EvaluateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("Invalid request body",
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	c.logger.Info("Received evaluation request",
		zap.String("expression", req.Expression),
	)

	eval, err := c.evalService.Evaluate(ctx, req.Expression)
	if err != nil {
		c.logger.Error("Evaluation failed",
			zap.String("expression", req.Expression),
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := EvaluateResponse{
		ID:         eval.ID,
		Expression: eval.Expression,
		Result:     eval.Result,
		Timestamp:  eval.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}

	if eval.Error != nil {
		response.Error = eval.Error.Error()
	}

	c.logger.Info("Evaluation successful",
		zap.String("id", eval.ID),
		zap.Float64("result", eval.Result),
	)

	ctx.JSON(http.StatusOK, response)
}

// GetHistory handles GET requests to retrieve evaluation history
// It supports pagination and returns the history in reverse chronological order
func (c *EvaluateController) GetHistory(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	c.logger.Info("Received history request",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
	)

	history, total, err := c.evalService.GetHistory(ctx, page, pageSize)
	if err != nil {
		c.logger.Error("Failed to retrieve history",
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	// Convert evaluations to response format
	response := make([]EvaluateResponse, len(history))
	for i, eval := range history {
		response[i] = EvaluateResponse{
			ID:         eval.ID,
			Expression: eval.Expression,
			Result:     eval.Result,
			Timestamp:  eval.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}
		if eval.Error != nil {
			response[i].Error = eval.Error.Error()
		}
	}

	c.logger.Info("Retrieved history successfully",
		zap.Int("count", len(history)),
		zap.Int("total", total),
	)

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}
