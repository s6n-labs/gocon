package gocon

import (
	"context"
	"reflect"
)

type Disposer interface {
	Dispose()
}

func dispose(rv reflect.Value) {
	if disposer, ok := rv.Interface().(Disposer); ok {
		disposer.Dispose()

		return
	}

	if rv.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < rv.NumField(); i++ {
		dispose(rv.Field(i))
	}
}

func Dispose[T any](service T) {
	dispose(reflect.ValueOf(service))
}

func DisposeAll(ctx context.Context) {

}
