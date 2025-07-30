package route

import (
	"reflect"
	"strings"

	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type RouteConfig struct {
	User   string `mapstructure:"user"`
	Board  string `mapstructure:"board"`
	School string `mapstructure:"school"`
	Search string `mapstructure:"search"`
}

var (
	Config   RouteConfig
	RouteMap = make(map[string]string)
)

type routeLoader struct{}

func (r *routeLoader) Init() error {
	err := config.NewYamlConfig("routeConfig", "yaml", "./config/route").
		Decode(&Config)
	if err != nil {
		return err
	}
	// 맵 초기화
	RouteMap = make(map[string]string)

	val := reflect.ValueOf(Config)
	typ := reflect.TypeOf(Config)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("mapstructure")

		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		value := val.Field(i).String()
		RouteMap[tag] = value
	}
	return nil
}

func init() {
	config.RegisterConfig(&routeLoader{})
}
