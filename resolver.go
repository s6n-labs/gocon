package gocon

import (
	"context"
)

type Resolver[T any] interface {
	Resolve(ctx context.Context, c Container) (T, error)
}

type value[T any] struct {
	v T
}

func Value[T any](v T) Resolver[T] {
	return &value[T]{
		v: v,
	}
}

func (r *value[T]) Resolve(context.Context, Container) (T, error) {
	return r.v, nil
}
