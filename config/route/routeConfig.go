package route

import (
	"github.com/Side-Project-for-Sparrows/gateway/config"
)

var (
	RouteMap = make(map[string]string)
)

type RouteLoader struct{}

func (r *RouteLoader) Init() error {
	err := config.NewYamlConfig("routeConfig", "yaml", "./config/route").
		Decode(&RouteMap)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	config.RegisterConfig(&RouteLoader{})
}
