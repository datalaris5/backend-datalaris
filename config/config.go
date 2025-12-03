package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	AppPort   string
	TenantKey string

	GCSBucket         string
	GCSProperties     string
	GCSPropertiesDocs string
	GCSPath           string
	GCSEnv            string
	JWTSecret         string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	EmailUser     string
	EmailPassword string
)

func LoadEnv() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found, fallback to system env")
	}

	AppPort = getEnv("APP_PORT", "8080")
	TenantKey = getEnv("TENANT_KEY", "tentant-key")

	// Database
	DBHost = getEnv("DB_HOST", "localhost")
	DBPort = getEnv("DB_PORT", "5432")
	DBUser = getEnv("DB_USER", "postgres")
	DBPassword = getEnv("DB_PASSWORD", "password")
	DBName = getEnv("DB_NAME", "app_db")
	DBSSLMode = getEnv("DB_SSLMODE", "disable")

	EmailUser = getEnv("EMAIL_USER", "app_db")
	EmailPassword = getEnv("EMAIL_PASSWORD", "disable")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
