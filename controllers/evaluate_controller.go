package controllers

import (
	"net/http"
	"strconv"

	"expression-eval-service/errors"
	"expression-eval-service/models"
	"expression-eval-service/services"

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
	evaluationService *services.EvaluationService
	logger            *zap.Logger
}

// NewEvaluateController creates a new evaluation controller
func NewEvaluateController(evaluationService *services.EvaluationService, logger *zap.Logger) *EvaluateController {
	return &EvaluateController{
		evaluationService: evaluationService,
		logger:            logger,
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

	eval, err := c.evaluationService.Evaluate(ctx, req.Expression)
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

	if eval.Error != "" {
		response.Error = eval.Error
	}

	c.logger.Info("Evaluation successful",
		zap.String("id", eval.ID),
		zap.Float64("result", eval.Result),
	)

	ctx.JSON(http.StatusOK, response)
}

// EvaluateBatch handles batch evaluation requests
func (c *EvaluateController) EvaluateBatch(ctx *gin.Context) {
	var req models.BatchEvaluationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.SendError(ctx, err)
		return
	}

	results := c.evaluationService.EvaluateBatch(ctx, req.Expressions)
	errors.SendSuccess(ctx, "Batch evaluation completed", results)
}

// GetHistory handles GET requests to retrieve evaluation history
// It supports pagination and returns the history in reverse chronological order
func (c *EvaluateController) GetHistory(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errors.SendError(ctx, err)
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		errors.SendError(ctx, err)
		return
	}

	history, total, err := c.evaluationService.GetHistory(ctx, page, pageSize)
	if err != nil {
		errors.SendError(ctx, err)
		return
	}

	errors.SendSuccess(ctx, "History retrieved successfully", gin.H{
		"history": history,
		"total":   total,
	})
}
