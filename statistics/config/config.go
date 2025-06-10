package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	Mongo        MongoConfig        `json:"mongo"`
	EventStoreDB EventStoreDBConfig `json:"eventstoredb"`
	Kafka        KafkaConfig        `json:"kafka"`
	HTTP         HTTPConfig         `json:"http"`
	Logging      LoggingConfig      `json:"logging"`
	Metrics      MetricsConfig      `json:"metrics"`
}

// MongoConfig holds MongoDB configuration
type MongoConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

// EventStoreDBConfig holds EventStoreDB configuration
type EventStoreDBConfig struct {
	ConnectionString string `json:"connection_string"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	BootstrapServers string   `json:"bootstrap_servers"`
	GroupID          string   `json:"group_id"`
	Topics           []string `json:"topics"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	ListenAddress string `json:"listen_address"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// MetricsConfig holds metrics calculation configuration
type MetricsConfig struct {
	CalculationInterval string `json:"calculation_interval"`
	RetentionPeriod     string `json:"retention_period"`
}

// LoadConfig loads the configuration from a file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Default configuration
	config := &Config{
		Mongo: MongoConfig{
			URI:      "mongodb://localhost:27017",
			Database: "polypass_statistics",
		},
		EventStoreDB: EventStoreDBConfig{
			ConnectionString: "esdb://localhost:2113?tls=false",
		},
		Kafka: KafkaConfig{
			BootstrapServers: "localhost:9092",
			GroupID:          "statistics-service",
			Topics:           []string{"user-events", "password-events", "security-events"},
		},
		HTTP: HTTPConfig{
			ListenAddress: ":8080",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Metrics: MetricsConfig{
			CalculationInterval: "1h",
			RetentionPeriod:     "90d",
		},
	}

	// Load configuration from file if it exists
	if configPath != "" {
		file, err := os.Open(configPath)
		if err == nil {
			defer file.Close()
			decoder := json.NewDecoder(file)
			err = decoder.Decode(config)
			if err != nil {
				return nil, err
			}
		}
	}

	// Override with environment variables
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		config.Mongo.URI = uri
	}
	if db := os.Getenv("MONGO_DATABASE"); db != "" {
		config.Mongo.Database = db
	}
	if connStr := os.Getenv("EVENTSTOREDB_CONNECTION_STRING"); connStr != "" {
		config.EventStoreDB.ConnectionString = connStr
	}
	if bootstrap := os.Getenv("KAFKA_BOOTSTRAP"); bootstrap != "" {
		config.Kafka.BootstrapServers = bootstrap
	}
	if groupID := os.Getenv("KAFKA_GROUP_ID"); groupID != "" {
		config.Kafka.GroupID = groupID
	}
	if topics := os.Getenv("KAFKA_TOPICS"); topics != "" {
		config.Kafka.Topics = strings.Split(topics, ",")
	}
	if addr := os.Getenv("HTTP_LISTEN_ADDRESS"); addr != "" {
		config.HTTP.ListenAddress = addr
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Logging.Format = format
	}
	if interval := os.Getenv("METRICS_CALCULATION_INTERVAL"); interval != "" {
		config.Metrics.CalculationInterval = interval
	}
	if retention := os.Getenv("METRICS_RETENTION_PERIOD"); retention != "" {
		config.Metrics.RetentionPeriod = retention
	}

	return config, nil
}

// GetCalculationInterval returns the metrics calculation interval as a time.Duration
func (c *Config) GetCalculationInterval() (time.Duration, error) {
	return time.ParseDuration(c.Metrics.CalculationInterval)
}

// GetRetentionPeriod returns the metrics retention period as a time.Duration
func (c *Config) GetRetentionPeriod() (time.Duration, error) {
	return time.ParseDuration(c.Metrics.RetentionPeriod)
}
