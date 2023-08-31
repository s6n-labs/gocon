package gocon

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

var (
	ErrNoContainer     = errors.New("no container in the context")
	ErrServiceNotFound = errors.New("service does not exist, or cannot be resolved")
)

func keyOfReflected(rt reflect.Type) string {
	name := rt.String()
	path := rt.PkgPath()
	parts := strings.Split(path, "/")
	pkg := parts[len(parts)-1]

	if strings.HasPrefix(name, pkg+".") {
		return path + strings.TrimPrefix(name, pkg)
	}

	return path + "." + name
}

func keyOf[T any]() string {
	var zero T
	return keyOfReflected(reflect.TypeOf(zero))
}

func GetDefinition[T any](ctx context.Context) (*Definition[T], error) {
	c := FromContext(ctx)
	if c == nil {
		return nil, ErrNoContainer
	}

	r, err := c.get(keyOf[T]())
	if err != nil {
		return nil, err
	}

	def, ok := r.(*Definition[T])
	if !ok {
		panic("BUG: type mismatch")
	}

	return def, nil
}

func Get[T any](ctx context.Context) (T, error) {
	def, err := GetDefinition[T](ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	return def.Resolver.Resolve(ctx)
}

func GetTaggedDefinitions(ctx context.Context, tag string) ([]AnyDefinition, error) {
	c := FromContext(ctx)
	if c == nil {
		return nil, ErrNoContainer
	}

	anyDefs, err := c.getTagged(tag)
	if err != nil {
		return nil, err
	}

	return anyDefs, nil
}

func GetTagged[I any](ctx context.Context, tag string) ([]I, error) {
	defs, err := GetTaggedDefinitions(ctx, tag)
	if err != nil {
		return nil, err
	}

	values := make([]I, 0, len(defs))
	for _, def := range defs {
		v, err := ResolveAs[I](ctx, def)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}

	return values, nil
}

func Set[T any](ctx context.Context, resolver Resolver[T], tags ...string) error {
	c := FromContext(ctx)
	if c == nil {
		return ErrNoContainer
	}

	return c.set(keyOf[T](), &Definition[T]{
		Resolver: resolver,
		Tags:     tags,
	})
}
