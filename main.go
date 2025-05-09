package main

import (
	"log"

	"go-service/controllers"
	"go-service/middlewares"
	"go-service/services"
	"go-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger, err := utils.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Initialize services
	evalService := services.NewEvaluationService(logger)

	// Initialize Gin router
	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())                    // Add recovery middleware
	router.Use(middlewares.BindCorrelationId())   // Bind correlation ID
	router.Use(middlewares.LatencyLogger(logger)) // Log latency

	// Initialize controllers
	healthController := controllers.NewHealthController()
	evaluateController := controllers.NewEvaluateController(evalService, logger)

	// API routes
	api := router.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", healthController.HealthCheck)

		// Expression evaluation endpoints
		api.POST("/evaluate", evaluateController.Evaluate)
		api.GET("/history", evaluateController.GetHistory)
	}

	// Start the server
	logger.Info(nil, "Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error(nil, "Failed to start server", err)
	}
}
