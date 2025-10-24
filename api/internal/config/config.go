package config

import (
	"os"
	"strconv"
	"strings"
)

type RedisConfig struct {
	// Standalone конфигурация
	Host     string
	Port     string
	Password string
	DB       int

	// Cluster конфигурация
	ClusterMode  bool
	ClusterNodes []string // ["host1:port1", "host2:port2", ...]

	// Общие настройки
	PoolSize     int
	MinIdleConns int
}

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

	// Redis configuration
	Redis RedisConfig

	// CORS configuration
	CORSAllowedOrigins []string
}

func LoadConfig() *Config {

	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:8080,http://localhost:8082")

	// Парсим Redis кластер ноды
	clusterNodes := strings.Split(getEnv("REDIS_CLUSTER_NODES", "127.0.0.1:6381,127.0.0.1:6382,127.0.0.1:6383,127.0.0.1:6384,127.0.0.1:6385,127.0.0.1:6386"), ",")
	if len(clusterNodes) == 1 && clusterNodes[0] == "" {
		clusterNodes = []string{}
	}

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		ServerPath:    getEnv("SERVER_PATH", "/api/v1"),
		ServerSwagger: getEnv("SERVER_SWAGGER", "enabled"),
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

		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvInt("REDIS_DB", 0),
			ClusterMode:  getEnvBool("REDIS_CLUSTER_MODE", true),
			ClusterNodes: clusterNodes,
			PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 5),
		},

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

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
