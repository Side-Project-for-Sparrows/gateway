package jwt

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type JwtConfig struct {
	PublicKeyUrl string `mapstructure:"publicKeyUrl"`
}

var Config JwtConfig

type jwtLoader struct{}

func (j *jwtLoader) Init(env string) error {
	viper.SetConfigName("jwtConfig-" + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/jwt")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("jwt config read error: %w", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("jwt config unmarshal error: %w", err)
	}
	return nil
}

func init() {
	config.RegisterConfig(&jwtLoader{})
}
