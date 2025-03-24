package config

import (
	"os"

	"github.com/firllyw/gogo/pkg/db"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database db.PostgresConfig
	OAuth    OAuthConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

// OAuthConfig holds the OAuth configuration
type OAuthConfig struct {
	Enabled      bool
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	JWTSecret    string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() Config {
	return Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: db.PostgresConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "user_service"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		OAuth: OAuthConfig{
			Enabled:      getEnvBool("OAUTH_ENABLED", false),
			ClientID:     getEnv("OAUTH_CLIENT_ID", ""),
			ClientSecret: getEnv("OAUTH_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/api/v1/auth/callback"),
			AuthURL:      getEnv("OAUTH_AUTH_URL", "https://accounts.google.com/o/oauth2/auth"),
			TokenURL:     getEnv("OAUTH_TOKEN_URL", "https://oauth2.googleapis.com/token"),
			UserInfoURL:  getEnv("OAUTH_USER_INFO_URL", "https://www.googleapis.com/oauth2/v3/userinfo"),
			JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		},
	}
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool reads a boolean environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}