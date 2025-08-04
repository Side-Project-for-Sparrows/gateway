package ratelimit

import (
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type SlidingWindowConfig struct {
	RequestsPerSecond float64       `mapstructure:"requestsPerSecond"`
	WindowSize        time.Duration `mapstructure:"windowSize"`
	CleanInterval     time.Duration `mapstructure:"cleanInterval"`
}

type TokenBucketConfig struct {
	Capacity      int64         `mapstructure:"bucketCapacity"`
	RefillRate    int64         `mapstructure:"bucketRefillRate"`
	CleanInterval time.Duration `mapstructure:"cleanInterval"`
}

type RateLimitConfig struct {
	SlidingWindow SlidingWindowConfig `mapstructure:"slidingWindow"`
	TokenBucket   TokenBucketConfig   `mapstructure:"tokenBucket"`
}

var Config RateLimitConfig

type rateLimitLoader struct{}

func (r *rateLimitLoader) Init() error {
	err := config.NewYamlConfig("ratelimit", "yaml", "./config/ratelimit").
		Decode(&Config)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	config.RegisterConfig(&rateLimitLoader{})
}
