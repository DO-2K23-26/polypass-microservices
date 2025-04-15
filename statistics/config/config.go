package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// Bootstrap is a flag to indicate if the application should start in bootstrap mode, meaning that the cycle should setup repositories e.g. for migrations or seeding
	Bootstrap bool `json:"bootstrap"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
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
