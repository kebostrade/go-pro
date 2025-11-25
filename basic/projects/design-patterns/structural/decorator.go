package structural

import "fmt"

/*
DECORATOR PATTERN

Purpose: Attach additional responsibilities to an object dynamically.

Use Cases:
- Adding logging to functions
- Adding caching to data access
- Adding encryption to data
- HTTP middleware

Go-Specific Implementation:
- Function wrapping
- Interface composition
*/

// Coffee interface
type Coffee interface {
	Cost() float64
	Description() string
}

// SimpleCoffee is the base implementation
type SimpleCoffee struct{}

func (c *SimpleCoffee) Cost() float64 {
	return 2.0
}

func (c *SimpleCoffee) Description() string {
	return "Simple Coffee"
}

// MilkDecorator adds milk to coffee
type MilkDecorator struct {
	coffee Coffee
}

func NewMilkDecorator(c Coffee) *MilkDecorator {
	return &MilkDecorator{coffee: c}
}

func (m *MilkDecorator) Cost() float64 {
	return m.coffee.Cost() + 0.5
}

func (m *MilkDecorator) Description() string {
	return m.coffee.Description() + ", Milk"
}

// SugarDecorator adds sugar to coffee
type SugarDecorator struct {
	coffee Coffee
}

func NewSugarDecorator(c Coffee) *SugarDecorator {
	return &SugarDecorator{coffee: c}
}

func (s *SugarDecorator) Cost() float64 {
	return s.coffee.Cost() + 0.2
}

func (s *SugarDecorator) Description() string {
	return s.coffee.Description() + ", Sugar"
}

// WhipDecorator adds whipped cream
type WhipDecorator struct {
	coffee Coffee
}

func NewWhipDecorator(c Coffee) *WhipDecorator {
	return &WhipDecorator{coffee: c}
}

func (w *WhipDecorator) Cost() float64 {
	return w.coffee.Cost() + 0.7
}

func (w *WhipDecorator) Description() string {
	return w.coffee.Description() + ", Whipped Cream"
}

// Function Decorator Example
type DataService func(string) (string, error)

// LoggingDecorator adds logging to a function
func LoggingDecorator(fn DataService) DataService {
	return func(input string) (string, error) {
		fmt.Printf("📝 Calling function with input: %s\n", input)
		result, err := fn(input)
		if err != nil {
			fmt.Printf("❌ Function returned error: %v\n", err)
		} else {
			fmt.Printf("✅ Function returned: %s\n", result)
		}
		return result, err
	}
}

// CachingDecorator adds caching
func CachingDecorator(fn DataService) DataService {
	cache := make(map[string]string)
	return func(input string) (string, error) {
		if cached, ok := cache[input]; ok {
			fmt.Printf("💾 Cache hit for: %s\n", input)
			return cached, nil
		}
		fmt.Printf("🔍 Cache miss for: %s\n", input)
		result, err := fn(input)
		if err == nil {
			cache[input] = result
		}
		return result, err
	}
}

