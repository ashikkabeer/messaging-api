package queue

import (
	"os"
	"strconv"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewConfig() *Config {
	port, _ := strconv.Atoi(getEnvOrDefault("RabbitMQ_PORT", "5672"))
	return &Config{
		Host:     getEnvOrDefault("RabbitMQ_HOST", "localhost"),
		Port:     port,
		User:     getEnvOrDefault("RabbitMQ_USER", "guest"),
		Password: getEnvOrDefault("RabbitMQ_PASSWORD", "guest"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}