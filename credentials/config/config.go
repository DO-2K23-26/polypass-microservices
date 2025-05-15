package config

import (
	"fmt"

	"github.com/DO-2K23-26/polypass-microservices/credentials/application/http"
	"github.com/DO-2K23-26/polypass-microservices/credentials/infrastructure/sql"
	"github.com/optique-dev/core"

	"github.com/spf13/viper"
)

type Config struct {
	// Bootstrap is a flag to indicate if the application should start in bootstrap mode, meaning that the cycle should setup repositories e.g. for migrations or seeding
	Bootstrap bool        `json:"bootstrap"`
	Database  sql.Config  `json:"database"`
	Server    http.Config `json:"server"`
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
		core.Error("Config file not found")
		panic(err)
	case viper.ConfigParseError:
		core.Error(fmt.Sprintf("Config file parse error : %s", err.Error()))
		panic(err)
	default:
		panic(err)
	}
}
