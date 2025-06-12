package config

import (
	"os"
)

type Config struct {
	SchemaRegistryURL string
	KafkaHost         string
	ClientId          string
	HttpPort          string
	GrpcPort          string
	DBHost            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBPort            string
}

func LoadConfig() (*Config, error) {
	return &Config{
		SchemaRegistryURL: getEnvOrDefault("SCHEMA_REGISTRY_URL", "http://localhost:8085"),
		KafkaHost:         getEnvOrDefault("KAFKA_HOST", "localhost:19092"),
		ClientId:          getEnvOrDefault("CLIENT_ID", "organization-service"),
		HttpPort:          getEnvOrDefault("HTTP_PORT", ":8000"),
		GrpcPort:          getEnvOrDefault("GRPC_PORT", ":50051"),
		DBHost:            getEnvOrDefault("DB_HOST", "localhost"),
		DBUser:            getEnvOrDefault("DB_USER", "postgres"),
		DBPassword:        getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:            getEnvOrDefault("DB_NAME", "postgres"),
		DBPort:            getEnvOrDefault("DB_PORT", "5432"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
