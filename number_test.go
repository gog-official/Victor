package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mustFraction(t *testing.T, numerator, denominator int64) Number {
	t.Helper()

	n, err := NewFraction(numerator, denominator)
	require.NoError(t, err)

	return n
}

func TestNewFraction(t *testing.T) {
	t.Run("integer fraction", func(t *testing.T) {
		n, err := NewFraction(6, 3)

		require.NoError(t, err)
		require.Equal(t, "2", n.String())
	})

	t.Run("will it reduce the fraction?", func(t *testing.T) {
		n, err := NewFraction(2, 4)

		require.NoError(t, err)

		require.Equal(t, "1/2", n.String())
	})
	t.Run("what about negative fraction", func(t *testing.T) {
		n, err := NewFraction(-2, 4)

		require.NoError(t, err)
		require.Equal(t, "-1/2", n.String())
	})
	t.Run("zero as denominator?", func(t *testing.T) {
		_, err := NewFraction(1, 0)
		require.ErrorIs(t, err, ErrZeroDenominator)
	})
}

func TestAddN(t *testing.T) {
	a := mustFraction(t, 1, 3)
	b := mustFraction(t, 1, 6)

	result := AddN(a, b)

	require.Equal(t, "1/2", result.String())
	require.Equal(t, "1/3", a.String())
	require.Equal(t, "1/6", b.String())
}

func TestSubN(t *testing.T) {
	a := mustFraction(t, 3, 4)
	b := mustFraction(t, 1, 4)

	result := SubN(a, b)

	require.Equal(t, "1/2", result.String())
	require.Equal(t, "3/4", a.String())
	require.Equal(t, "1/4", b.String())
}

func TestMulN(t *testing.T) {
	a := mustFraction(t, 2, 3)
	b := mustFraction(t, 3, 4)

	result := MulN(a, b)

	require.Equal(t, "1/2", result.String())
	require.Equal(t, "2/3", a.String())
	require.Equal(t, "3/4", b.String())
}

func TestDivN(t *testing.T) {
	t.Run("normal division", func(t *testing.T) {
		a := mustFraction(t, 1, 2)
		b := mustFraction(t, 1, 4)

		result, err := DivN(a, b)

		require.NoError(t, err)
		require.Equal(t, "2", result.String())
		require.Equal(t, "1/2", a.String())
		require.Equal(t, "1/4", b.String())
	})

	t.Run("division by zero", func(t *testing.T) {
		a := mustFraction(t, 1, 2)
		b := mustFraction(t, 0, 1)

		_, err := DivN(a, b)

		require.ErrorIs(t, err, ErrDivisionByZero)
	})
}
