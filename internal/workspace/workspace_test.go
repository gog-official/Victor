package workspace

import (
	"testing"

	"github.com/gog-official/victor/internal/number"
	"github.com/gog-official/victor/internal/vector"
	"github.com/stretchr/testify/require"
)

func vec(x, y int64) vector.Vector2D {
	return vector.Vector2D{X: number.NewInteger(x), Y: number.NewInteger(y)}
}

func TestWorkspace(t *testing.T) {
	w := NewWorkspace()

	err := w.Define("      A    ", vec(6, 9))
	require.NoError(t, err)

	v, exists := w.Get("a")
	require.True(t, exists)
	require.Equal(t, vec(6, 9), v)

	_, exists = w.Get("WOOOOOOOOOOOOOOOOOOOOOOOOO")
	require.False(t, exists)

	require.NoError(t, w.Define("c", vec(6, 7)))
	require.NoError(t, w.Define("b", vector.Vector2D{
		X: number.NewFloat64(6.6),
		Y: number.NewFloat64(9.9),
	}))

	require.Equal(t, []string{"a", "b", "c"}, w.List())

	require.True(t, w.Delete("A"))
	_, exists = w.Get("a")
	require.False(t, exists)

	require.False(t, w.Delete("Alakazam"))

	require.NoError(t, w.Define("A", vec(0, 0)))
	require.NoError(t, w.Define("A", vec(6, 9)))
	v, exists = w.Get("A")
	require.True(t, exists)
	require.Equal(t, vec(6, 9), v)

	require.Equal(t, []string{"a", "b", "c"}, w.List())
}

func TestWorkspaceLimit(t *testing.T) {
	w := NewWorkspace()

	for i := 0; i < MaxVectors; i++ {
		require.NoError(t, w.Define(string(rune('a'+i)), vec(int64(i), 0)))
	}

	err := w.Define("f", vec(5, 5))
	require.ErrorIs(t, err, ErrWorkspaceFull)

	require.NoError(t, w.Define("a", vec(9, 9)))
	require.Len(t, w.List(), MaxVectors)
}
