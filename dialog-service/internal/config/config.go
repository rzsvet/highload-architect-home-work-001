package config

import (
	"os"
)

type Config struct {
	ServerPort string
	MongoDBURI string
	DBName     string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		MongoDBURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:     getEnv("DB_NAME", "dialog_service"),
		JWTSecret:  getEnv("JWT_SECRET", "secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
