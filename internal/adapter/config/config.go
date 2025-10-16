package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database *DataBaseConfig
	Http     *HttpConfig
	Env      string
}

type HttpConfig struct {
	AllowedOrigins []string
	Port           string
	Url            string
}

type DataBaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	TimeZone string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func Init() (*Config, error) {
	_ = godotenv.Load()

	database := &DataBaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "postgres"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		TimeZone: getEnv("DB_TIMEZONE", "UTC"),
	}

	var allowedOriginList []string

	allowedOrigins := getEnv("HTTP_ALLOWED_ORIGINS", "*")
	if allowedOrigins == "*" && len(allowedOrigins) > 0 {
		allowedOriginList = strings.Split(allowedOrigins, ",")
	}

	http := &HttpConfig{
		AllowedOrigins: allowedOriginList,
		Port:           getEnv("HTTP_PORT", "8080"),
		Url:            getEnv("HTTP_URL", "http://localhost"),
	}

	return &Config{
		Database: database,
		Http:     http,
		Env:      getEnv("ENV", "development"),
	}, nil
}
