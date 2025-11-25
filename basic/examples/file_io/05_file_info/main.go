package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Exercise: Get file information
// Learn how to retrieve file metadata using os.Stat

func main() {
	// Create a sample file
	filename := "info_test.txt"
	content := []byte("Sample file for testing file info")
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	// Get file info
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	// Display file information
	fmt.Println("File Information:")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Name:         %s\n", fileInfo.Name())
	fmt.Printf("Size:         %d bytes\n", fileInfo.Size())
	fmt.Printf("Permissions:  %s\n", fileInfo.Mode())
	fmt.Printf("Modified:     %s\n", fileInfo.ModTime().Format(time.RFC3339))
	fmt.Printf("Is Directory: %v\n", fileInfo.IsDir())

	// Clean up
	os.Remove(filename)
}
