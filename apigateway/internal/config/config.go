package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}
