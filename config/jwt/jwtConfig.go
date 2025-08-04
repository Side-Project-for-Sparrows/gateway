package jwt

import (
	"github.com/Side-Project-for-Sparrows/gateway/config"
)

type JwtConfig struct {
	PublicKeyUrl  string   `mapstructure:"publicKeyUrl"`
	ExcludedPaths []string `mapstructure:"excludedPaths"`
}

var Config JwtConfig

type jwtLoader struct{}

func (r *jwtLoader) Init() error {
	return config.NewYamlConfig("jwtConfig", "yaml", "./config/jwt").
		Decode(&Config)
}

func init() {
	config.RegisterConfig(&jwtLoader{})
}
