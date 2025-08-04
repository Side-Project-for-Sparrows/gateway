package circuitbreaker

import (
	"log"
	"math/rand"
	"time"
)

type CircuitStatusType int

const (
	Closed CircuitStatusType = iota
	Open
	HalfOpen
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type CircuitState struct {
	stateType CircuitStatusType
}

func (cs *CircuitState) Next(c *Circuit) CircuitStatusType {
	if cs.stateType == Closed {
		log.Printf("GRADE!!!!!!!!!!!!! %.4f", c.getGrade())

		log.Printf("STATUS!!!!!!!!!!!!! CLOSED")

		if c.getGrade() < float64(0.9) {
			return Open
		}
		return Closed
	}

	if cs.stateType == Open {
		log.Printf("GRADE!!!!!!!!!!!!! %.4f", c.getGrade())

		log.Printf("STATUS!!!!!!!!!!!!! OPEN")

		if c.getGrade() > float64(0.9) {
			return HalfOpen
		}
		return Open
	}

	if cs.stateType == HalfOpen {
		log.Printf("GRADE!!!!!!!!!!!!! %.4f", c.getGrade())

		log.Printf("STATUS!!!!!!!!!!!!! HALF")

		if c.getGrade() > float64(1) {
			return Closed
		}

		if c.getGrade() < float64(0.9) {
			return Open
		}

		return HalfOpen
	}

	return Open
}

func (cs *CircuitState) IsHealthy() bool {
	if cs.stateType == Closed {
		log.Printf("THIS IS CLOSE")
		return true
	}

	if cs.stateType == Open {
		log.Printf("THIS IS OPEN")
		return false
	}

	if cs.stateType == HalfOpen {
		log.Printf("THIS IS HALF")
		return rnd.Intn(2) == 0 // 50% 확률로 true 반환
	}

	return false
}
