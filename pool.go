package go_lib_core

import (
	"sync"
	"sync/atomic"
)

type (
	// Struct representing the pool
	referenceCountedPool struct {
		pool       *sync.Pool
		factory    func() ReferenceCountable
		returned   uint32
		allocated  uint32
		referenced uint32
	}
)

// NewReferenceCountedPool Method to create a new pool
func NewReferenceCountedPool(factory func(referenceCounter ReferenceCounter) ReferenceCountable,
	reset func(interface{}) error) *referenceCountedPool {
	p := new(referenceCountedPool)
	p.pool = new(sync.Pool)
	p.pool.New = func() interface{} {
		// Incrementing allocated count
		atomic.AddUint32(&p.allocated, 1)
		c := factory(ReferenceCounter{
			count:       new(uint32),
			destination: p.pool,
			released:    &p.returned,
			reset:       reset,
			id:          p.allocated,
		})
		return c
	}
	return p
}

// Get Method to get new object
func (p *referenceCountedPool) Get() ReferenceCountable {
	c := p.pool.Get().(ReferenceCountable)
	c.SetInstance(c)
	atomic.AddUint32(&p.referenced, 1)
	c.IncrementReferenceCount()
	return c
}

// Stats Method to return reference counted pool stats
func (p *referenceCountedPool) Stats() map[string]interface{} {
	return map[string]interface{}{
		"allocated":  p.allocated,
		"referenced": p.referenced,
		"returned":   p.returned,
	}
}
