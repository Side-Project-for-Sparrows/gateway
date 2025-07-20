package config

import (
	"log"
	"os"
)

var configLoaders []Configurable

func RegisterConfig(c Configurable) {
	configLoaders = append(configLoaders, c)
}

func InitAll() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // 기본값 fallback
		log.Println("ENV not set, using default: dev")
	}

	for _, loader := range configLoaders {
		if err := loader.Init(env); err != nil {
			log.Fatalf("config init failed: %v", err)
		}
	}
}
