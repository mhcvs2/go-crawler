package middleware

import (
	"math"
	"sync"
)

type IdGenerator interface {
	GetUint32() uint32
}

func NewIdGenerator() IdGenerator {
	return &cyclicIdGenerator{}
}

type cyclicIdGenerator struct {
	sn    uint32
	ended bool
	mutex sync.Mutex
}

func (gen *cyclicIdGenerator) GetUint32() uint32 {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()
	if gen.ended {
		defer func() { gen.ended = false }()
		gen.sn = 0
		return gen.sn
	}
	id := gen.sn
	if id < math.MaxUint32 {
		gen.sn++
	} else {
		gen.ended = true
	}
	return id
}
