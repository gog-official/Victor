package vector

import (
	"errors"
	"math"

	"github.com/gog-official/victor/internal/number"
)

var ErrDivisionByZero error = errors.New("cannot divide by 0")

type Vector2D struct {
	X number.Number
	Y number.Number
}

func (v Vector2D) Magnitude() float64 {
	return math.Sqrt(v.MagnitudeSquared().Float64())
}

func (v Vector2D) MagnitudeSquared() number.Number {
	return number.AddN(number.MulN(v.X, v.X), number.MulN(v.Y, v.Y))
}

func (v Vector2D) Normalize() (Vector2D, error) {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2D{}, ErrDivisionByZero
	}

	magNum := number.NewFloat64(mag)
	nX, err := number.DivN(v.X, magNum)
	if err != nil {
		return Vector2D{}, err
	}
	nY, err := number.DivN(v.Y, magNum)
	if err != nil {
		return Vector2D{}, err
	}

	return Vector2D{nX, nY}, nil
}

func Add(vectors ...Vector2D) Vector2D {
	sum := Vector2D{number.NewInteger(0), number.NewInteger(0)}

	for _, vector := range vectors {
		sum.X = number.AddN(sum.X, vector.X)
		sum.Y = number.AddN(sum.Y, vector.Y)
	}

	return sum
}

func Subtract(vec1, vec2 Vector2D) Vector2D {
	return Vector2D{
		X: number.SubN(vec1.X, vec2.X),
		Y: number.SubN(vec1.Y, vec2.Y),
	}
}

func Scale(vector Vector2D, scalar number.Number) Vector2D {
	return Vector2D{
		X: number.MulN(vector.X, scalar),
		Y: number.MulN(vector.Y, scalar),
	}
}

func Divide(vector Vector2D, scalar number.Number) (Vector2D, error) {
	if scalar.IsZero() {
		return Vector2D{}, ErrDivisionByZero
	}

	qX, err := number.DivN(vector.X, scalar)
	if err != nil {
		return Vector2D{}, err
	}
	qY, err := number.DivN(vector.Y, scalar)
	if err != nil {
		return Vector2D{}, err
	}

	return Vector2D{X: qX, Y: qY}, nil
}

func DotProduct(vec1, vec2 Vector2D) number.Number {
	return number.AddN(number.MulN(vec1.X, vec2.X), number.MulN(vec1.Y, vec2.Y))
}
