package main

import (
	"errors"
	"math"
)

var ErrDivisionByZero error = errors.New("cannot divide by 0")

type Vector2D struct {
	X Number
	Y Number
}

func (v Vector2D) Magnitude() float64 {
	return math.Sqrt(v.MagnitudeSquared().Float64())
}

func (v Vector2D) MagnitudeSquared() Number {
	return AddN(MulN(v.X, v.X), MulN(v.Y, v.Y))
}

func (v Vector2D) Normalize() (Vector2D, error) {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2D{}, ErrDivisionByZero
	}

	magNum := NewFloat64(mag)
	nX, err := DivN(v.X, magNum)
	if err != nil {
		return Vector2D{}, err
	}
	nY, err := DivN(v.Y, magNum)
	if err != nil {
		return Vector2D{}, err
	}

	return Vector2D{nX, nY}, nil
}

func Add(vectors ...Vector2D) Vector2D {
	sum := Vector2D{NewInteger(0), NewInteger(0)}

	for _, vector := range vectors {
		sum.X = AddN(sum.X, vector.X)
		sum.Y = AddN(sum.Y, vector.Y)
	}

	return sum
}

func Subtract(vec1, vec2 Vector2D) Vector2D {
	return Vector2D{
		X: SubN(vec1.X, vec2.X),
		Y: SubN(vec1.Y, vec2.Y),
	}
}

func Scale(vector Vector2D, scalar Number) Vector2D {
	return Vector2D{
		X: MulN(vector.X, scalar),
		Y: MulN(vector.Y, scalar),
	}
}

func Divide(vector Vector2D, scalar Number) (Vector2D, error) {
	if scalar.IsZero() {
		return Vector2D{}, ErrDivisionByZero
	}

	qX, err := DivN(vector.X, scalar)
	if err != nil {
		return Vector2D{}, err
	}
	qY, err := DivN(vector.Y, scalar)
	if err != nil {
		return Vector2D{}, err
	}

	return Vector2D{X: qX, Y: qY}, nil
}

func DotProduct(vec1, vec2 Vector2D) Number {
	return AddN(MulN(vec1.X, vec2.X), MulN(vec1.Y, vec2.Y))
}
