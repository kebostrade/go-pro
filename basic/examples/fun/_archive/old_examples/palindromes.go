//go:build ignore

package main

import "fmt"

// Task: Write a function that takes a slice of strings and returns
// the count of strings that are palindromes (read the same forwards and backwards).
// The function should check strings concurrently using goroutines.

// Example input: ["radar", "hello", "deified", "golang"]
// Expected output: 2 (radar and deified are palindromes)

// Function signature to implement:
// func countPalindromes(words []string) int {
// Implement this
// }

func countPalindromes(words []string) int {
	// Create a channel to receive results from goroutines
	results := make(chan bool, len(words))

	// Launch a goroutine for each word
	for _, word := range words {
		go func(w string) {
			results <- isPalindrome(w)
		}(word)
	}

	// Count the palindromes
	count := 0
	for i := 0; i < len(words); i++ {
		if <-results {
			count++
		}
	}

	return count
}

// Helper function to check if a string is a palindrome
func isPalindrome(word string) bool {
	for i := 0; i < len(word)/2; i++ {
		if word[i] != word[len(word)-1-i] {
			return false
		}
	}
	return true
}

func main() {
	words := []string{"radar", "hello", "deified", "golang"}
	result := countPalindromes(words)
	fmt.Printf("Input words: %v\n", words) // Fixed: print words array instead of result
	fmt.Printf("Number of palindromes: %d\n", result)
}
