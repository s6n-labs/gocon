package gocon_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siketyan/gocon"
)

type Greeter struct {
	message string
}

func (g Greeter) Greet() string {
	return g.message
}

type AnotherGreeter struct {
	name string
}

func (g AnotherGreeter) Greet() string {
	return fmt.Sprintf("Hello, %s!", g.name)
}

type GreeterLike interface {
	Greet() string
}

func Test_SimpleGetSet(t *testing.T) {
	container := gocon.NewContainer(nil)
	ctx := gocon.WithContainer(context.Background(), container)

	greeter := Greeter{message: "Hello, world!"}
	err := gocon.Set(ctx, gocon.Value(greeter))
	require.NoError(t, err)

	def, err := gocon.GetDefinition[Greeter](ctx)
	require.NoError(t, err)
	assert.Implements(t, new(gocon.Resolver[Greeter]), def.Resolver)
	assert.Len(t, def.Tags, 0)

	actual, err := gocon.Get[Greeter](ctx)
	require.NoError(t, err)
	assert.Equal(t, "Hello, world!", actual.message)
}

func Test_TaggedGetSet(t *testing.T) {
	container := gocon.NewContainer(nil)
	ctx := gocon.WithContainer(context.Background(), container)

	greeter1 := Greeter{message: "Hello, world!"}
	greeter2 := AnotherGreeter{name: "John"}

	require.NoError(t, gocon.Set(ctx, gocon.Value(greeter1), "greeter"))
	require.NoError(t, gocon.Set(ctx, gocon.Value(greeter2), "greeter"))

	defs, err := gocon.GetTaggedDefinitions(ctx, "greeter")
	require.NoError(t, err)
	assert.Len(t, defs, 2)

	greeters, err := gocon.GetTagged[GreeterLike](ctx, "greeter")
	require.NoError(t, err)
	assert.Len(t, greeters, 2)

	for _, g := range greeters {
		message := g.Greet()
		assert.True(t, message == "Hello, world!" || message == "Hello, John!")
	}
}