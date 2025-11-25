package main

import (
	"fmt"
	"os"
)

// Exercise: Append data to an existing file
// Learn how to open files in append mode

func main() {
	filename := "log.txt"

	// Create initial file
	err := os.WriteFile(filename, []byte("Log entry 1\n"), 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	// Open file in append mode
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Append data
	entries := []string{
		"Log entry 2\n",
		"Log entry 3\n",
		"Log entry 4\n",
	}

	for _, entry := range entries {
		_, err := file.WriteString(entry)
		if err != nil {
			fmt.Println("Error appending to file:", err)
			return
		}
	}

	fmt.Println("Successfully appended to", filename)

	// Read and display
	content, _ := os.ReadFile(filename)
	fmt.Println("\nFile contents:")
	fmt.Println(string(content))

	// Clean up
	os.Remove(filename)
}
