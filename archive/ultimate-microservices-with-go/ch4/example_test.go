package main

import (
	"testing"
)

func TestGetBigger(t *testing.T) {
	type args struct {
		a float64
		b float64
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "a is the bigger one",
			args: args{
				a: 9,
				b: 6,
			},
			want: 9,
		},
		{
			name: "b is the bigger one",
			args: args{
				a: 3,
				b: 7,
			},
			want: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBigger(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("GetBigger() = %v,want %v", got, tt.want)
			}
		})
	}
}

