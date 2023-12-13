package gocon_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s6n-labs/gocon"
)

func TestAlias(t *testing.T) {
	container := gocon.NewContainer(nil)
	ctx := gocon.WithContainer(context.Background(), container)

	require.NoError(t, gocon.Set(ctx, gocon.Value(&Greeter{message: "Hello, world!"})))
	require.NoError(t, gocon.Set(ctx, gocon.Alias[*Greeter]("my_greeter")))

	greeter, err := gocon.GetBy[*Greeter](ctx, "my_greeter")
	require.NoError(t, err)
	assert.Equal(t, "Hello, world!", greeter.Greet())
}
