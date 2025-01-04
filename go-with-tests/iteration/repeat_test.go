package iteration

import (
	"fmt"
	"testing"
)

func ExampleRepeat() {
	str := Repeat("a", 3)
	fmt.Println(str)
	// Output: aaa
}
func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 6)
	expected := "aaaaaa"

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}
