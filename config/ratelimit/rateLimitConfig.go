package ratelimit

import (
	"fmt"

	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/spf13/viper"
)

type RateLimitConfig struct {
	Limit int `mapstructure:"limit"`
}

var Config RateLimitConfig

type rateLimitLoader struct{}

func (r *rateLimitLoader) Init(env string) error {
	viper.SetConfigName("rateLimit-" + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/ratelimit")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("rateLimit config read error: %w", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("rateLimit config unmarshal error: %w", err)
	}

	fmt.Printf("[DEBUG] Rate Limit config loaded: %+v\n", Config)
	return nil
}

func init() {
	config.RegisterConfig(&rateLimitLoader{})
}
