package gocon

import (
	"context"
	"reflect"
)

func Bind[T, I any]() *Definition {
	var zero T

	rt := reflect.TypeOf(zero)

	return &Definition{
		Key:  keyOf(typeOf[I]()),
		Type: rt,
		resolveFunc: func(ctx context.Context, c Container) (*reflect.Value, error) {
			def, err := c.Get(keyOf(rt))
			if err != nil {
				return nil, err
			}

			return def.Resolve(ctx, c)
		},
	}
}
