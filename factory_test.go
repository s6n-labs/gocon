package gocon_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s6n-labs/gocon"
)

func TestFactory(t *testing.T) {
	t.Parallel()

	ctx := gocon.WithContainer(context.Background(), gocon.NewContainer(nil))
	require.NoError(t, gocon.Set(ctx, gocon.Value(Namer{name: "John"})))
	require.NoError(t, gocon.Set(ctx, gocon.Factory(func(c gocon.Container) (*NamedGreeter, error) {
		namer, err := gocon.GetFrom[Namer](ctx, c)
		if err != nil {
			return nil, err
		}

		return &NamedGreeter{
			Namer: namer,
		}, nil
	})))

	greeter, err := gocon.Get[*NamedGreeter](ctx)
	require.NoError(t, err)
	assert.Equal(t, "Hello, John!", greeter.Greet())
}
