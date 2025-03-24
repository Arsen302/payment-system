package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Metrics  MetricsConfig
}

// ServerConfig holds the configuration for the HTTP server
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds the configuration for the database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds the configuration for JWT tokens
type JWTConfig struct {
	Secret     string
	ExpireMins int
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

	expireMins, _ := strconv.Atoi(getEnv("JWT_EXPIRE_MINUTES", "60"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "auth"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-here"),
			ExpireMins: expireMins,
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