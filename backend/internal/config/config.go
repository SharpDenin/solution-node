package config

import (
	"os"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	ServerPort string

	UploadDir string
	BaseURL   string

	CorsAllowedOrigins   []string
	CorsAllowCredentials bool
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "report_db"),

		ServerPort: getEnv("SERVER_PORT", "8080"),

		UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
		BaseURL:   getEnv("BASE_URL", "http://localhost:8090"),

		CorsAllowedOrigins:   strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3010"), ","),
		CorsAllowCredentials: strings.ToLower(getEnv("CORS_ALLOW_CREDENTIALS", "true")) == "true",
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
