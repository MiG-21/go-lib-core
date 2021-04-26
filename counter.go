package go_lib_core

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	// ReferenceCountable Interface following reference countable interface
	// We have provided inbuilt embeddable implementation of the reference countable pool
	// This interface just provides the extensibility for the implementation
	ReferenceCountable interface {
		// SetInstance Method to set the current instance
		SetInstance(i interface{})
		// IncrementReferenceCount Method to increment the reference count
		IncrementReferenceCount()
		// DecrementReferenceCount Method to decrement reference count
		DecrementReferenceCount() error
	}

	// ReferenceCounter Struct representing reference
	// This struct is supposed to be embedded inside the object to be pooled
	// Along with that incrementing and decrementing the references is highly important specifically around routines
	ReferenceCounter struct {
		Instance    interface{} `sql:"-" json:"-" yaml:"-"`
		destination *sync.Pool
		id          uint32
		count       *uint32
		released    *uint32
		reset       func(interface{}) error
	}
)

// IncrementReferenceCount Method to increment a reference
func (r ReferenceCounter) IncrementReferenceCount() {
	atomic.AddUint32(r.count, 1)
}

// DecrementReferenceCount Method to decrement a reference
// If the reference count goes to zero, the object is put back inside the pool
func (r ReferenceCounter) DecrementReferenceCount() error {
	if atomic.LoadUint32(r.count) == 0 {
		return fmt.Errorf("this should not happen => " + reflect.TypeOf(r.Instance).String())
	}
	if atomic.AddUint32(r.count, ^uint32(0)) == 0 {
		atomic.AddUint32(r.released, 1)
		if err := r.reset(r.Instance); err != nil {
			return fmt.Errorf("error while resetting an instance => " + err.Error())
		}
		r.destination.Put(r.Instance)
		r.Instance = nil
	}
	return nil
}

// SetInstance Method to set the current instance
func (r *ReferenceCounter) SetInstance(i interface{}) {
	r.Instance = i
}
