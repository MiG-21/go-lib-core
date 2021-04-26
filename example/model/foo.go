package model

import (
	"fmt"
	go_lib_core "github.com/MiG-21/go-lib-core"
)

const (
	byteCap = 1024
)

// Foo representing some test struct
type Foo struct {
	go_lib_core.ReferenceCounter

	Stack []byte
}

// Reset method to reset event
func (e *Foo) Reset() {
	e.Stack = e.Stack[:0]
}

// ResetFoo method to reset Foo
// Used by reference countable pool
func ResetFoo(i interface{}) error {
	obj, ok := i.(*Foo)
	if !ok {
		return fmt.Errorf("illegal object sent to ResetFoo")
	}
	if len(obj.Stack) > byteCap {
		return fmt.Errorf("reset condition has been failed")
	}
	obj.Reset()
	return nil
}

// NewFoo method to create new event
func NewFoo() *Foo {
	return GetFoo()
}

// Foo pool
var fooPool = go_lib_core.NewReferenceCountedPool(
	func(counter go_lib_core.ReferenceCounter) go_lib_core.ReferenceCountable {
		br := new(Foo)
		br.ReferenceCounter = counter
		return br
	}, ResetFoo)

// GetFoo method to get new Foo
func GetFoo() *Foo {
	return fooPool.Get().(*Foo)
}

func GetStat() map[string]interface{} {
	return fooPool.Stats()
}
