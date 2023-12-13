package gocon_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s6n-labs/gocon"
)

func TestDispose(t *testing.T) {
	t.Parallel()

	greeter := &Greeter{}

	gocon.Dispose(greeter)
	assert.True(t, greeter.disposed)
}

func TestContainer_DisposeAll(t *testing.T) {
	t.Parallel()

	greeter := &Greeter{}

	container := gocon.NewContainer(nil)
	require.NoError(t, container.Set(gocon.Value(greeter)))

	container.DisposeAll()
	assert.True(t, greeter.disposed)
}
