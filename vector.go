package main

import (
	"errors"
	"math"
)

var ErrDivisionByZero error = errors.New("cannot divide by 0")

type Vector2D struct {
	X float64
	Y float64
}

func (v Vector2D) Magnitude() float64 {
	mag := math.Sqrt(v.MagnitudeSquared())
	return mag
}

func (v Vector2D) MagnitudeSquared() float64 {
	magSq := v.X*v.X + v.Y*v.Y
	return magSq
}

func (v Vector2D) Normalize() (Vector2D, error) {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2D{}, ErrDivisionByZero
	}
	nX, nY := v.X/mag, v.Y/mag
	return Vector2D{
		nX, nY,
	}, nil
}

func Add(vectors ...Vector2D) Vector2D {
	sumX := 0.00
	sumY := 0.00

	for _, vector := range vectors {
		sumX += vector.X
		sumY += vector.Y
	}

	sumVector := Vector2D{
		X: sumX,
		Y: sumY,
	}

	return sumVector
}

func Subtract(vec1, vec2 Vector2D) Vector2D {
	difX := vec1.X - vec2.X
	difY := vec1.Y - vec2.Y
	difVector := Vector2D{
		X: difX,
		Y: difY,
	}

	return difVector
}

func Scale(vector Vector2D, scalar float64) Vector2D {
	scaledX := vector.X * scalar
	scaledY := vector.Y * scalar
	return Vector2D{
		X: scaledX,
		Y: scaledY,
	}
}

func Divide(vector Vector2D, scalar float64) (Vector2D, error) {
	if scalar == 0 {
		return Vector2D{}, ErrDivisionByZero
	}
	qX, qY := vector.X/scalar, vector.Y/scalar

	return Vector2D{
		X: qX,
		Y: qY,
	}, nil
}

func DotProduct(vec1, vec2 Vector2D) float64 {
	product := vec1.X*vec2.X + vec1.Y*vec2.Y
	return product
}
