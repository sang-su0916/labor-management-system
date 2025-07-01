package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server settings
	Port    string
	Host    string
	GinMode string

	// Database settings
	DBType     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPath     string

	// JWT settings
	JWTSecret      string
	JWTExpiresHours int

	// File settings
	UploadPath    string
	DocumentsPath string

	// SMTP settings
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string

	// Company settings
	CompanyName    string
	CompanyAddress string
	CompanyPhone   string

	// Security settings
	AllowedOrigins string
	RateLimitReqs  int
	RateLimitWindow int

	// Logging
	LogLevel string
	LogFile  string
}

func LoadConfig() *Config {
	return &Config{
		// Server
		Port:    getEnv("PORT", "8080"),
		Host:    getEnv("HOST", "0.0.0.0"),
		GinMode: getEnv("GIN_MODE", "debug"),

		// Database
		DBType:     getEnv("DB_TYPE", "sqlite"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "labor_management"),
		DBPath:     getEnv("DB_PATH", "./labor_management.db"),

		// JWT
		JWTSecret:       getEnv("JWT_SECRET", "your_jwt_secret_key"),
		JWTExpiresHours: getEnvAsInt("JWT_EXPIRES_HOURS", 24),

		// Files
		UploadPath:    getEnv("UPLOAD_PATH", "./uploads"),
		DocumentsPath: getEnv("DOCUMENTS_PATH", "./documents"),

		// SMTP
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),

		// Company
		CompanyName:    getEnv("COMPANY_NAME", "테스트 회사"),
		CompanyAddress: getEnv("COMPANY_ADDRESS", "서울특별시 강남구 테헤란로 123"),
		CompanyPhone:   getEnv("COMPANY_PHONE", "02-1234-5678"),

		// Security
		AllowedOrigins:  getEnv("ALLOWED_ORIGINS", "*"),
		RateLimitReqs:   getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow: getEnvAsInt("RATE_LIMIT_WINDOW", 60),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),
		LogFile:  getEnv("LOG_FILE", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}