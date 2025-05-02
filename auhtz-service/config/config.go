package config

import (
	"log"
	"os"
)

type Config struct {
	AuthzedHost   string `json:"authzed_host"`
	AuthzedApiKey string `json:"authzed_api_key"`
	KafkaHost     string `json:"kafka_host"`
	KafkaClientId string `json:"kafka_client_id"`
}

// LoadConfig loads configuration values from environment variables
func LoadConfig() *Config {
	missingVars := checkMissingVars([]string{"AUTHZED_HOST", "AUTHZED_API_KEY", "KAFKA_CLIENT_ID", "KAFKA_HOST"})

	if len(missingVars) > 0 {
		log.Fatalf("The following environment variables are missing: %v", missingVars)
	}

	config := &Config{
		AuthzedHost:   os.Getenv("AUTHZED_HOST"),
		AuthzedApiKey: os.Getenv("AUTHZED_API_KEY"),
		KafkaHost:     os.Getenv("KAFKA_HOST"),
		KafkaClientId: os.Getenv("KAFKA_CLIENT_ID"),
	}

	return config
}

// checkMissingVars checks for missing environment variables and returns a list of their names
func checkMissingVars(vars []string) []string {
	missingVars := []string{}
	for _, v := range vars {
		if os.Getenv(v) == "" {
			missingVars = append(missingVars, v)
		}
	}
	return missingVars
}
