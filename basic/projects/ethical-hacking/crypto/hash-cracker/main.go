package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ⚠️ LEGAL WARNING ⚠️
// This tool is for EDUCATIONAL PURPOSES ONLY.
// Only crack hashes you own or have explicit permission to test.

type HashType string

const (
	MD5    HashType = "md5"
	SHA1   HashType = "sha1"
	SHA256 HashType = "sha256"
	BCRYPT HashType = "bcrypt"
)

func main() {
	printBanner()

	// Parse flags
	hash := flag.String("hash", "", "Hash to crack")
	hashType := flag.String("type", "md5", "Hash type (md5, sha1, sha256, bcrypt)")
	wordlist := flag.String("wordlist", "", "Path to wordlist file")
	flag.Parse()

	if *hash == "" {
		fmt.Println("❌ Error: Hash is required")
		fmt.Println("Usage: hash-cracker -hash <hash> -type md5 -wordlist wordlist.txt")
		return
	}

	fmt.Printf("🎯 Hash: %s\n", *hash)
	fmt.Printf("🔧 Type: %s\n", *hashType)

	var words []string

	if *wordlist != "" {
		fmt.Printf("📖 Wordlist: %s\n", *wordlist)
		var err error
		words, err = loadWordlist(*wordlist)
		if err != nil {
			fmt.Printf("❌ Error loading wordlist: %v\n", err)
			return
		}
		fmt.Printf("📊 Loaded %d words\n\n", len(words))
	} else {
		// Use common passwords
		words = getCommonPasswords()
		fmt.Printf("📊 Using %d common passwords\n\n", len(words))
	}

	// Start cracking
	fmt.Println("🔨 Starting hash cracking...")
	start := time.Now()

	result := crackHash(*hash, HashType(*hashType), words)

	elapsed := time.Since(start)

	if result != "" {
		fmt.Printf("\n✅ HASH CRACKED!\n")
		fmt.Printf("   Password: %s\n", result)
		fmt.Printf("   Time: %s\n", elapsed)
		fmt.Printf("   Speed: %.0f hashes/sec\n", float64(len(words))/elapsed.Seconds())
	} else {
		fmt.Printf("\n❌ Hash not cracked\n")
		fmt.Printf("   Tried %d passwords in %s\n", len(words), elapsed)
		fmt.Printf("   Speed: %.0f hashes/sec\n", float64(len(words))/elapsed.Seconds())
	}
}

func crackHash(targetHash string, hashType HashType, words []string) string {
	for i, word := range words {
		if i%1000 == 0 {
			fmt.Printf("\r🔨 Progress: %d/%d (%.1f%%)", i, len(words), float64(i)/float64(len(words))*100)
		}

		var hash string
		switch hashType {
		case MD5:
			hash = hashMD5(word)
		case SHA1:
			hash = hashSHA1(word)
		case SHA256:
			hash = hashSHA256(word)
		case BCRYPT:
			if compareBcrypt(targetHash, word) {
				return word
			}
			continue
		default:
			return ""
		}

		if hash == targetHash {
			return word
		}
	}

	return ""
}

func hashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func hashSHA1(text string) string {
	hash := sha1.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func hashSHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func compareBcrypt(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func loadWordlist(filename string) ([]string, error) {
	// Prevent path traversal attacks
	if strings.Contains(filename, "..") || filepath.IsAbs(filename) {
		return nil, fmt.Errorf("invalid path: path traversal detected")
	}

	// Only allow files in current directory
	baseDir, _ := os.Getwd()
	fullPath := filepath.Join(baseDir, filename)
	realPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Ensure the resolved path is within the current directory
	if !strings.HasPrefix(realPath, baseDir) {
		return nil, fmt.Errorf("invalid path: path traversal detected")
	}

	file, err := os.Open(realPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}

	return words, scanner.Err()
}

func getCommonPasswords() []string {
	return []string{
		"password", "123456", "12345678", "qwerty", "abc123",
		"monkey", "1234567", "letmein", "trustno1", "dragon",
		"baseball", "111111", "iloveyou", "master", "sunshine",
		"ashley", "bailey", "passw0rd", "shadow", "123123",
		"654321", "superman", "qazwsx", "michael", "football",
		"password1", "password123", "admin", "root", "toor",
		"pass", "test", "guest", "info", "adm", "mysql",
		"user", "administrator", "oracle", "ftp", "pi",
		"puppet", "ansible", "ec2-user", "vagrant", "azureuser",
		"demo", "changeme", "welcome", "login", "default",
	}
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  🔨 Hash Cracker                            ║
║                                                              ║
║              Password Hash Cracking Tool                    ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

⚠️  LEGAL WARNING: Only crack hashes you own or have explicit
   permission to test. Unauthorized password cracking may be illegal.

⚠️  EDUCATIONAL PURPOSE: This tool demonstrates hash cracking concepts.
   Real-world password cracking requires specialized tools and hardware.

`
	fmt.Println(banner)
}
