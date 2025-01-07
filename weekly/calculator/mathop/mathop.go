package mathops

import "errors"

var DivisionByZero = errors.New("division by zero is not possible")

func Add(a, b int) float32 {
	return float32(a + b)
}

func Subtract(a, b int) float32 {
	return float32(a - b)
}

func Multiply(a, b int) float32 {
	return float32(a * b)
}

func Divide(a, b int) (float32, error) {
	if b == 0 {
		return 0, DivisionByZero
	}
	return float32(a) / float32(b), nil
}
