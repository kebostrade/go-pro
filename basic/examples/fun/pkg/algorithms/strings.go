package algorithms

import (
	"strings"
	"unicode"
)

// IsPalindrome checks if a string is a palindrome (reads the same forwards and backwards)
// Ignores case and non-alphanumeric characters
// Time complexity: O(n)
func IsPalindrome(s string) bool {
	// Clean the string: remove non-alphanumeric and convert to lowercase
	cleaned := cleanString(s)

	left := 0
	right := len(cleaned) - 1

	for left < right {
		if cleaned[left] != cleaned[right] {
			return false
		}
		left++
		right--
	}

	return true
}

// cleanString removes non-alphanumeric characters and converts to lowercase
func cleanString(s string) string {
	var builder strings.Builder
	for _, ch := range s {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			builder.WriteRune(unicode.ToLower(ch))
		}
	}
	return builder.String()
}

// CountPalindromes counts how many strings in a slice are palindromes
// Time complexity: O(n * m) where n is number of strings and m is average length
func CountPalindromes(words []string) int {
	count := 0
	for _, word := range words {
		if IsPalindrome(word) {
			count++
		}
	}
	return count
}

// FindPalindromes returns all palindromes from a slice of strings
// Time complexity: O(n * m)
func FindPalindromes(words []string) []string {
	palindromes := make([]string, 0)
	for _, word := range words {
		if IsPalindrome(word) {
			palindromes = append(palindromes, word)
		}
	}
	return palindromes
}

// ReverseString reverses a string
// Time complexity: O(n)
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CountWords counts the number of words in a string
// Time complexity: O(n)
func CountWords(text string) int {
	return len(strings.Fields(text))
}

// WordFrequency counts the frequency of each word in a text
// Returns a map with words as keys and their counts as values
// Case-insensitive and ignores punctuation
// Time complexity: O(n)
func WordFrequency(text string) map[string]int {
	frequency := make(map[string]int)

	// Extract words
	words := extractWords(text)

	// Count each word
	for _, word := range words {
		if word != "" {
			frequency[word]++
		}
	}

	return frequency
}

// extractWords extracts words from text, converting to lowercase and removing punctuation
func extractWords(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Split into words
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	return words
}

// MostFrequentWord finds the most frequently occurring word in a text
// Returns the word and its count
// Time complexity: O(n)
func MostFrequentWord(text string) (string, int) {
	frequency := WordFrequency(text)

	maxWord := ""
	maxCount := 0

	for word, count := range frequency {
		if count > maxCount {
			maxWord = word
			maxCount = count
		}
	}

	return maxWord, maxCount
}

// IsAnagram checks if two strings are anagrams of each other
// Ignores case and spaces
// Time complexity: O(n)
func IsAnagram(s1, s2 string) bool {
	// Remove spaces and convert to lowercase
	s1 = strings.ReplaceAll(strings.ToLower(s1), " ", "")
	s2 = strings.ReplaceAll(strings.ToLower(s2), " ", "")

	if len(s1) != len(s2) {
		return false
	}

	// Count character frequencies
	charCount := make(map[rune]int)

	for _, ch := range s1 {
		charCount[ch]++
	}

	for _, ch := range s2 {
		charCount[ch]--
		if charCount[ch] < 0 {
			return false
		}
	}

	return true
}

// LongestCommonPrefix finds the longest common prefix among a slice of strings
// Time complexity: O(n * m) where n is number of strings and m is length of shortest string
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	if len(strs) == 1 {
		return strs[0]
	}

	// Use first string as reference
	prefix := strs[0]

	for i := 1; i < len(strs); i++ {
		// Reduce prefix until it matches the start of current string
		for !strings.HasPrefix(strs[i], prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}

	return prefix
}

// CountVowels counts the number of vowels in a string
// Time complexity: O(n)
func CountVowels(s string) int {
	count := 0
	vowels := "aeiouAEIOU"

	for _, ch := range s {
		if strings.ContainsRune(vowels, ch) {
			count++
		}
	}

	return count
}

// CountConsonants counts the number of consonants in a string
// Time complexity: O(n)
func CountConsonants(s string) int {
	count := 0

	for _, ch := range s {
		if unicode.IsLetter(ch) && !strings.ContainsRune("aeiouAEIOU", ch) {
			count++
		}
	}

	return count
}

// RemoveDuplicates removes duplicate characters from a string while preserving order
// Time complexity: O(n)
func RemoveDuplicates(s string) string {
	seen := make(map[rune]bool)
	var result strings.Builder

	for _, ch := range s {
		if !seen[ch] {
			seen[ch] = true
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// IsSubsequence checks if s is a subsequence of t
// Time complexity: O(n)
func IsSubsequence(s, t string) bool {
	if len(s) == 0 {
		return true
	}

	sIdx := 0
	for _, ch := range t {
		if sIdx < len(s) && rune(s[sIdx]) == ch {
			sIdx++
			if sIdx == len(s) {
				return true
			}
		}
	}

	return sIdx == len(s)
}

// LevenshteinDistance calculates the edit distance between two strings
// Time complexity: O(n * m)
// Space complexity: O(n * m)
func LevenshteinDistance(s1, s2 string) int {
	m, n := len(s1), len(s2)

	// Create a 2D slice for dynamic programming
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Initialize base cases
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// Fill the dp table
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}

	return dp[m][n]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// CapitalizeWords capitalizes the first letter of each word
// Time complexity: O(n)
func CapitalizeWords(s string) string {
	return strings.Title(strings.ToLower(s))
}

// TruncateString truncates a string to a maximum length and adds ellipsis if needed
// Time complexity: O(n)
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return s[:maxLen]
	}

	return s[:maxLen-3] + "..."
}
