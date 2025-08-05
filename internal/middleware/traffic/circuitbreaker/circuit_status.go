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
		log.Printf("CLOSED!!!!!!!!!!!!! %.6f", c.getGrade())

		if c.getGrade() < float64(0.0011) { // 최근 5회중 2회 이상 성공 못할시 Open 상태 전이 0,1
			return Open
		}
		return Closed
	}

	if cs.stateType == Open {
		log.Printf("OPEN!!!!!!!!!!!!! %.6f", c.getGrade())

		if c.getGrade() >= float64(0.0111) { // 최근 5회 요청중 3회이상 성공시 Half Open 상태 전환 3,4,5
			return HalfOpen
		}
		return Open
	}

	if cs.stateType == HalfOpen {
		log.Printf("HALF!!!!!!!!!!!!! %.6f", c.getGrade())

		if c.getGrade() >= float64(1.1111) { // 최근 5회요청 모두 성공하면 Close 상태 전환 5
			return Closed
		}

		if c.getGrade() < float64(0.0011) { // 최근 5회 요청중 2회 이상 성공하지 못했으면 Open 상태 전환 0,1
			return Open
		}

		return HalfOpen // 최근 요청 5회중 2,3,4 인경우 half open 유지
	}

	return Open
}

func (cs *CircuitState) IsHealthy() bool {
	if cs.stateType == Closed {
		return true
	}

	if cs.stateType == Open {
		return false
	}

	if cs.stateType == HalfOpen {
		return rnd.Intn(2) == 0 // 50% 확률로 true 반환
	}

	return false
}
