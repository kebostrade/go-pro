//go:build ignore

package main

import (
	"fmt"
	"strings"
)

// Solution: Reverse a String

func ReverseStringSolution(s string) string {
	// Convert to rune slice to handle Unicode correctly
	runes := []rune(s)

	// Reverse the rune slice using two pointers
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func main() {
	testCases := []string{
		"Hello",
		"Go Programming",
		"12345",
		"Hello, 世界",
	}

	fmt.Println("String Reversal Solution")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()

	for _, test := range testCases {
		reversed := ReverseStringSolution(test)
		fmt.Printf("Original: %s\n", test)
		fmt.Printf("Reversed: %s\n\n", reversed)
	}
}
