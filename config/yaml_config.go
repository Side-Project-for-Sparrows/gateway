package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type yamlConfig struct {
	key  string
	typ  string
	path string
}

func (y *yamlConfig) Key() string {
	return y.key
}

func (y *yamlConfig) Decode(target interface{}) error {
	env := os.Getenv("ENV")
	viper.SetConfigName(y.key + "-" + env)
	viper.SetConfigType(y.typ)
	viper.AddConfigPath(y.path)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read error: %w", err)
	}
	if err := viper.Unmarshal(target); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	return nil
}

func NewYamlConfig(key string, typ string, path string) *yamlConfig {
	return &yamlConfig{
		key:  key,
		typ:  typ,
		path: path,
	}
}
