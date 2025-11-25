package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

const (
	strengthVeryWeak   = "Very Weak"
	strengthWeak       = "Weak"
	strengthFair       = "Fair"
	strengthGood       = "Good"
	strengthStrong     = "Strong"
	strengthVeryStrong = "Very Strong"
)

type PasswordStrength struct {
	Password    string
	Score       int
	Strength    string
	Entropy     float64
	TimeToCrack string
	Issues      []string
	Suggestions []string
}

func main() {
	printBanner()

	// Parse flags
	password := flag.String("password", "", "Password to analyze")
	file := flag.String("file", "", "File containing passwords (one per line)")
	flag.Parse()

	if *password == "" && *file == "" {
		fmt.Println("❌ Error: Password or file is required")
		fmt.Println("Usage: password-analyzer -password mypassword")
		fmt.Println("   or: password-analyzer -file passwords.txt")
		return
	}

	if *password != "" {
		// Analyze single password
		result := analyzePassword(*password)
		displayResult(result)
	} else {
		// Analyze passwords from file
		analyzePasswordFile(*file)
	}
}

func analyzePassword(password string) PasswordStrength {
	result := PasswordStrength{
		Password: password,
		Score:    0,
	}

	// Check length
	scoreLength(&result, password)

	// Character variety
	scoreCharacterVariety(&result, password)

	// Common patterns
	scoreCommonPatterns(&result, password)

	// Calculate entropy
	result.Entropy = calculateEntropy(password)

	// Determine strength
	determineStrength(&result)

	// Estimate time to crack
	result.TimeToCrack = estimateTimeToCrack(result.Entropy)

	// Add suggestions
	addSuggestions(&result)

	return result
}

func scoreLength(result *PasswordStrength, password string) {
	length := len(password)
	if length < 8 {
		result.Issues = append(result.Issues, "Password is too short (minimum 8 characters)")
	} else if length >= 8 && length < 12 {
		result.Score += 10
		result.Suggestions = append(result.Suggestions, "Consider using 12+ characters for better security")
	} else if length >= 12 && length < 16 {
		result.Score += 20
	} else {
		result.Score += 30
	}
}

func scoreCharacterVariety(result *PasswordStrength, password string) {
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}

	if hasLower {
		result.Score += 10
	} else {
		result.Issues = append(result.Issues, "No lowercase letters")
	}

	if hasUpper {
		result.Score += 10
	} else {
		result.Issues = append(result.Issues, "No uppercase letters")
	}

	if hasDigit {
		result.Score += 10
	} else {
		result.Issues = append(result.Issues, "No numbers")
	}

	if hasSpecial {
		result.Score += 20
	} else {
		result.Issues = append(result.Issues, "No special characters")
	}
}

func scoreCommonPatterns(result *PasswordStrength, password string) {
	if isCommonPassword(password) {
		result.Score -= 50
		result.Issues = append(result.Issues, "Password is in common password list")
	}

	if hasSequentialChars(password) {
		result.Score -= 10
		result.Issues = append(result.Issues, "Contains sequential characters (e.g., abc, 123)")
	}

	if hasRepeatingChars(password) {
		result.Score -= 10
		result.Issues = append(result.Issues, "Contains repeating characters (e.g., aaa, 111)")
	}
}

func determineStrength(result *PasswordStrength) {
	if result.Score >= 80 {
		result.Strength = strengthVeryStrong
	} else if result.Score >= 60 {
		result.Strength = strengthStrong
	} else if result.Score >= 40 {
		result.Strength = strengthFair
	} else if result.Score >= 20 {
		result.Strength = strengthWeak
	} else {
		result.Strength = strengthVeryWeak
	}
}

