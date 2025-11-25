//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Solution: URL Shortener

type URLShortenerSolution struct {
	urlMap map[string]string // short code -> long URL
}

func NewURLShortenerSolution() *URLShortenerSolution {
	return &URLShortenerSolution{
		urlMap: make(map[string]string),
	}
}

func (us *URLShortenerSolution) Shorten(longURL string) string {
	// Check if URL already exists
	for short, long := range us.urlMap {
		if long == longURL {
			return short
		}
	}

	// Generate unique short code
	var shortCode string
	for {
		shortCode = generateShortCodeSolution(6)
		if _, exists := us.urlMap[shortCode]; !exists {
			break
		}
	}

	us.urlMap[shortCode] = longURL
	return shortCode
}

func (us *URLShortenerSolution) Expand(shortCode string) (string, bool) {
	longURL, exists := us.urlMap[shortCode]
	return longURL, exists
}

func generateShortCodeSolution(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	shortener := NewURLShortenerSolution()

	urls := []string{
		"https://www.example.com/very/long/url/path/to/resource",
		"https://github.com/golang/go",
		"https://go.dev/doc/",
	}

	fmt.Println("URL Shortener Solution")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	shortCodes := make([]string, 0)

	// Shorten URLs
	for _, url := range urls {
		short := shortener.Shorten(url)
		shortCodes = append(shortCodes, short)
		fmt.Printf("Long:  %s\n", url)
		fmt.Printf("Short: %s\n\n", short)
	}

	// Expand URLs
	fmt.Println("\nExpanding short codes:")
	fmt.Println(strings.Repeat("-", 60))
	for _, code := range shortCodes {
		if long, found := shortener.Expand(code); found {
			fmt.Printf("Code: %s\n", code)
			fmt.Printf("URL:  %s\n\n", long)
		}
	}
}
