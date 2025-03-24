package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Kafka    KafkaConfig
	Email    EmailConfig
	SMS      SMSConfig
	Metrics  MetricsConfig
}

// ServerConfig holds the configuration for the HTTP server
type ServerConfig struct {
	Port string
}

// KafkaConfig holds the configuration for Kafka
type KafkaConfig struct {
	Brokers []string
	Topics  []string
	GroupID string
}

// EmailConfig holds the configuration for email notifications
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
}

// SMSConfig holds the configuration for SMS notifications
type SMSConfig struct {
	Provider  string
	APIKey    string
	APISecret string
	FromPhone string
}

// MetricsConfig holds the configuration for Prometheus metrics
type MetricsConfig struct {
	Enabled bool
	Port    string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
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
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			Topics:  strings.Split(getEnv("KAFKA_TOPICS", "payments,auth"), ","),
			GroupID: getEnv("KAFKA_GROUP_ID", "notification-service"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@payment-system.com"),
		},
		SMS: SMSConfig{
			Provider:  getEnv("SMS_PROVIDER", "twilio"),
			APIKey:    getEnv("SMS_API_KEY", ""),
			APISecret: getEnv("SMS_API_SECRET", ""),
			FromPhone: getEnv("SMS_FROM_PHONE", ""),
		},
		Metrics: MetricsConfig{
			Enabled: getEnvAsBool("METRICS_ENABLED", true),
			Port:    getEnv("METRICS_PORT", "9090"),
		},
	}, nil
}

// Helper function to get an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as a boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
} 