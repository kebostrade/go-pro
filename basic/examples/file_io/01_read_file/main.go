package main

import (
	"fmt"
	"os"
)

// Exercise: Read the entire contents of a file
// Learn how to read files using os.ReadFile

func main() {
	// Create a sample file first
	sampleContent := []byte("Hello, Go!\nThis is a sample file.\nLearning file I/O is fun!")
	err := os.WriteFile("sample.txt", sampleContent, 0644)
	if err != nil {
		fmt.Println("Error creating sample file:", err)
		return
	}

	// Read the file
	content, err := os.ReadFile("sample.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("File contents:")
	fmt.Println(string(content))

	// Clean up
	os.Remove("sample.txt")
}
