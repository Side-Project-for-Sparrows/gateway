package config

type Configurable interface {
	Init(env string) error
}