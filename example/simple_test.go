package example_test

import (
	"testing"

	"github.com/MiG-21/go-lib-core/example/model"
)

//go:noinline
func inc(foo *model.Foo) { foo.Stack = append(foo.Stack, '1') }

func BenchmarkWithoutPool(b *testing.B) {
	var foo *model.Foo
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			foo = &model.Foo{}
			b.StopTimer()
			inc(foo)
			b.StartTimer()
		}
	}
}

func BenchmarkWithPool(b *testing.B) {
	var foo *model.Foo
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			foo = model.GetFoo()
			b.StopTimer()
			inc(foo)
			b.StartTimer()
			_ = foo.DecrementReferenceCount()
		}
	}
}
