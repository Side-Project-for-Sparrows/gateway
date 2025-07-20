package route

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type RouteConfig struct {
	User   string `mapstructure:"user"`
	Board  string `mapstructure:"board"`
	School string `mapstructure:"school"`
	Search string `mapstructure:"search"`
}

var Config RouteConfig

type routeLoader struct{}

func (r *routeLoader) Init(env string) error {
	viper.SetConfigName("routeConfig-" + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/route")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("route config read error: %w", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("route config unmarshal error: %w", err)
	}

	fmt.Printf("[DEBUG] Route config loaded: %+v\n", Config)
	return nil
}

func init() {
	config.RegisterConfig(&routeLoader{})
}
