package config

type Constructable interface {
	Construct() error
}
