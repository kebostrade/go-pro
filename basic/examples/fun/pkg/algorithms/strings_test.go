package algorithms

import (
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"simple palindrome", "racecar", true},
		{"with spaces", "A man a plan a canal Panama", true},
		{"with punctuation", "A man, a plan, a canal: Panama", true},
		{"not palindrome", "hello", false},
		{"single char", "a", true},
		{"empty", "", true},
		{"numbers", "12321", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(tt.s)
			if got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"simple", "hello", "olleh"},
		{"single char", "a", "a"},
		{"empty", "", ""},
		{"with unicode", "Hello, 世界", "界世 ,olleH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReverseString(tt.s)
			if got != tt.want {
				t.Errorf("ReverseString(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{"simple sentence", "Hello world", 2},
		{"multiple spaces", "Hello    world", 2},
		{"empty", "", 0},
		{"single word", "Hello", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountWords(tt.text)
			if got != tt.want {
				t.Errorf("CountWords(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}

func TestWordFrequency(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog the"
	freq := WordFrequency(text)

	if freq["the"] != 3 {
		t.Errorf("WordFrequency: 'the' count = %d, want 3", freq["the"])
	}
	if freq["quick"] != 1 {
		t.Errorf("WordFrequency: 'quick' count = %d, want 1", freq["quick"])
	}
}

func TestMostFrequentWord(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog the"
	word, count := MostFrequentWord(text)

	if word != "the" || count != 3 {
		t.Errorf("MostFrequentWord() = (%q, %d), want (\"the\", 3)", word, count)
	}
}

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{"simple anagrams", "listen", "silent", true},
		{"with spaces", "conversation", "voices rant on", true},
		{"not anagrams", "hello", "world", false},
		{"empty strings", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAnagram(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("IsAnagram(%q, %q) = %v, want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want string
	}{
		{"common prefix", []string{"flower", "flow", "flight"}, "fl"},
		{"no common", []string{"dog", "racecar", "car"}, ""},
		{"empty array", []string{}, ""},
		{"single string", []string{"hello"}, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestCommonPrefix(tt.strs)
			if got != tt.want {
				t.Errorf("LongestCommonPrefix() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCountVowels(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"simple", "hello", 2},
		{"all vowels", "aeiou", 5},
		{"no vowels", "bcdfg", 0},
		{"mixed case", "AEIOUaeiou", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountVowels(tt.s)
			if got != tt.want {
				t.Errorf("CountVowels(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestCountConsonants(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"simple", "hello", 3},
		{"no consonants", "aeiou", 0},
		{"all consonants", "bcdfg", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountConsonants(tt.s)
			if got != tt.want {
				t.Errorf("CountConsonants(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"with duplicates", "hello", "helo"},
		{"no duplicates", "world", "world"},
		{"all same", "aaa", "a"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDuplicates(tt.s)
			if got != tt.want {
				t.Errorf("RemoveDuplicates(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsSubsequence(t *testing.T) {
	tests := []struct {
		name string
		s    string
		t    string
		want bool
	}{
		{"is subsequence", "abc", "ahbgdc", true},
		{"not subsequence", "axc", "ahbgdc", false},
		{"empty s", "", "hello", true},
		{"equal strings", "hello", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSubsequence(tt.s, tt.t)
			if got != tt.want {
				t.Errorf("IsSubsequence(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
		})
	}
}

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"same strings", "hello", "hello", 0},
		{"one char diff", "hello", "hallo", 1},
		{"different", "kitten", "sitting", 3},
		{"empty strings", "", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LevenshteinDistance(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		maxLen int
		want   string
	}{
		{"no truncation", "hello", 10, "hello"},
		{"truncate", "hello world", 8, "hello..."},
		{"exactly max", "hello", 5, "hello"},
		{"very short", "hello", 2, "he"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateString(tt.s, tt.maxLen)
			if got != tt.want {
				t.Errorf("TruncateString(%q, %d) = %q, want %q", tt.s, tt.maxLen, got, tt.want)
			}
		})
	}
}

// Benchmarks
func BenchmarkIsPalindrome(b *testing.B) {
	s := "A man a plan a canal Panama"
	for i := 0; i < b.N; i++ {
		IsPalindrome(s)
	}
}

func BenchmarkWordFrequency(b *testing.B) {
	text := "the quick brown fox jumps over the lazy dog"
	for i := 0; i < b.N; i++ {
		WordFrequency(text)
	}
}

func BenchmarkLevenshteinDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LevenshteinDistance("kitten", "sitting")
	}
}
