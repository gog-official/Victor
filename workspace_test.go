package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkspace(t *testing.T) {
	w := NewWorkspace()

	err := w.Define("      A    ", Vector2D{6, 9})
	require.NoError(t, err)

	v, exists := w.Get("a")
	require.True(t, exists)
	require.Equal(t, Vector2D{6, 9}, v)

	_, exists = w.Get("WOOOOOOOOOOOOOOOOOOOOOOOOO")
	require.False(t, exists)

	require.NoError(t, w.Define("c", Vector2D{6, 7}))
	require.NoError(t, w.Define("b", Vector2D{6.6, 9.9}))

	require.Equal(t, []string{"a", "b", "c"}, w.List())

	require.True(t, w.Delete("A"))
	_, exists = w.Get("a")
	require.False(t, exists)

	require.False(t, w.Delete("Alakazam"))

	require.NoError(t, w.Define("A", Vector2D{0, 0}))
	require.NoError(t, w.Define("A", Vector2D{6, 9}))
	v, exists = w.Get("A")
	require.True(t, exists)
	require.Equal(t, Vector2D{6, 9}, v)

	require.Equal(t, []string{"a", "b", "c"}, w.List())
}
