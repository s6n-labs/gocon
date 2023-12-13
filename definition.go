package gocon

import (
	"context"
	"reflect"
)

type Definition struct {
	Key   string
	Tags  []string
	Type  reflect.Type
	Value *reflect.Value

	resolveFunc func(ctx context.Context, c Container) (*reflect.Value, error)
}

func (d *Definition) WithTags(tags ...string) *Definition {
	d.Tags = tags

	return d
}

func (d *Definition) Resolve(ctx context.Context, c Container) (*reflect.Value, error) {
	if d.Value != nil {
		return d.Value, nil
	}

	rv, err := d.resolveFunc(ctx, c)
	if err != nil {
		return nil, err
	}

	d.Value = rv

	return rv, nil
}

func resolve[T any](ctx context.Context, c Container) (*reflect.Value, error) {
	def, err := c.Get(keyOf(typeOf[T]()))
	if err != nil {
		return nil, err
	}

	return def.Resolve(ctx, c)
}

func resolveAs[T any](ctx context.Context, c Container, def *Definition) (T, error) {
	var zero T

	rv, err := def.Resolve(ctx, c)
	if err != nil {
		return zero, err
	}

	v, ok := rv.Interface().(T)
	if !ok {
		return zero, ErrServiceNotFound
	}

	return v, nil
}
