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
	KafkaBootstrapServers   string
	KafkaPaymentTopic       string
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

	port, err := strconv.Atoi(getEnv("PORT", "50052"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	grpcMaxSize, err := strconv.Atoi(getEnv("GRPC_MAX_RECEIVE_MESSAGE_SIZE", "4194304")) // 4MB
	if err != nil {
		return nil, fmt.Errorf("invalid gRPC max message size: %w", err)
	}

	return &Config{
		Port:                    port,
		PostgresURL:             getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/payment_db?sslmode=disable"),
		KafkaBootstrapServers:   getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaPaymentTopic:       getEnv("KAFKA_PAYMENT_TOPIC", "payments"),
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