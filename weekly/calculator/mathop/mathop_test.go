package mathops

import "testing"

func assertOperation(t *testing.T, got float32, want float32, err error, wantErr error) {
	t.Helper()
	if err != wantErr {
		t.Errorf("expected error %v, got %v", wantErr, err)
	}
	if got != want {
		t.Errorf("expected result %.2f, got %.2f", want, got)
	}
}

func TestOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation func(int, int) (float32, error)
		a, b      int
		want      float32
		wantErr   error
	}{
		{"Addition", func(a, b int) (float32, error) { return Add(a, b), nil }, 4, 2, 6, nil},
		{"Subtraction", func(a, b int) (float32, error) { return Subtract(a, b), nil }, 4, 2, 2, nil},
		{"Multiplication", func(a, b int) (float32, error) { return Multiply(a, b), nil }, 4, 2, 8, nil},
		{"Valid Division", Divide, 4, 2, 2, nil},
		{"Negative Division", Divide, 8, -4, -2, nil},
		{"Division by Zero", Divide, 8, 0, 0, DivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.operation(tt.a, tt.b)
			assertOperation(t, got, tt.want, err, tt.wantErr)
		})
	}
}
