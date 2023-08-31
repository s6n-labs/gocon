package gocon

import (
	"context"
	"reflect"
)

type bind[T, I any] struct{}

func Bind[T, I any]() Resolver[I] {
	return &bind[T, I]{}
}

func (r *bind[T, I]) Resolve(ctx context.Context, c Container) (I, error) {
	var zero I

	v, err := GetFrom[T](ctx, c)
	if err != nil {
		return zero, err
	}

	i, ok := reflect.ValueOf(v).Interface().(I)
	if !ok {
		return zero, ErrServiceNotFound
	}

	return i, nil
}
