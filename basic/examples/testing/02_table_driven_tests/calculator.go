package calculator

// Calculate performs basic arithmetic operations
func Calculate(a, b int, operation string) int {
	switch operation {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b == 0 {
			return 0
		}
		return a / b
	default:
		return 0
	}
}

// IsEven checks if a number is even
func IsEven(n int) bool {
	return n%2 == 0
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
