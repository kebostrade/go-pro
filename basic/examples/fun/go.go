package main

import (
	"strconv"
	"strings"
)

// fizzbuzz takes a string and returns:
// - "3" if input is a positive number divisible by 3
// - "5" if input is a positive number divisible by 5
// - "15" if input is a positive number divisible by 15
// - the input string itself otherwise
func fizzbuzz(s string) string {
	// Try to parse the string as an integer
	num, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || num <= 0 {
		// Not a valid positive number
		return s
	}

	// Check divisibility (15 first to avoid false positives)
	if num%15 == 0 {
		return "15"
	}
	if num%3 == 0 {
		return "3"
	}
	if num%5 == 0 {
		return "5"
	}

	// Not divisible by 3, 5, or 15
	return s
}