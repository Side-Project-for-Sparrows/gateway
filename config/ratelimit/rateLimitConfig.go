package ratelimit

import (
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type RateLimitConfig struct {
	RequestsPerSecond float64       `mapstructure:"requestsPerSecond"`
	WindowSize        time.Duration `mapstructure:"windowSize"`
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
