//go:build ignore

package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Task: Write a function that counts the frequency of each word in a text.
// The function should be case-insensitive and ignore punctuation.
// Return a map with words as keys and their counts as values.

// Example input: "Hello world! Hello Go. Go is awesome."
// Expected output: map[hello:2 world:1 go:2 is:1 awesome:1]

// Function signature to implement:
// func countWords(text string) map[string]int

func countWords(text string) map[string]int {
	wordCount := make(map[string]int)

	// Convert to lowercase and split into words
	words := extractWords(text)

	// Count each word
	for _, word := range words {
		if word != "" {
			wordCount[word]++
		}
	}

	return wordCount
}

// extractWords extracts words from text, removing punctuation
func extractWords(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Build words by filtering out non-letter characters
	var words []string
	var currentWord strings.Builder

	for _, char := range text {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			currentWord.WriteRune(char)
		} else {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}
	}

	// Add the last word if exists
	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// findMostFrequent returns the most frequent word and its count
func findMostFrequent(wordCount map[string]int) (string, int) {
	var maxWord string
	maxCount := 0

	for word, count := range wordCount {
		if count > maxCount {
			maxWord = word
			maxCount = count
		}
	}

	return maxWord, maxCount
}

// printWordFrequency prints words and their frequencies in a formatted way
func printWordFrequency(wordCount map[string]int) {
	fmt.Println("\nWord Frequency:")
	fmt.Println(strings.Repeat("-", 40))
	for word, count := range wordCount {
		fmt.Printf("%-20s: %d\n", word, count)
	}
}

func main() {
	// Test case 1
	text1 := "Hello world! Hello Go. Go is awesome."
	fmt.Println("Text:", text1)
	wordCount1 := countWords(text1)
	fmt.Println("Word count:", wordCount1)

	mostFrequent, count := findMostFrequent(wordCount1)
	fmt.Printf("Most frequent word: '%s' (appears %d times)\n", mostFrequent, count)

	// Test case 2
	text2 := `Go is an open source programming language that makes it easy to build 
	simple, reliable, and efficient software. Go is expressive, concise, clean, 
	and efficient. Go is great for building web applications.`

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Analyzing longer text...")
	wordCount2 := countWords(text2)
	printWordFrequency(wordCount2)

	mostFrequent2, count2 := findMostFrequent(wordCount2)
	fmt.Printf("\nMost frequent word: '%s' (appears %d times)\n", mostFrequent2, count2)
}
