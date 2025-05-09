package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-service/config"
	"go-service/controllers"
	"go-service/middlewares"
	"go-service/services"
	"go-service/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Initialize logger
	logger, err := utils.NewLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize services
	evalService := services.NewEvaluationService(logger.Logger)

	// Initialize controllers
	evaluateController := controllers.NewEvaluateController(evalService, logger.Logger)

	// Configure Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.LoggerMiddleware(logger.Logger))
	router.Use(middlewares.RateLimitMiddleware(float64(cfg.Security.RateLimit), float64(cfg.Security.RateLimit), logger.Logger))
	router.Use(middlewares.CORSMiddleware(cfg.Security.AllowedOrigins))

	// Configure API routes
	api := router.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Evaluation endpoints
		api.POST("/evaluate", evaluateController.Evaluate)
		api.GET("/history", evaluateController.GetHistory)
	}

	// Configure HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Logger.Info("Starting server",
			zap.Int("port", cfg.Server.Port),
			zap.Duration("read_timeout", cfg.Server.ReadTimeout),
			zap.Duration("write_timeout", cfg.Server.WriteTimeout),
			zap.Duration("idle_timeout", cfg.Server.IdleTimeout),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("Failed to start server",
				zap.Error(err),
			)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Error("Server forced to shutdown",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Logger.Info("Server exited properly")
}
