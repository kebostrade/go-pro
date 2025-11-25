package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Exercise: URL Shortener
// Build a simple URL shortener service
// Requirements:
// 1. Generate short codes for long URLs
// 2. Store URL mappings
// 3. Retrieve original URLs from short codes
// 4. Handle collisions

type URLShortener struct {
	// TODO: Add fields to store URL mappings
	// Hint: Use a map[string]string to store shortCode -> longURL
}

// NewURLShortener creates a new URL shortener instance
func NewURLShortener() *URLShortener {
	// TODO: Initialize the shortener
	// Hint: Initialize the map
	return &URLShortener{}
}

// Shorten creates a short code for a long URL
func (us *URLShortener) Shorten(longURL string) string {
	// TODO: Implement URL shortening
	// 1. Generate a random short code (e.g., 6 characters)
	// 2. Check if it already exists (handle collisions)
	// 3. Store the mapping
	// 4. Return the short code
	return ""
}

// Expand retrieves the original URL from a short code
func (us *URLShortener) Expand(shortCode string) (string, bool) {
	// TODO: Implement URL expansion
	// Return the long URL and true if found, empty string and false otherwise
	return "", false
}

// generateShortCode generates a random short code
func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// TODO: Generate random code
	// Hint: Use rand.Intn() to pick random characters from charset
	return ""
}

func main() {
	rand.Seed(time.Now().UnixNano())

	shortener := NewURLShortener()

	// Test URLs
	urls := []string{
		"https://www.example.com/very/long/url/path/to/resource",
		"https://github.com/golang/go",
		"https://go.dev/doc/",
	}

	fmt.Println("URL Shortener Challenge")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Shorten URLs
	shortCodes := make([]string, 0)
	for _, url := range urls {
		short := shortener.Shorten(url)
		shortCodes = append(shortCodes, short)
		if short == "" {
			fmt.Printf("Long:  %s\n", url)
			fmt.Printf("Short: (not implemented)\n\n")
		} else {
			fmt.Printf("Long:  %s\n", url)
			fmt.Printf("Short: %s\n\n", short)
		}
	}

	// Expand URLs
	fmt.Println("\nExpanding short codes:")
	fmt.Println(strings.Repeat("-", 60))
	for _, code := range shortCodes {
		if code != "" {
			if longURL, found := shortener.Expand(code); found {
				fmt.Printf("Short: %s\n", code)
				fmt.Printf("Long:  %s\n\n", longURL)
			}
		}
	}
}
