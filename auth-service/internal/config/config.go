package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	Port                    int
	PostgresURL             string
	RedisURL                string
	JWTSecret               string
	JWTExpiryHours          int
	LogLevel                string
	GrpcMaxReceiveMessageSize int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		// Try to find and load .env files in order of precedence
		possibleEnvFiles := []string{
			".env",
			".env.local",
			filepath.Join("..", "..", ".env"),
			filepath.Join("..", "..", ".env.local"),
		}

		for _, file := range possibleEnvFiles {
			if _, err := os.Stat(file); err == nil {
				_ = godotenv.Load(file)
				break
			}
		}
	} else {
		_ = godotenv.Load(envFile)
	}

	port, err := strconv.Atoi(getEnv("PORT", "50051"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiry hours: %w", err)
	}

	grpcMaxSize, err := strconv.Atoi(getEnv("GRPC_MAX_RECEIVE_MESSAGE_SIZE", "4194304")) // 4MB
	if err != nil {
		return nil, fmt.Errorf("invalid gRPC max message size: %w", err)
	}

	return &Config{
		Port:                    port,
		PostgresURL:             getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"),
		RedisURL:                getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:               getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHours:          jwtExpiryHours,
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		GrpcMaxReceiveMessageSize: grpcMaxSize,
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
} 