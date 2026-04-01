package greeting

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreet(t *testing.T) {
	t.Run("greet with name", func(t *testing.T) {
		result := Greet("Alice", 1)
		assert.Equal(t, "Hello, Alice!\n", result)
	})

	t.Run("greet multiple times", func(t *testing.T) {
		result := Greet("Bob", 3)
		expected := "Hello, Bob!\n" + "Hello, Bob!\n" + "Hello, Bob!\n"
		assert.Equal(t, expected, result)
	})

	t.Run("greet with empty name uses World", func(t *testing.T) {
		result := Greet("", 1)
		assert.Equal(t, "Hello, World!\n", result)
	})

	t.Run("greet with zero times", func(t *testing.T) {
		result := Greet("Test", 0)
		assert.Equal(t, "Hello, Test!\n", result)
	})

	t.Run("greet with negative times", func(t *testing.T) {
		result := Greet("Test", -5)
		assert.Equal(t, "Hello, Test!\n", result)
	})
}

func TestGetGreetingTemplate(t *testing.T) {
	template := GetGreetingTemplate()
	assert.Equal(t, "Hello, %s!\n", template)
	assert.Contains(t, template, "%s")
}

func TestGreetWithLock(t *testing.T) {
	svc := NewGreetingService()

	t.Run("thread-safe greet", func(t *testing.T) {
		result := GreetWithLock(svc, "ThreadSafe", 2)
		assert.Equal(t, "Hello, ThreadSafe!\n"+"Hello, ThreadSafe!\n", result)
	})
}

func TestGreetContainsName(t *testing.T) {
	testCases := []struct {
		name  string
		times int
	}{
		{"Alice", 1},
		{"Bob", 2},
		{"Charlie", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Greet(tc.name, tc.times)
			assert.True(t, strings.Contains(result, tc.name))
			assert.True(t, strings.Contains(result, "Hello"))
			assert.True(t, strings.Contains(result, "!"))
		})
	}
}
