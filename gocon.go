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

func GetDefinitionFrom[T any](c Container) (*Definition, error) {
	def, err := c.get(keyOf(typeOf[T]()))
	if err != nil {
		return nil, err
	}

	return def, nil
}

func GetDefinition[T any](ctx context.Context) (*Definition, error) {
	c := FromContext(ctx)
	if c == nil {
		return nil, ErrNoContainer
	}

	return GetDefinitionFrom[T](c)
}

func GetFrom[T any](ctx context.Context, c Container) (T, error) {
	def, err := GetDefinitionFrom[T](c)
	if err != nil {
		var zero T
		return zero, err
	}

	return ResolveAs[T](ctx, c, def)
}

func Get[T any](ctx context.Context) (T, error) {
	c := FromContext(ctx)
	if c == nil {
		var zero T
		return zero, ErrNoContainer
	}

	return GetFrom[T](ctx, c)
}

func GetTaggedDefinitions(ctx context.Context, tag string) ([]*Definition, error) {
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

func GetTaggedFrom[I any](ctx context.Context, c Container, tag string) ([]I, error) {
	defs, err := GetTaggedDefinitions(ctx, tag)
	if err != nil {
		return nil, err
	}

	values := make([]I, 0, len(defs))
	for _, def := range defs {
		v, err := ResolveAs[I](ctx, c, def)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}

	return values, nil
}

func GetTagged[I any](ctx context.Context, tag string) ([]I, error) {
	c := FromContext(ctx)
	if c == nil {
		return nil, ErrNoContainer
	}

	return GetTaggedFrom[I](ctx, c, tag)
}

func Set(ctx context.Context, def *Definition) error {
	c := FromContext(ctx)
	if c == nil {
		return ErrNoContainer
	}

	return c.set(def)
}
