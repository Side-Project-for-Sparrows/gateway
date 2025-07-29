package route

import (
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

func (r *routeLoader) Init() error {
	return config.NewYamlConfig("routeConfig", "yaml", "./config/route").
		Decode(&Config)
}

func init() {
	config.RegisterConfig(&routeLoader{})
}
