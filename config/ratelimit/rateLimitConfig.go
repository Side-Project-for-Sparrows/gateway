package ratelimit

import (
	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type RateLimitConfig struct {
	Limit int `mapstructure:"limit"`
}

var Config RateLimitConfig

type rateLimitLoader struct{}

func (r *rateLimitLoader) Init() error {
	return config.NewYamlConfig("ratelimit", "yaml", "./config/ratelimit").
		Decode(&Config)
}

func init() {
	config.RegisterConfig(&rateLimitLoader{})
}