func addSuggestions(result *PasswordStrength) {
	if len(result.Issues) == 0 {
		result.Suggestions = append(result.Suggestions, "Password looks good!")
	} else {
		// Check for missing character types
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(result.Password)
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(result.Password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(result.Password)
		hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(result.Password)
		length := len(result.Password)

		if !hasLower || !hasUpper {
			result.Suggestions = append(result.Suggestions, "Use both uppercase and lowercase letters")
		}
		if !hasDigit {
			result.Suggestions = append(result.Suggestions, "Add numbers to your password")
		}
		if !hasSpecial {
			result.Suggestions = append(result.Suggestions, "Add special characters (!@#$%^&*)")
		}
		if length < 12 {
			result.Suggestions = append(result.Suggestions, "Increase password length to 12+ characters")
		}
	}
}

func calculateEntropy(password string) float64 {
	// Calculate character set size
	charsetSize := 0
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	if hasLower {
		charsetSize += 26
	}
	if hasUpper {
		charsetSize += 26
	}
	if hasDigit {
		charsetSize += 10
	}
	if hasSpecial {
		charsetSize += 32
	}

	// Entropy = log2(charset^length)
	return float64(len(password)) * math.Log2(float64(charsetSize))
}

func estimateTimeToCrack(entropy float64) string {
	// Assume 1 billion guesses per second
	guessesPerSecond := 1e9
	combinations := math.Pow(2, entropy)
	seconds := combinations / guessesPerSecond

	if seconds < 1 {
		return "Instant"
	} else if seconds < 60 {
		return fmt.Sprintf("%.0f seconds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.0f minutes", seconds/60)
	} else if seconds < 86400 {
		return fmt.Sprintf("%.0f hours", seconds/3600)
	} else if seconds < 31536000 {
		return fmt.Sprintf("%.0f days", seconds/86400)
	} else if seconds < 31536000*100 {
		return fmt.Sprintf("%.0f years", seconds/31536000)
	} else {
		return "Centuries"
	}
}

func isCommonPassword(password string) bool {
	common := []string{
		"password", "123456", "12345678", "qwerty", "abc123",
		"monkey", "letmein", "dragon", "111111", "baseball",
		"iloveyou", "trustno1", "1234567", "sunshine", "master",
	}

	lower := strings.ToLower(password)
	for _, p := range common {
		if lower == p {
			return true
		}
	}
	return false
}

func hasSequentialChars(password string) bool {
	sequential := []string{
		"abc", "bcd", "cde", "def", "efg", "fgh", "ghi", "hij",
		"123", "234", "345", "456", "567", "678", "789",
		"qwerty", "asdfgh", "zxcvbn",
	}

	lower := strings.ToLower(password)
	for _, seq := range sequential {
		if strings.Contains(lower, seq) {
			return true
		}
	}
	return false
}

func hasRepeatingChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i] == password[i+2] {
			return true
		}
	}
	return false
}

func displayResult(result PasswordStrength) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 PASSWORD ANALYSIS")
	fmt.Println(strings.Repeat("=", 80))

	// Mask password
	masked := strings.Repeat("*", len(result.Password))
	fmt.Printf("\nPassword: %s (length: %d)\n", masked, len(result.Password))

	// Strength indicator
	fmt.Printf("\n🔒 Strength: ")
	switch result.Strength {
	case strengthVeryStrong:
		fmt.Printf("✅ %s\n", result.Strength)
	case strengthStrong:
		fmt.Printf("✅ %s\n", result.Strength)
	case strengthFair:
		fmt.Printf("⚠️  %s\n", result.Strength)
	case strengthWeak:
		fmt.Printf("❌ %s\n", result.Strength)
	case strengthVeryWeak:
		fmt.Printf("❌ %s\n", result.Strength)
	}

	fmt.Printf("📊 Score: %d/100\n", result.Score)
	fmt.Printf("🔢 Entropy: %.2f bits\n", result.Entropy)
	fmt.Printf("⏱️  Time to crack: %s\n", result.TimeToCrack)

	// Issues
	if len(result.Issues) > 0 {
		fmt.Println("\n❌ ISSUES:")
		for _, issue := range result.Issues {
			fmt.Printf("  • %s\n", issue)
		}
	}

	// Suggestions
	if len(result.Suggestions) > 0 {
		fmt.Println("\n💡 SUGGESTIONS:")
		for _, suggestion := range result.Suggestions {
			fmt.Printf("  • %s\n", suggestion)
		}
	}

	fmt.Println(strings.Repeat("=", 80))
}

func analyzePasswordFile(filename string) {
	// Prevent path traversal attacks
	if strings.Contains(filename, "..") || filepath.IsAbs(filename) {
		fmt.Printf("❌ Error: invalid path - path traversal detected\n")
		return
	}

	// Only allow files in current directory
	baseDir, _ := os.Getwd()
	fullPath := filepath.Join(baseDir, filename)
	realPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		fmt.Printf("❌ Error: invalid path: %v\n", err)
		return
	}

	// Ensure the resolved path is within the current directory
	if !strings.HasPrefix(realPath, baseDir) {
		fmt.Printf("❌ Error: invalid path - path traversal detected\n")
		return
	}

	file, err := os.Open(realPath)
	if err != nil {
		fmt.Printf("❌ Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Printf("📖 Analyzing passwords from: %s\n", filename)

	scanner := bufio.NewScanner(file)
	count := 0
	weak := 0
	medium := 0
	strong := 0

	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		if password == "" {
			continue
		}

		count++
		result := analyzePassword(password)

		switch result.Strength {
		case strengthVeryWeak, strengthWeak:
			weak++
		case strengthFair:
			medium++
		case strengthStrong, strengthVeryStrong:
			strong++
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Total passwords: %d\n", count)
	fmt.Printf("  Weak: %d (%.1f%%)\n", weak, float64(weak)/float64(count)*100)
	fmt.Printf("  Medium: %d (%.1f%%)\n", medium, float64(medium)/float64(count)*100)
	fmt.Printf("  Strong: %d (%.1f%%)\n", strong, float64(strong)/float64(count)*100)
	fmt.Println(strings.Repeat("=", 80))
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  🔐 Password Analyzer                       ║
║                                                              ║
║              Password Strength Analysis Tool                ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
