package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	MongoURI   string
	MongoDB    string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8083),
		MongoURI:   getEnv("MONGODB_URL", "mongodb+srv://openpasswordopen:open@cluster0.dwuo6sl.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"),
		MongoDB:    getEnv("MONGO_DB", "user_service"),
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
