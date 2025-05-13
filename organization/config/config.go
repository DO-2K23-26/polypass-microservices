package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// The port on which the service will listen for incoming requests.
	GrpcPort   string `json:"grpcPort"`
	HttpPort   string `json:"httpPort"`
	KafkaHost  string `json:"kafkaHost"`
	ClientId   string `json:"clientId"`
	PgHost     string `json:"pgHost"`
	PgUser     string `json:"pgUser"`
	PgPassword string `json:"pgPassword"`
	PgDBName   string `json:"pgDBName"`
	PgPort     string `json:"pgPort"`
	MongoHost  string `json:"mongoHost"`
	MongoUser  string `json:"mongoUser"`
	MongoPassword string `json:"mongoPassword"`
	MongoDBName   string `json:"mongoDBName"`
	MongoPort     string `json:"mongoPort"`
	SchemaRegistryURL string `json:"schemaRegistryURL"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func HandleError(err error) {

	switch err.(type) {
	case viper.ConfigFileNotFoundError:
		log.Fatal("Config file not found")
		panic(err)
	case viper.ConfigParseError:
		log.Fatalf("Config file parse error : %s", err.Error())
		panic(err)
	default:
		panic(err)
	}
}
