package gocon

import (
	"context"
	"errors"
	"reflect"
)

var ErrNotStruct = errors.New("autowire only supports structs")

func autowire(rt reflect.Type) *Definition {
	return &Definition{
		Key:  keyOf(rt),
		Type: rt,
		configureFunc: func(c Container) error {
			var configure func(rt reflect.Type) error

			configure = func(rt reflect.Type) error {
				if rt.Kind() == reflect.Pointer {
					return configure(rt.Elem())
				}

				if rt.Kind() != reflect.Struct {
					return ErrNotStruct
				}

				for i := 0; i < rt.NumField(); i++ {
					ft := rt.Field(i).Type

					if _, err := c.Get(keyOf(ft)); err == nil || !errors.As(err, new(ServiceNotFoundError)) {
						continue
					}

					for ft.Kind() == reflect.Pointer {
						ft = ft.Elem()
					}

					if ft.Kind() != reflect.Struct {
						continue
					}

					if err := c.Set(autowire(ft)); err != nil {
						return err
					}
				}

				return nil
			}

			return configure(rt)
		},
		resolveFunc: func(ctx context.Context, c Container) (*reflect.Value, error) {
			var resolve func(rt reflect.Type) (*reflect.Value, error)

			resolve = func(rt reflect.Type) (*reflect.Value, error) {
				if rt.Kind() == reflect.Pointer {
					rv, err := resolve(rt.Elem())
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

			return resolve(rt)
		},
	}
}

func Autowire[T any]() *Definition {
	return autowire(typeOf[T]())
}
