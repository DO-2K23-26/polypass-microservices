package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    SchemaRegistryURL string
    KafkaHost         string
    ClientId          string
    HttpPort          string
    GrpcPort          string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("json")
    viper.AddConfigPath("./configs")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    return &Config{
        SchemaRegistryURL: viper.GetString("schemaRegistryURL"),
        KafkaHost:         viper.GetString("kafkaHost"),
        ClientId:          viper.GetString("clientId"),
        HttpPort:          viper.GetString("httpPort"),
        GrpcPort:          viper.GetString("grpcPort"),
    }, nil
}
