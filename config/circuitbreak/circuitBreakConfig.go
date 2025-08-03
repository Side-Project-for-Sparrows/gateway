package circuitbreak

import (
	"fmt"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type CircuitBreakConfig struct {
	PingInterval time.Duration `mapstructure:"pingInterval"`
}

var Config CircuitBreakConfig

type circuitBreakLoader struct{}

func (r *circuitBreakLoader) Init() error {
	err := config.NewYamlConfig("circuitbreak", "yaml", "./config/circuitbreak").
		Decode(&Config)
	if err != nil {
		return err
	}

	fmt.Printf("PingInterval loaded = %.2f\n", Config.PingInterval.Seconds())
	return nil
}

func init() {
	config.RegisterConfig(&circuitBreakLoader{})
}
