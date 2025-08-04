package config

type Configurable interface {
	Init() error
}
