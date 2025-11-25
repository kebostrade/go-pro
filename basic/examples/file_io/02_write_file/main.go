package main

import (
	"fmt"
	"os"
)

// Exercise: Write data to a file
// Learn how to create and write to files

func main() {
	// Data to write
	data := []byte("This is written by Go!\nLine 2\nLine 3")

	// Write to file
	err := os.WriteFile("output.txt", data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Successfully wrote to output.txt")

	// Read it back to verify
	content, err := os.ReadFile("output.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("\nFile contents:")
	fmt.Println(string(content))

	// Clean up
	os.Remove("output.txt")
}
