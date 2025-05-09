package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetQueryInt gets an integer query parameter with a default value
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetQueryFloat gets a float query parameter with a default value
func GetQueryFloat(c *gin.Context, key string, defaultValue float64) float64 {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return floatValue
}

// GetQueryBool gets a boolean query parameter with a default value
func GetQueryBool(c *gin.Context, key string, defaultValue bool) bool {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
