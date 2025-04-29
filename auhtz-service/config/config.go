package config

import (
	"log"
	"os"
)

type Config struct {
	AuthzedHost   string `json:"authzed_host"`
	AuthzedApiKey string `json:"authzed_api_key"`
}

// LoadConfig loads configuration values from environment variables
func LoadConfig() *Config {
	missingVars := checkMissingVars([]string{"AUTHZED_HOST", "AUTHZED_API_KEY"})

	if len(missingVars) > 0 {
		log.Fatalf("The following environment variables are missing: %v", missingVars)
	}

	config := &Config{
		AuthzedHost:   os.Getenv("AUTHZ_HOST"),
		AuthzedApiKey: os.Getenv("AUTHZ_API_KEY"),
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
