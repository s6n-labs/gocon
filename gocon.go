package gocon

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

var ErrServiceNotFound = errors.New("service does not exist, or cannot be resolved")

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func keyOf(rt reflect.Type) string {
	name := rt.String()
	path := rt.PkgPath()
	parts := strings.Split(path, "/")
	pkg := parts[len(parts)-1]

	if strings.HasPrefix(name, pkg+".") {
		return path + strings.TrimPrefix(name, pkg)
	}

	return path + "." + name
}

func KeyOf[T any]() string {
	return keyOf(typeOf[T]())
}

func Resolve[T any](ctx context.Context, c Container, key string) (T, error) {
	def, err := c.Get(key)
	if err != nil {
		var zero T
		return zero, err
	}

	return resolveAs[T](ctx, c, def)
}

func ResolveTagged[I any](ctx context.Context, c Container, tag string) ([]I, error) {
	defs, err := c.GetTagged(tag)
	if err != nil {
		return nil, err
	}

	values := make([]I, 0, len(defs))
	for _, def := range defs {
		v, err := resolveAs[I](ctx, c, def)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}

	return values, nil
}

func GetFrom[T any](ctx context.Context, c Container) (T, error) {
	return Resolve[T](ctx, c, KeyOf[T]())
}

func Get[T any](ctx context.Context) (T, error) {
	return GetFrom[T](ctx, FromContext(ctx))
}

func GetBy[T any](ctx context.Context, key string) (T, error) {
	return Resolve[T](ctx, FromContext(ctx), key)
}

func GetTagged[I any](ctx context.Context, tag string) ([]I, error) {
	return ResolveTagged[I](ctx, FromContext(ctx), tag)
}

func Set(ctx context.Context, def *Definition) error {
	return FromContext(ctx).Set(def)
}
