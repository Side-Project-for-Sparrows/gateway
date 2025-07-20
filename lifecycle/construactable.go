package lifecycle

type Constructable interface {
	Construct() error
}