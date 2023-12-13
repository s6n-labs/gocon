package gocon

import (
	"context"
	"fmt"
	"reflect"
)

type Definition struct {
	Key     string
	Type    reflect.Type
	Tags    []string
	Resolve func(ctx context.Context, c Container) (*reflect.Value, error)
}

func (d *Definition) WithTags(tags ...string) *Definition {
	d.Tags = tags

	return d
}

func ResolveAs[T any](ctx context.Context, c Container, def *Definition) (T, error) {
	var zero T

	rv, err := def.Resolve(ctx, c)
	if err != nil {
		return zero, err
	}

	v, ok := rv.Interface().(T)
	if !ok {
		fmt.Printf("Type: %T\n", rv.Interface())

		return zero, ErrServiceNotFound
	}

	return v, nil
}
