package gocon

import (
	"context"
	"reflect"
)

type Definition[T any] struct {
	Resolver Resolver[T]
	Tags     []string
}

type AnyDefinition interface {
	Resolve(ctx context.Context) (any, error)
	GetTags() []string

	asAnyDefinition() AnyDefinition
}

func (d *Definition[T]) reflectType() reflect.Type {
	var zero T
	return reflect.TypeOf(zero)
}

func (d *Definition[T]) asAnyDefinition() AnyDefinition {
	return d
}

func (d *Definition[T]) Resolve(ctx context.Context) (any, error) {
	return d.Resolver.Resolve(ctx)
}

func (d *Definition[T]) GetTags() []string {
	return d.Tags
}

func ResolveAs[I any](ctx context.Context, def AnyDefinition) (I, error) {
	var zero I

	anyValue, err := def.Resolve(ctx)
	if err != nil {
		return zero, err
	}

	v, ok := anyValue.(I)
	if !ok {
		return zero, ErrServiceNotFound
	}

	return v, nil
}