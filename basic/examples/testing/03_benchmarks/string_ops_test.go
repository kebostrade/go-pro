package stringops

import "testing"

// Exercise: Write benchmark tests
// Learn how to measure performance using benchmarks
// Run with: go test -bench=.

var testStrings = []string{"Hello", " ", "World", " ", "from", " ", "Go", "!"}

func BenchmarkConcatWithPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatWithPlus(testStrings)
	}
}

func BenchmarkConcatWithBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatWithBuilder(testStrings)
	}
}

func BenchmarkConcatWithJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatWithJoin(testStrings)
	}
}

// Test correctness
func TestConcatMethods(t *testing.T) {
	expected := "Hello World from Go!"

	tests := []struct {
		name string
		fn   func([]string) string
	}{
		{"Plus", ConcatWithPlus},
		{"Builder", ConcatWithBuilder},
		{"Join", ConcatWithJoin},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(testStrings)
			if result != expected {
				t.Errorf("%s = %s; want %s", tt.name, result, expected)
			}
		})
	}
}
