package gocon_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s6n-labs/gocon"
)

func TestBind(t *testing.T) {
	container := gocon.NewContainer(nil)
	ctx := gocon.WithContainer(context.Background(), container)

	require.NoError(t, gocon.Set(ctx, gocon.Value(Greeter{message: "Hello, world!"})))
	require.NoError(t, gocon.Set(ctx, gocon.Bind[Greeter, GreeterLike]()))

	greeter, err := gocon.Get[GreeterLike](ctx)
	require.NoError(t, err)
	assert.Equal(t, "Hello, world!", greeter.Greet())
}
