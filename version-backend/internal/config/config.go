package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Osquery  OsqueryConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// OsqueryConfig holds osquery configuration
type OsqueryConfig struct {
	SocketPath    string
	QueryInterval int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "7070"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "osquery"),
			Password: getEnv("DB_PASSWORD", "osquery_password"),
			DBName:   getEnv("DB_NAME", "osquery_data"),
		},
		Osquery: OsqueryConfig{
			SocketPath:    getEnv("OSQUERY_SOCKET", "/var/osquery/osquery.em"),
			QueryInterval: getEnvAsInt("QUERY_INTERVAL", 300),
		},
	}, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
