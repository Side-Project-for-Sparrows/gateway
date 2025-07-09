package config

import (
	"log"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	User   string `mapstructure:"user"`
	Board  string `mapstructure:"board"`
	School string `mapstructure:"school"`
	Search string `mapstructure:"search"`
}

type AppConfig struct {
	Env    string
	Routes map[string]RouteConfig `mapstructure:"routes"`
	JwtConfig map[string]JwtConfig `mapstructure:"jwt"`
}

type JwtConfig struct{
	PublicKeyUrl string `mapstructure:"publicKeyUrl"`
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
