package gocon_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s6n-labs/gocon"
)

type Namer struct {
	name string
}

type NamedGreeter struct {
	Namer Namer
}

func (g *NamedGreeter) Greet() string {
	return fmt.Sprintf("Hello, %s!", g.Namer.name)
}

func TestAutowire(t *testing.T) {
	container := gocon.NewContainer(nil)
	ctx := gocon.WithContainer(context.Background(), container)

	require.NoError(t, gocon.Set(ctx, gocon.Value(Namer{name: "John"})))

	greeter, err := gocon.Autowire[NamedGreeter]().Resolve(ctx, container)
	require.NoError(t, err)
	assert.Equal(t, "Hello, John!", greeter.Greet())
}
