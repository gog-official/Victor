package vector

import (
	"fmt"
	"testing"

	"github.com/gog-official/victor/internal/number"
	"github.com/stretchr/testify/require"
)

func vec(x, y int64) Vector2D {
	return Vector2D{number.NewInteger(x), number.NewInteger(y)}
}

func mustFraction(t *testing.T, numerator, denominator int64) number.Number {
	t.Helper()

	n, err := number.NewFraction(numerator, denominator)
	require.NoError(t, err)

	return n
}

func formatVector(v Vector2D) string {
	return fmt.Sprintf("(%s, %s)", v.X.String(), v.Y.String())
}

func TestVectorMath(t *testing.T) {
	a := vec(3, 4)
	b := vec(2, 1)
	c := a

	require.Equal(t, vec(5, 5), Add(a, b))
	require.Equal(t, vec(8, 9), Add(a, b, c))
	require.Equal(t, vec(1, 3), Subtract(a, b))
	require.Equal(t, vec(6, 8), Scale(a, number.NewInteger(2)))
	require.Equal(t, "10", DotProduct(a, b).String())
	require.Equal(t, "25", a.MagnitudeSquared().String())
	require.Equal(t, 5.0, a.Magnitude())

	div, err := Divide(a, number.NewInteger(2))
	require.NoError(t, err)
	require.Equal(t, "3/2", div.X.String())
	require.Equal(t, "2", div.Y.String())

	normalized, err := a.Normalize()
	require.NoError(t, err)
	require.Equal(t, "3/5", normalized.X.String())
	require.Equal(t, "4/5", normalized.Y.String())
}

func TestDivisionByZero(t *testing.T) {
	_, err := Divide(vec(3, 4), number.NewInteger(0))
	require.ErrorIs(t, err, ErrDivisionByZero)
}

func TestNormalization(t *testing.T) {
	_, err := Vector2D{}.Normalize()
	require.ErrorIs(t, err, ErrDivisionByZero)
}

func TestVectorRationalMath(t *testing.T) {
	a := Vector2D{mustFraction(t, 1, 2), mustFraction(t, 3, 4)}
	b := Vector2D{mustFraction(t, 1, 3), mustFraction(t, 1, 4)}

	require.Equal(t, "(5/6, 1)", formatVector(Add(a, b)))
	require.Equal(t, "(1/6, 1/2)", formatVector(Subtract(a, b)))
	require.Equal(t, "(1, 3/2)", formatVector(Scale(a, number.NewInteger(2))))
	require.Equal(t, "17/48", DotProduct(a, b).String())

	div, err := Divide(a, number.NewInteger(2))
	require.NoError(t, err)
	require.Equal(t, "(1/4, 3/8)", formatVector(div))
}

func TestVectorDivisionByZero(t *testing.T) {
	_, err := Divide(Vector2D{number.NewInteger(1), number.NewInteger(2)}, number.NewInteger(0))
	require.ErrorIs(t, err, ErrDivisionByZero)
}
