package gocon

import (
	"context"
	"reflect"
)

func Value[T any](v T) *Definition {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	return &Definition{
		Key:  keyOf(rt),
		Type: rt,
		Resolve: func(ctx context.Context, c Container) (*reflect.Value, error) {
			return &rv, nil
		},
	}
}
