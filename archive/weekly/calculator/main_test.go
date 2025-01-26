package main

import "testing"

func assertDivision(t *testing.T, got float32, want float32, err error, wantErr error) {
	t.Helper()
	if err != wantErr {
		t.Errorf("expected error %v, got %v", wantErr, err)
	}
	if got != want {
		t.Errorf("expected result %v, got %v", want, got)
	}
}

func TestDivision(t *testing.T) {
	t.Run("Valid Division", func(t *testing.T) {
		got, err := Divide(4, 2)
		want := float32(2)
		assertDivision(t, got, want, err, nil)
	})

	t.Run("Negative Result", func(t *testing.T) {
		got, err := Divide(8, -4)
		want := float32(-2)
		assertDivision(t, got, want, err, nil)
	})

	t.Run("Division by Zero", func(t *testing.T) {
		got, err := Divide(8, 0)
		want := float32(0) // Default value for division by zero
		assertDivision(t, got, want, err, DivisionByZero)
	})
}
func TestMultiplication(t *testing.T) {
	t.Run("Multiplication", func(t *testing.T) {
		got := Multiply(4, 2)
		want := float32(8)
		assertOperation(t, got, want)
	})

	t.Run("Negative Result", func(t *testing.T) {
		got := Multiply(-2, 4)
		want := float32(-8)
		assertOperation(t, got, want)
	})
}
func TestSubtraction(t *testing.T) {
	t.Run("Subtraction", func(t *testing.T) {
		got := Subtract(4, 2)
		want := float32(2)
		assertOperation(t, got, want)
	})

	t.Run("Negative Result", func(t *testing.T) {
		got := Subtract(2, 4)
		want := float32(-2)
		assertOperation(t, got, want)
	})
}

func TestAdd(t *testing.T) {
	t.Run("Addition", func(t *testing.T) {
		got := Add(2, 4)
		want := float32(6)
		assertOperation(t, got, want)
	})

	t.Run("Negative Numbers", func(t *testing.T) {
		got := Add(-2, -4)
		want := float32(-6)
		assertOperation(t, got, want)
	})
}

func assertOperation(t testing.TB, got, want float32) {
	t.Helper()
	if got != want {
		t.Errorf("got %f want %f", got, want)
	}
}
