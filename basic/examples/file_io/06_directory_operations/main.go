package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Exercise: Work with directories
// Learn how to create, read, and remove directories

func main() {
	// Create a directory
	dirName := "test_directory"
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	fmt.Printf("Created directory: %s\n", dirName)

	// Create nested directories
	nestedDir := filepath.Join(dirName, "subdir1", "subdir2")
	err = os.MkdirAll(nestedDir, 0755)
	if err != nil {
		fmt.Println("Error creating nested directories:", err)
		return
	}
	fmt.Printf("Created nested directories: %s\n", nestedDir)

	// Create some files in the directory
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, filename := range files {
		path := filepath.Join(dirName, filename)
		err := os.WriteFile(path, []byte("Sample content"), 0644)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}
	}

	// Read directory contents
	fmt.Println("\nDirectory contents:")
	entries, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		fileType := "File"
		if entry.IsDir() {
			fileType = "Dir "
		}
		fmt.Printf("  [%s] %s\n", fileType, entry.Name())
	}

	// Clean up
	err = os.RemoveAll(dirName)
	if err != nil {
		fmt.Println("Error removing directory:", err)
		return
	}
	fmt.Printf("\nCleaned up: removed %s\n", dirName)
}
