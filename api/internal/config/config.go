package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort    string
	ServerPath    string
	ServerSwagger string
	JWTSecret     string

	// Database configurations
	WriteDBHost     string
	WriteDBPort     string
	WriteDBName     string
	WriteDBUser     string
	WriteDBPassword string

	ReadDBHost     string
	ReadDBPort     string
	ReadDBName     string
	ReadDBUser     string
	ReadDBPassword string

	// CORS configuration
	CORSAllowedOrigins []string
}

func LoadConfig() *Config {

	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:8080,http://localhost:8082")

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		ServerPath:    getEnv("SERVER_PATH", "/api/v1"),
		ServerSwagger: getEnv("SERVER_SWAGGER", "disabled"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),

		WriteDBHost:     getEnv("WRITE_DB_HOST", "localhost"),
		WriteDBPort:     getEnv("WRITE_DB_PORT", "5432"),
		WriteDBName:     getEnv("WRITE_DB_NAME", "root"),
		WriteDBUser:     getEnv("WRITE_DB_USER", "root"),
		WriteDBPassword: getEnv("WRITE_DB_PASSWORD", "password"),

		ReadDBHost:     getEnv("READ_DB_HOST", "localhost"),
		ReadDBPort:     getEnv("READ_DB_PORT", "5432"),
		ReadDBName:     getEnv("READ_DB_NAME", "root"),
		ReadDBUser:     getEnv("READ_DB_USER", "root"),
		ReadDBPassword: getEnv("READ_DB_PASSWORD", "password"),

		CORSAllowedOrigins: strings.Split(corsOrigins, ","),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
