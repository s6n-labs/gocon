package gocon

import (
	"context"
	"errors"
	"reflect"
)

var ErrNotStruct = errors.New("autowire only supports structs")

func Autowire[T any]() *Definition {
	var zero T

	rt := reflect.TypeOf(zero)

	return &Definition{
		Key:  keyOf(rt),
		Type: rt,
		resolveFunc: func(ctx context.Context, c Container) (*reflect.Value, error) {
			return resolveAutowire(ctx, c, rt)
		},
	}
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

		def, err := c.Get(keyOf(field.Type))
		if err != nil {
			if errors.Is(err, ErrServiceNotFound) {
				return resolveAutowire(ctx, c, field.Type)
			}

			return nil, err
		}

		v, err := def.Resolve(ctx, c)
		if err != nil {
			return nil, err
		}

		rv.Field(i).Set(*v)
	}

	return &rv, nil
}
