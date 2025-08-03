package circuitbreaker

type CircuitStatus int

const (
	CLOSED CircuitStatus = iota
	OPEN
	HALF_OPEN
)

func (s CircuitStatus) String() string {
	switch s {
	case CLOSED:
		return "CLOSED"
	case OPEN:
		return "OPEN"
	case HALF_OPEN:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}
