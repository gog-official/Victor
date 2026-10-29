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

func TestNewInteger(t *testing.T) {
	t.Run("positive integer", func(t *testing.T) {
		n := NewInteger(9)
		require.Equal(t, "9", n.String())
	})
	t.Run("negative integer", func(t *testing.T) {
		n := NewInteger(-6)
		require.Equal(t, "-6", n.String())
	})
	t.Run("zero", func(t *testing.T) {
		n := NewInteger(0)
		require.Equal(t, "0", n.String())
	})
}

func TestIsZero(t *testing.T) {
	t.Run("integer zero", func(t *testing.T) {
		n := NewInteger(0)
		require.True(t, n.IsZero())
	})

	t.Run("zero fraction", func(t *testing.T) {
		n := mustFraction(t, 0, 5)
		require.True(t, n.IsZero())
	})

	t.Run("nonzero fraction", func(t *testing.T) {
		n := mustFraction(t, 1, 3)
		require.False(t, n.IsZero())
	})

	t.Run("negative number", func(t *testing.T) {
		n := NewInteger(-2)
		require.False(t, n.IsZero())
	})
}

func TestFloat64(t *testing.T) {
	t.Run("exact decimal conversion", func(t *testing.T) {
		n := mustFraction(t, 1, 4)
		require.Equal(t, 0.25, n.Float64())
	})

	t.Run("integer conversion", func(t *testing.T) {
		n := NewInteger(5)
		require.Equal(t, 5.0, n.Float64())
	})

	t.Run("repeating decimal conversion", func(t *testing.T) {
		n := mustFraction(t, 1, 3)
		require.InDelta(t, 0.3333333333333333, n.Float64(), 1e-15)
	})

	t.Run("negative conversion", func(t *testing.T) {
		n := mustFraction(t, -1, 2)
		require.Equal(t, -0.5, n.Float64())
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

func TestNewFloat64(t *testing.T) {
	t.Run("exact fraction", func(t *testing.T) {
		n := NewFloat64(0.5)
		require.Equal(t, "1/2", n.String())
	})

	t.Run("whole number", func(t *testing.T) {
		n := NewFloat64(2)
		require.Equal(t, "2", n.String())
	})

	t.Run("negative fraction", func(t *testing.T) {
		n := NewFloat64(-1.25)
		require.Equal(t, "-5/4", n.String())
	})
}
