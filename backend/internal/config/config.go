package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	return &Config{
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "report_db"),
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		UploadDir:            getEnv("UPLOAD_DIR", "./uploads"),
		BaseURL:              getEnv("BASE_URL", "http://localhost:8090"),
		CorsAllowedOrigins:   strings.Split(getEnv("CORS_ALLOWED_ORIGINS", ""), ","),
		CorsAllowCredentials: getEnv("CORS_ALLOW_CREDENTIALS", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
