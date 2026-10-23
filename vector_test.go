package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVectorMath(t *testing.T) {
	a := Vector2D{3, 4}
	b := Vector2D{2, 1}
	c := a

	require.Equal(t, Vector2D{5, 5}, Add(a, b))
	require.Equal(t, Vector2D{8, 9}, Add(a, b, c))
	require.Equal(t, Vector2D{1, 3}, Subtract(a, b))
	require.Equal(t, Vector2D{6, 8}, Scale(a, 2))
	require.Equal(t, 10.0, DotProduct(a, b))
	require.Equal(t, 25.0, a.MagnitudeSquared())
	require.Equal(t, 5.0, a.Magnitude())

	div, err := Divide(a, 2)
	require.NoError(t, err)
	require.Equal(t, Vector2D{1.5, 2}, div)

	normalized, err := a.Normalize()
	require.NoError(t, err)
	require.InDelta(t, 0.6, normalized.X, 1e-9)
	require.InDelta(t, 0.8, normalized.Y, 1e-9)
}

func TestDivisionByZero(t *testing.T) {
	_, err := Divide(Vector2D{3, 4}, 0)
	require.ErrorIs(t, err, ErrDivisionByZero)
}

func TestNormalization(t *testing.T) {
	_, err := Vector2D{}.Normalize()
	require.ErrorIs(t, err, ErrDivisionByZero)
}
