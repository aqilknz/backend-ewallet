package config

import (
	"os"
)

type Config struct {
	AppPort       string
	FrontendURL   string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPassword  string
	SMTPFromEmail string
}

func LoadConfig() *Config {
	return &Config{
		AppPort:       getEnv("PORT", "9000"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUser:      os.Getenv("SMTP_USER"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		SMTPFromEmail: os.Getenv("SMTP_FROM_EMAIL"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
