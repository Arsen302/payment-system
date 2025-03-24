package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents service configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Kafka    KafkaConfig
	Metrics  MetricsConfig
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// AuthConfig represents auth service configuration
type AuthConfig struct {
	URL string
}

// KafkaConfig represents Kafka configuration
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled bool
	Port    string
}

// Load loads configuration from environment variables
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "payment"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			URL: getEnv("AUTH_SERVICE_URL", "http://auth-service:8080"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKER", "kafka:9092")},
			Topic:   getEnv("KAFKA_TOPIC", "payment-events"),
		},
		Metrics: MetricsConfig{
			Enabled: getBoolEnv("METRICS_ENABLED", true),
			Port:    getEnv("METRICS_PORT", "9090"),
		},
	}
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get boolean environment variable with default value
func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	
	return boolValue
} 