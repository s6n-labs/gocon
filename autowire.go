package gocon

import (
	"context"
	"errors"
	"reflect"
)

var ErrNotStruct = errors.New("autowire only supports structs")

type autowire[T any] struct{}

func Autowire[T any]() Resolver[T] {
	return &autowire[T]{}
}

func (r *autowire[T]) Resolve(ctx context.Context, c Container) (T, error) {
	var zero T

	rv, err := resolveAutowire(ctx, c, reflect.TypeOf(zero))
	if err != nil {
		return zero, err
	}

	v, ok := rv.Interface().(T)
	if !ok {
		panic("BUG: type mismatch")
	}

	return v, nil
}

func resolveAutowire(ctx context.Context, c Container, rt reflect.Type) (*reflect.Value, error) {
	if rt.Kind() == reflect.Pointer {
		rv, err := resolveAutowire(ctx, c, rt.Elem())
		if err != nil {
			return nil, err
		}

		a := rv.Addr()

		return &a, nil
	}

	if rt.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	rv := reflect.New(rt).Elem()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		def, err := c.get(keyOfReflected(field.Type))
		if err != nil {
			return nil, err
		}

		v, err := def.Resolve(ctx, c)
		if err != nil {
			return nil, err
		}

		rv.Field(i).Set(reflect.ValueOf(v))
	}

	return &rv, nil
}
