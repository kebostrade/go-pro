package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Exercise: Read a file line by line
// Learn how to use bufio.Scanner for efficient line-by-line reading

func main() {
	// Create a sample file
	content := "Line 1: Hello\nLine 2: World\nLine 3: Go\nLine 4: Programming"
	err := os.WriteFile("lines.txt", []byte(content), 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	// Open file
	file, err := os.Open("lines.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create scanner
	scanner := bufio.NewScanner(file)

	lineNum := 1
	fmt.Println("Reading file line by line:")
	fmt.Println(strings.Repeat("-", 40))

	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%d: %s\n", lineNum, line)
		lineNum++
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// Clean up
	os.Remove("lines.txt")
}
