package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort       int
	GRPCPort         int
	JWTSecret        string
	TokenExpiryHours int
	MongoURI         string
	MongoDB          string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:       getEnvAsInt("SERVER_PORT", 8081),
		GRPCPort:         getEnvAsInt("GRPC_PORT", 50051),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiryHours: getEnvAsInt("TOKEN_EXPIRY_HOURS", 24),
		MongoURI:         getEnv("MONGODB_URL", "mongodb+srv://openpasswordopen:open@cluster0.dwuo6sl.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"),
		MongoDB:          getEnv("MONGO_DB", "auth_service"),
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
