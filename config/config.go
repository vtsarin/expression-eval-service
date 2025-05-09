package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Logging  LoggingConfig
	History  HistoryConfig
	Security SecurityConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level      string
	Format     string
	OutputPath string
}

// HistoryConfig holds history-related configuration
type HistoryConfig struct {
	MaxSize int
	TTL     time.Duration
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	RateLimit      int
	MaxRequestSize int64
	AllowedOrigins []string
}

// New creates a new Config with values from environment variables
func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Logging: LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			OutputPath: getEnv("LOG_OUTPUT", "stdout"),
		},
		History: HistoryConfig{
			MaxSize: getEnvAsInt("HISTORY_MAX_SIZE", 1000),
			TTL:     getEnvAsDuration("HISTORY_TTL", 24*time.Hour),
		},
		Security: SecurityConfig{
			RateLimit:      getEnvAsInt("RATE_LIMIT", 100),
			MaxRequestSize: getEnvAsInt64("MAX_REQUEST_SIZE", 1024*1024), // 1MB
			AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
		},
	}
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return []string{value}
	}
	return defaultValue
}
