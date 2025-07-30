package lifecycle

import (
	"log"
	"os"
)

var constructables []Constructable

func Register(c Constructable) {
	constructables = append(constructables, c)
}

func ConstructAll() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
		log.Println("[Construct] ENV not set, using default: dev")
	}

	for _, c := range constructables {
		if err := c.Construct(); err != nil {
			log.Fatalf("construct failed: %v", err)
		}
	}
}
