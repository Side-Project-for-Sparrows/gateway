package config

import (
	"log"

	"github.com/spf13/viper"
)

type RouteConfig struct {
	Board  string `mapstructure:"board"`
	School string `mapstructure:"school"`
	Search string `mapstructure:"search"`
}

type AppConfig struct {
	Env    string                 `mapstructure:"env"`
	Routes map[string]RouteConfig `mapstructure:"routes"`
}

var Conf AppConfig

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("config read error: %v", err)
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Fatalf("config unmarshal error: %v", err)
	}
}
