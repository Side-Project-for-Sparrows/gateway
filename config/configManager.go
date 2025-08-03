package config

import (
	"fmt"
	"log"
)

var configLoaders []Configurable

func RegisterConfig(c Configurable) {
	configLoaders = append(configLoaders, c)
}

func InitAll() {
	fmt.Print("INIT TOTAL START")
	for _, loader := range configLoaders {
		if err := loader.Init(); err != nil {
			log.Fatalf("config init failed: %v", err)
		}
	}
}
