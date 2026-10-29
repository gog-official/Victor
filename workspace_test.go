package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkspace(t *testing.T) {
	w := NewWorkspace()

	err := w.Define("      A    ", Vector2D{NewInteger(6), NewInteger(9)})
	require.NoError(t, err)

	v, exists := w.Get("a")
	require.True(t, exists)
	require.Equal(t, Vector2D{NewInteger(6), NewInteger(9)}, v)

	_, exists = w.Get("WOOOOOOOOOOOOOOOOOOOOOOOOO")
	require.False(t, exists)

	require.NoError(t, w.Define("c", Vector2D{NewInteger(6), NewInteger(7)}))
	require.NoError(t, w.Define("b", Vector2D{NewFloat64(6.6), NewFloat64(9.9)}))

	require.Equal(t, []string{"a", "b", "c"}, w.List())

	require.True(t, w.Delete("A"))
	_, exists = w.Get("a")
	require.False(t, exists)

	require.False(t, w.Delete("Alakazam"))

	require.NoError(t, w.Define("A", Vector2D{NewInteger(0), NewInteger(0)}))
	require.NoError(t, w.Define("A", Vector2D{NewInteger(6), NewInteger(9)}))
	v, exists = w.Get("A")
	require.True(t, exists)
	require.Equal(t, Vector2D{NewInteger(6), NewInteger(9)}, v)

	require.Equal(t, []string{"a", "b", "c"}, w.List())
}

func TestWorkspaceLimit(t *testing.T) {
	w := NewWorkspace()

	for i := 0; i < MaxVectors; i++ {
		require.NoError(t, w.Define(string(rune('a'+i)), Vector2D{NewInteger(int64(i)), NewInteger(0)}))
	}

	err := w.Define("f", Vector2D{NewInteger(5), NewInteger(5)})
	require.ErrorIs(t, err, ErrWorkspaceFull)

	require.NoError(t, w.Define("a", Vector2D{NewInteger(9), NewInteger(9)}))
	require.Len(t, w.List(), MaxVectors)
}
