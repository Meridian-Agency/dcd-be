package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            int
	DBHost          string
	DBPort          int
	DBDatabase      string
	DBUsername      string
	DBPassword      string
	DBSchema        string
	SMTPHost        string
	SMTPPort        int
	S3Endpoint      string
	WhatsAppMockURL string
}

func Load() (*Config, error) {
	// Attempt to load .env file if it exists (for local development)
	_ = godotenv.Load()

	cfg := &Config{}
	var err error

	cfg.Port, err = getEnvInt("PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost, err = getEnvRequired("DB_HOST")
	if err != nil {
		return nil, err
	}

	cfg.DBPort, err = getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DBDatabase, err = getEnvRequired("DB_DATABASE")
	if err != nil {
		return nil, err
	}

	cfg.DBUsername, err = getEnvRequired("DB_USERNAME")
	if err != nil {
		return nil, err
	}

	cfg.DBPassword, err = getEnvRequired("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	cfg.DBSchema = getEnv("DB_SCHEMA", "public")
	cfg.SMTPHost = getEnv("SMTP_HOST", "mailpit")

	cfg.SMTPPort, err = getEnvInt("SMTP_PORT", 1025)
	if err != nil {
		return nil, err
	}

	cfg.S3Endpoint = getEnv("S3_ENDPOINT", "http://minio:9000")
	cfg.WhatsAppMockURL = getEnv("WHATSAPP_MOCK_URL", "http://whatsapp-mock:3000/api/messages")

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvRequired(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return "", fmt.Errorf("required environment variable %s is missing", key)
	}
	return val, nil
}

func getEnvInt(key string, defaultVal int) (int, error) {
	valStr, ok := os.LookupEnv(key)
	if !ok || valStr == "" {
		return defaultVal, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be a valid integer: %w", key, err)
	}
	return val, nil
}
