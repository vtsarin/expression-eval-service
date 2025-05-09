package controllers

import (
	"github.com/gin-gonic/gin"
)

// HealthController handles health check related endpoints
type HealthController struct{}

// NewHealthController creates a new instance of HealthController
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck handles the health check endpoint
func (h *HealthController) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
