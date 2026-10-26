package main

import (
	"errors"
	"math/big"
)

var ErrZeroDenominator error = errors.New("cannot use 0 as the denominator")

type Number struct {
	value big.Rat
}

func NewFraction(numerator, denominator int64) (Number, error) {
	if denominator == 0 {
		return Number{}, ErrZeroDenominator
	}
	num := Number{}
	num.value.SetFrac64(numerator, denominator)
	return num, nil
}

func NewInteger(value int64) Number {
	num := Number{}
	num.value.SetFrac64(value, 1)
	return num
}

func (n Number) String() string {
	return n.value.RatString()
}

func (n Number) IsZero() bool {
	return n.value.Sign() == 0
}

func (n Number) Float64() float64 {
	val, _ := n.value.Float64()
	return val
}

func AddN(a, b Number) Number {
	n := Number{}
	n.value.Add(&a.value, &b.value)
	return n
}

func SubN(a, b Number) Number {
	n := Number{}
	n.value.Sub(&a.value, &b.value)
	return n
}

func MulN(a, b Number) Number {
	n := Number{}
	n.value.Mul(&a.value, &b.value)
	return n
}

func DivN(a, b Number) (Number, error) {
	n := Number{}
	if b.value.Sign() == 0 {
		return Number{}, ErrDivisionByZero
	}
	n.value.Quo(&a.value, &b.value)
	return n, nil
}
