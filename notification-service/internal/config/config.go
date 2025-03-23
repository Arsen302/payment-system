package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	KafkaBootstrapServers   string
	KafkaPaymentTopic       string
	KafkaGroupID            string
	LogLevel                string
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

	return &Config{
		KafkaBootstrapServers:   getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaPaymentTopic:       getEnv("KAFKA_PAYMENT_TOPIC", "payments"),
		KafkaGroupID:            getEnv("KAFKA_GROUP_ID", "notification-service"),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
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