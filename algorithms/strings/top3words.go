package strings

import (
	"regexp"
	"sort"
	"strings"
)

// WordCount represents a word and its frequency count
type WordCount struct {
	Word  string
	Count int
}

// Top3Words finds the 3 most frequently occurring words in text.
// Words are case-insensitive and punctuation is ignored.
func Top3Words(text string) []WordCount {
	// Convert to lowercase for case-insensitive matching
	text = strings.ToLower(text)

	// Extract words (sequences of letters only)
	re := regexp.MustCompile(`[a-z]+`)
	words := re.FindAllString(text, -1)

	// Count word frequencies
	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}

	// Convert to slice for sorting
	var counts []WordCount
	for word, count := range wordCounts {
		counts = append(counts, WordCount{word, count})
	}

	// Sort by count descending, then by word ascending for ties
	sort.Slice(counts, func(i, j int) bool {
		if counts[i].Count != counts[j].Count {
			return counts[i].Count > counts[j].Count
		}
		return counts[i].Word < counts[j].Word
	})

	// Return top 3 (or fewer if text has less than 3 unique words)
	if len(counts) < 3 {
		return counts
	}
	return counts[:3]
}
