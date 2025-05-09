package routes

import (
	"expression-eval-service/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, evaluateController *controllers.EvaluateController) {
	// API routes
	api := router.Group("/api")
	{
		// Evaluation endpoints
		eval := api.Group("/evaluate")
		{
			// Single expression evaluation
			eval.POST("/single", evaluateController.Evaluate)
			// Batch expression evaluation
			eval.POST("/batch", evaluateController.EvaluateBatch)
			// History endpoint
			eval.GET("/history", evaluateController.GetHistory)
		}
	}
}
