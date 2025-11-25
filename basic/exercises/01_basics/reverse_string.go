//go:build ignore

package main

import (
	"fmt"
	"strings"
)

// Exercise: Reverse a String
// Write a function that reverses a string
// Handle Unicode characters correctly

// TODO: Implement the ReverseString function
func ReverseString(s string) string {
	// Your code here
	// Hint: Convert string to rune slice to handle Unicode properly
	// Step 1: Convert string to []rune
	// Step 2: Reverse the slice using two pointers
	// Step 3: Convert back to string
	return ""
}

func main() {
	testCases := []string{
		"Hello",
		"Go Programming",
		"12345",
		"Hello, 世界",
	}

	fmt.Println("String Reversal Challenge")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()

	for _, test := range testCases {
		reversed := ReverseString(test)
		fmt.Printf("Original: %s\n", test)
		if reversed == "" {
			fmt.Printf("Reversed: (not implemented)\n\n")
		} else {
			fmt.Printf("Reversed: %s\n\n", reversed)
		}
	}
}
