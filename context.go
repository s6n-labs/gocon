package gocon

import (
	"context"
)

type contextKeyContainer struct{}

func WithContainer(ctx context.Context, container Container) context.Context {
	return context.WithValue(ctx, contextKeyContainer{}, container)
}

func FromContext(ctx context.Context) Container {
	c, ok := ctx.Value(contextKeyContainer{}).(Container)
	if !ok {
		return NewContainer(nil)
	}

	return c
}
