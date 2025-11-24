// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides curriculum lessons 11-15
package service

import "go-pro-backend/internal/domain"

// Lesson 11: Working with Files (I/O)
func (s *curriculumService) getComprehensiveLessonData11() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          11,
		Title:       "Working with Files (I/O)",
		Description: "Master file operations, directory handling, buffered I/O, and structured file formats in Go.",
		Duration:    "6-7 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Building Real Applications",
		Objectives: []string{
			"Read and write files efficiently",
			"Manage file permissions and metadata",
			"Work with directories and paths",
			"Use buffered I/O with scanners",
			"Handle JSON and XML files",
			"Implement proper error handling for file operations",
			"Work with temporary files and cleanup",
		},
		Theory: `# Working with Files (I/O)

## File Operations Fundamentals

File I/O is one of the most common tasks in real-world programs. Go provides powerful and simple APIs for working with files through the standard library's ` + "`os`" + ` and ` + "`io`" + ` packages.

## Reading Files

### Basic File Reading

The simplest way to read an entire file is using ` + "`ioutil.ReadFile`" + `:

` + "```go" + `
package main

import (
    "fmt"
    "os"
)

func main() {
    // Read entire file
    data, err := os.ReadFile("file.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    fmt.Println(string(data))
}
` + "```" + `

### Streaming File Reading

For large files, use buffered reading to avoid loading everything into memory:

` + "```go" + `
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("large_file.txt")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Scanner error:", err)
    }
}
` + "```" + `

## Writing Files

### Creating and Writing Files

` + "```go" + `
package main

import (
    "fmt"
    "os"
)

func main() {
    // Create new file (overwrite if exists)
    file, err := os.Create("output.txt")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    // Write content
    _, err = file.WriteString("Hello, World!\n")
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
}
` + "```" + `

### Appending to Files

To append instead of overwrite:

` + "```go" + `
package main

import (
    "fmt"
    "os"
)

func main() {
    // Open file for appending
    file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // Append data
    _, err = file.WriteString("New log entry\n")
    if err != nil {
        fmt.Println("Error writing to file:", err)
    }
}
` + "```" + `

### Buffered Writing

For better performance when writing many small pieces:

` + "```go" + `
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Create("output.txt")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    writer := bufio.NewWriter(file)

    // Write multiple items
    for i := 1; i <= 1000; i++ {
        fmt.Fprintf(writer, "Line %d\n", i)
    }

    // Flush any remaining data
    writer.Flush()
}
` + "```" + `

## File Paths and Directories

### Working with Paths

` + "```go" + `
package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    path := "/home/user/documents/file.txt"

    // Extract components
    dir := filepath.Dir(path)           // "/home/user/documents"
    file := filepath.Base(path)         // "file.txt"
    ext := filepath.Ext(path)           // ".txt"

    // Join paths (platform-safe)
    newPath := filepath.Join(dir, "new_file.txt")

    // Absolute path
    abs, _ := filepath.Abs(file)

    fmt.Println("Directory:", dir)
    fmt.Println("File:", file)
    fmt.Println("Extension:", ext)
    fmt.Println("New path:", newPath)
    fmt.Println("Absolute:", abs)
}
` + "```" + `

### Directory Operations

` + "```go" + `
package main

import (
    "fmt"
    "os"
)

func main() {
    // Create directory
    err := os.Mkdir("new_dir", 0755)
    if err != nil && !os.IsExist(err) {
        fmt.Println("Error creating directory:", err)
    }

    // Create nested directories
    err = os.MkdirAll("path/to/nested/dir", 0755)
    if err != nil {
        fmt.Println("Error creating directories:", err)
    }

    // List directory contents
    entries, err := os.ReadDir(".")
    if err != nil {
        fmt.Println("Error reading directory:", err)
        return
    }

    for _, entry := range entries {
        fmt.Println(entry.Name(), "- IsDir:", entry.IsDir())
    }

    // Remove directory
    err = os.RemoveAll("new_dir")
    if err != nil {
        fmt.Println("Error removing directory:", err)
    }
}
` + "```" + `

## File Permissions and Metadata

` + "```go" + `
package main

import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Stat("file.txt")
    if err != nil {
        fmt.Println("Error getting file info:", err)
        return
    }

    fmt.Println("Name:", file.Name())
    fmt.Println("Size:", file.Size(), "bytes")
    fmt.Println("Modified:", file.ModTime())
    fmt.Println("Permissions:", file.Mode())
    fmt.Println("Is Directory:", file.IsDir())

    // Change file permissions
    err = os.Chmod("file.txt", 0644)
    if err != nil {
        fmt.Println("Error changing permissions:", err)
    }
}
` + "```" + `

## JSON File Handling

### Reading JSON Files

` + "```go" + `
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Person struct {
    Name string ` + "`json:\"name\"`" + `
    Age  int    ` + "`json:\"age\"`" + `
    Email string ` + "`json:\"email\"`" + `
}

func main() {
    // Read JSON file
    data, err := os.ReadFile("person.json")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    var person Person
    err = json.Unmarshal(data, &person)
    if err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    fmt.Printf("Name: %s, Age: %d, Email: %s\n", person.Name, person.Age, person.Email)
}
` + "```" + `

### Writing JSON Files

` + "```go" + `
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Person struct {
    Name string ` + "`json:\"name\"`" + `
    Age  int    ` + "`json:\"age\"`" + `
}

func main() {
    person := Person{
        Name: "Alice",
        Age:  30,
    }

    // Marshal to JSON
    jsonData, err := json.MarshalIndent(person, "", "  ")
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        return
    }

    // Write to file
    err = os.WriteFile("person.json", jsonData, 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
    }
}
` + "```" + `

## Error Handling in File Operations

Common errors when working with files:

` + "```go" + `
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("missing_file.txt")
    if err != nil {
        // Check specific error types
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("File does not exist")
        } else if errors.Is(err, os.ErrPermission) {
            fmt.Println("Permission denied")
        } else {
            fmt.Println("Other error:", err)
        }
        return
    }
    defer file.Close()
}
` + "```" + `

## Best Practices

1. **Always defer file.Close()**: Ensures files are properly closed even if errors occur
2. **Use buffered I/O for large files**: Avoid memory issues with huge files
3. **Check errors immediately**: Don't skip error handling in file operations
4. **Use filepath.Join**: Platform-independent path construction
5. **Validate file permissions**: Check read/write before operating
6. **Clean up temporary files**: Remove temp files when done
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "file-read-write",
				Title:       "Basic File Read and Write",
				Description: "Create a program that reads a file and writes its content to another file.",
				Requirements: []string{
					"Read content from input.txt",
					"Write content to output.txt",
					"Handle errors appropriately",
					"Display file size information",
				},
				InitialCode: `package main

import "fmt"

func main() {
	// TODO: Read from input.txt
	// TODO: Write to output.txt
	// TODO: Handle errors
	// TODO: Print file size
}
`,
				Solution: `package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Read file
	data, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Write to new file
	err = os.WriteFile("output.txt", data, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	// Display file size
	fileInfo, err := os.Stat("output.txt")
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	fmt.Printf("File written successfully. Size: %d bytes\n", fileInfo.Size())
}
`,
			},
			{
				ID:          "line-counter",
				Title:       "Line Counter",
				Description: "Create a program that counts lines in a file using a scanner.",
				Requirements: []string{
					"Read file using bufio.Scanner",
					"Count total lines",
					"Count non-empty lines",
					"Display statistics",
				},
				InitialCode: `package main

import "fmt"

func main() {
	// TODO: Open file
	// TODO: Use scanner to read lines
	// TODO: Count lines
	// TODO: Print statistics
}
`,
				Solution: `package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("data.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalLines := 0
	nonEmptyLines := 0

	for scanner.Scan() {
		totalLines++
		line := scanner.Text()
		if len(line) > 0 {
			nonEmptyLines++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Non-empty lines: %d\n", nonEmptyLines)
}
`,
			},
			{
				ID:          "log-appender",
				Title:       "Log File Appender",
				Description: "Create a program that appends timestamped log entries to a file.",
				Requirements: []string{
					"Append to log.txt file",
					"Include timestamp for each entry",
					"Handle file creation if missing",
					"Use buffered writing",
				},
				InitialCode: `package main

import (
	"fmt"
	"time"
)

func LogEntry(message string) {
	// TODO: Get current time
	// TODO: Open/create log file
	// TODO: Append timestamped entry
	// TODO: Handle errors
}

func main() {
	LogEntry("Application started")
	LogEntry("Processing data")
	LogEntry("Application ended")
	fmt.Println("Logs written to log.txt")
}
`,
				Solution: `package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func LogEntry(message string) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(writer, "[%s] %s\n", timestamp, message)
	writer.Flush()
}

func main() {
	LogEntry("Application started")
	LogEntry("Processing data")
	LogEntry("Application ended")
	fmt.Println("Logs written to log.txt")
}
`,
			},
			{
				ID:          "directory-walker",
				Title:       "Directory Walker",
				Description: "Create a program that walks through a directory and lists all files.",
				Requirements: []string{
					"Read directory contents",
					"Distinguish files from directories",
					"Calculate total size",
					"Display formatted output",
				},
				InitialCode: `package main

import "fmt"

func main() {
	// TODO: Read current directory
	// TODO: Loop through entries
	// TODO: Check if file or directory
	// TODO: Calculate and display total size
}
`,
				Solution: `package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	var totalSize int64
	fileCount := 0
	dirCount := 0

	for _, entry := range entries {
		info, _ := entry.Info()
		if entry.IsDir() {
			fmt.Printf("[DIR]  %s\n", entry.Name())
			dirCount++
		} else {
			fmt.Printf("[FILE] %s (%d bytes)\n", entry.Name(), info.Size())
			totalSize += info.Size()
			fileCount++
		}
	}

	fmt.Printf("\nTotal: %d files, %d directories, %d bytes\n", fileCount, dirCount, totalSize)
}
`,
			},
			{
				ID:          "json-config",
				Title:       "JSON Configuration File",
				Description: "Create a program that reads and writes a JSON configuration file.",
				Requirements: []string{
					"Define struct with JSON tags",
					"Read JSON from file",
					"Modify configuration",
					"Write updated JSON to file",
					"Pretty-print JSON",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
)

// TODO: Define Config struct with JSON tags

func main() {
	// TODO: Read config.json
	// TODO: Modify configuration
	// TODO: Write updated config
	// TODO: Pretty-print
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppName    string ` + "`json:\"app_name\"`" + `
	Version    string ` + "`json:\"version\"`" + `
	Port       int    ` + "`json:\"port\"`" + `
	Debug      bool   ` + "`json:\"debug\"`" + `
	MaxWorkers int    ` + "`json:\"max_workers\"`" + `
}

func main() {
	// Read config
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Modify configuration
	config.Port = 8080
	config.Debug = false
	config.MaxWorkers = 10

	// Write updated config
	updated, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile("config.json", updated, 0644)
	if err != nil {
		log.Fatalf("Error writing config: %v", err)
	}

	// Display
	fmt.Println("Updated configuration:")
	fmt.Println(string(updated))
}
`,
			},
			{
				ID:          "csv-processor",
				Title:       "CSV File Processor",
				Description: "Create a program that reads a CSV file and processes its contents.",
				Requirements: []string{
					"Read CSV file line by line",
					"Parse comma-separated values",
					"Store in appropriate data structure",
					"Display statistics",
					"Handle quoted fields",
				},
				InitialCode: `package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	// TODO: Open CSV file
	// TODO: Read lines
	// TODO: Parse CSV
	// TODO: Store data
	// TODO: Display statistics
}
`,
				Solution: `package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type Record struct {
	Name   string
	Age    string
	Email  string
}

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineCount := 0

	// Skip header
	_, _ = reader.Read()

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		record := Record{
			Name:  row[0],
			Age:   row[1],
			Email: row[2],
		}
		records = append(records, record)
		lineCount++
	}

	fmt.Printf("Processed %d records:\n", lineCount)
	for _, r := range records {
		fmt.Printf("  %s (Age: %s) - %s\n", r.Name, r.Age, r.Email)
	}
}
`,
			},
			{
				ID:          "temp-file-handler",
				Title:       "Temporary File Handler",
				Description: "Create a program that safely works with temporary files.",
				Requirements: []string{
					"Create temporary file",
					"Write data to temp file",
					"Read from temp file",
					"Delete temp file",
					"Verify cleanup",
				},
				InitialCode: `package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// TODO: Create temporary file
	// TODO: Write data
	// TODO: Read data
	// TODO: Delete file
	// TODO: Verify deletion
}
`,
				Solution: `package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "temp-*.txt")
	if err != nil {
		log.Fatalf("Error creating temp file: %v", err)
	}
	tempName := tmpFile.Name()
	defer os.Remove(tempName)

	// Write data
	_, err = tmpFile.WriteString("Temporary data for processing\n")
	if err != nil {
		log.Fatalf("Error writing to temp file: %v", err)
	}
	tmpFile.Close()

	// Read data back
	data, err := os.ReadFile(tempName)
	if err != nil {
		log.Fatalf("Error reading temp file: %v", err)
	}
	fmt.Println("Temp file content:", string(data))

	// File will be deleted by defer
	_, err = os.Stat(tempName)
	if os.IsNotExist(err) {
		fmt.Println("Temp file successfully cleaned up")
	}
}
`,
			},
		},
	}
}

// Lesson 12: Concurrency Basics (Goroutines)
func (s *curriculumService) getComprehensiveLessonData12() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          12,
		Title:       "Concurrency Basics (Goroutines)",
		Description: "Master goroutines, concurrent execution, synchronization, and common pitfalls in concurrent Go programs.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Building Real Applications",
		Objectives: []string{
			"Understand concurrent vs parallel execution",
			"Create and manage goroutines",
			"Synchronize goroutines with WaitGroups",
			"Identify and fix race conditions",
			"Understand scheduler behavior",
			"Apply goroutines to real problems",
			"Debug concurrent programs",
		},
		Theory: `# Concurrency Basics (Goroutines)

## What Are Goroutines?

A goroutine is a lightweight thread managed by the Go runtime. Unlike operating system threads, which are heavyweight and expensive to create, goroutines are designed to be cheap and plentiful. You can have thousands or even millions of goroutines running concurrently.

Goroutines are not threads, and they are not processes. They are abstractions provided by the Go runtime that allow you to write concurrent code that is simpler, more efficient, and more scalable.

## Starting Goroutines

Creating a goroutine is remarkably simple—just use the ` + "`go`" + ` keyword:

` + "```go" + `
package main

import (
    "fmt"
    "time"
)

func greeting(name string) {
    for i := 1; i <= 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Start goroutines
    go greeting("Alice")
    go greeting("Bob")

    // Wait for goroutines to finish
    time.Sleep(1 * time.Second)
    fmt.Println("Done!")
}
` + "```" + `

## Concurrent vs Parallel

**Concurrency** is about managing multiple tasks. On a single-core processor, goroutines take turns running. **Parallelism** is about executing multiple tasks simultaneously on multiple cores.

Go's scheduler multiplexes many goroutines onto a smaller number of OS threads, allowing efficient concurrent execution even on single-core systems.

## The Race Condition Problem

When multiple goroutines access the same variable without synchronization, you get race conditions:

` + "```go" + `
package main

import "fmt"

var counter = 0

func increment() {
    for i := 0; i < 1000; i++ {
        counter++  // NOT SAFE - Race condition!
    }
}

func main() {
    go increment()
    go increment()
    go increment()

    fmt.Println(counter)  // Unpredictable result
}
` + "```" + `

The problem: ` + "`counter++`" + ` is not atomic. It's actually three operations: read, increment, write. Multiple goroutines can interleave, causing lost updates.

## Synchronization with WaitGroup

Use ` + "`sync.WaitGroup`" + ` to ensure all goroutines complete before the main goroutine exits:

` + "```go" + `
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d starting\n", id)
    // Do work...
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
    fmt.Println("All workers finished")
}
` + "```" + `

How it works:
1. ` + "`wg.Add(1)`" + ` increments the counter
2. ` + "`defer wg.Done()`" + ` decrements when goroutine finishes
3. ` + "`wg.Wait()`" + ` blocks until counter reaches 0

## Protecting Shared Data with Mutex

Use ` + "`sync.Mutex`" + ` to protect shared data:

` + "```go" + `
package main

import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    var counter Counter
    var wg sync.WaitGroup

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                counter.Increment()
            }
        }()
    }

    wg.Wait()
    fmt.Printf("Final count: %d\n", counter.Value())
}
` + "```" + `

## Detecting Race Conditions

Go includes a race detector. Run your tests with the ` + "`-race`" + ` flag:

` + "```bash" + `
go run -race main.go
go test -race ./...
` + "```" + `

This adds instrumentation to detect race conditions at runtime.

## Common Goroutine Patterns

### Fan-Out Pattern

Create multiple goroutines to process data in parallel:

` + "```go" + `
func processItems(items []int) {
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        go func(value int) {
            defer wg.Done()
            result := expensiveOperation(value)
            fmt.Printf("Result: %d\n", result)
        }(item)
    }

    wg.Wait()
}
` + "```" + `

### Worker Pool Pattern

Reuse a fixed number of goroutines to avoid creating too many:

` + "```go" + `
func workerPool(numWorkers int, jobs <-chan int) {
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                fmt.Printf("Worker %d processing job %d\n", workerID, job)
            }
        }(i)
    }

    wg.Wait()
}
` + "```" + `

## Goroutine Leaks

A goroutine leak occurs when a goroutine never exits, preventing garbage collection:

` + "```go" + `
// BAD - goroutine leak
func badListen() {
    for {
        fmt.Println("Listening...")
        time.Sleep(1 * time.Second)
    }
}

func main() {
    go badListen()  // This goroutine will run forever
    time.Sleep(2 * time.Second)
}

// GOOD - goroutine can be stopped
func goodListen(done <-chan struct{}) {
    for {
        select {
        case <-done:
            return
        default:
            fmt.Println("Listening...")
            time.Sleep(1 * time.Second)
        }
    }
}

func main() {
    done := make(chan struct{})
    go goodListen(done)
    time.Sleep(2 * time.Second)
    close(done)
}
` + "```" + `

## Goroutine Best Practices

1. **Always synchronize**: Use WaitGroup, channels, or context to ensure clean shutdown
2. **Avoid goroutine leaks**: Make sure all goroutines eventually exit
3. **Use buffered channels for fan-out**: Prevent goroutines from blocking
4. **Keep critical sections small**: Minimize time locks are held
5. **Run with -race flag**: During development and testing
6. **Document goroutine lifetime**: Make it clear when goroutines start/stop
7. **Use context for cancellation**: Pass context through goroutines for coordinated shutdown
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "basic-goroutine",
				Title:       "Basic Goroutine Creation",
				Description: "Create multiple goroutines and wait for them to complete.",
				Requirements: []string{
					"Create 5 worker goroutines",
					"Each worker processes tasks",
					"Use WaitGroup for synchronization",
					"Print completion messages",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
)

func worker(id int, wg *sync.WaitGroup) {
	// TODO: Process work
	// TODO: Signal completion
}

func main() {
	// TODO: Create WaitGroup
	// TODO: Start goroutines
	// TODO: Wait for all to complete
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Duration(id*100) * time.Millisecond)
	fmt.Printf("Worker %d finished\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers completed")
}
`,
			},
			{
				ID:          "safe-counter",
				Title:       "Thread-Safe Counter",
				Description: "Create a counter that multiple goroutines can safely increment.",
				Requirements: []string{
					"Define Counter with Mutex",
					"Implement Increment method",
					"Implement Value method",
					"Test with concurrent increments",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
)

// TODO: Define Counter struct with Mutex

// TODO: Implement Increment method

// TODO: Implement Value method

func main() {
	// TODO: Create counter
	// TODO: Start goroutines
	// TODO: Print final count
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	counter := &Counter{}
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.Value())
}
`,
			},
			{
				ID:          "parallel-processing",
				Title:       "Parallel Data Processing",
				Description: "Process a list of items in parallel using goroutines.",
				Requirements: []string{
					"Process multiple items concurrently",
					"Collect results safely",
					"Handle synchronization",
					"Measure performance",
				},
				InitialCode: `package main

import (
	"fmt"
	"time"
)

func processItem(item int) int {
	// Simulate work
	time.Sleep(100 * time.Millisecond)
	return item * 2
}

func main() {
	items := []int{1, 2, 3, 4, 5}

	// TODO: Process items in parallel
	// TODO: Collect results
	// TODO: Print results
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

func processItem(item int) int {
	time.Sleep(100 * time.Millisecond)
	return item * 2
}

func main() {
	items := []int{1, 2, 3, 4, 5}
	results := make([]int, len(items))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, item := range items {
		wg.Add(1)
		go func(index, value int) {
			defer wg.Done()
			result := processItem(value)
			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, item)
	}

	wg.Wait()
	fmt.Println("Results:", results)
}
`,
			},
			{
				ID:          "worker-pool",
				Title:       "Worker Pool",
				Description: "Implement a worker pool pattern with a fixed number of workers.",
				Requirements: []string{
					"Create fixed number of workers",
					"Distribute jobs to workers",
					"Workers process jobs concurrently",
					"Wait for all jobs to complete",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	// TODO: Receive jobs
	// TODO: Process jobs
	// TODO: Signal completion
}

func main() {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	// TODO: Create workers
	// TODO: Send jobs
	// TODO: Wait for completion

	// Close channel to signal workers
	close(jobs)
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Done()
}

func main() {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	for j := 1; j <= 10; j++ {
		jobs <- j
	}

	close(jobs)
	wg.Wait()
	fmt.Println("All jobs completed")
}
`,
			},
			{
				ID:          "goroutine-coordination",
				Title:       "Goroutine Coordination",
				Description: "Coordinate multiple goroutines using channels and synchronization.",
				Requirements: []string{
					"Start multiple goroutines",
					"Coordinate their execution",
					"Handle completion signals",
					"Ensure clean shutdown",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
)

func task(id int, ready, done chan struct{}, wg *sync.WaitGroup) {
	// TODO: Wait for ready signal
	// TODO: Do work
	// TODO: Signal completion
}

func main() {
	// TODO: Create channels
	// TODO: Start tasks
	// TODO: Coordinate execution
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

func task(id int, ready, done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ready
	fmt.Printf("Task %d started\n", id)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Task %d completed\n", id)
	done <- struct{}{}
}

func main() {
	ready := make(chan struct{})
	done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go task(i, ready, done, &wg)
	}

	time.Sleep(500 * time.Millisecond)
	close(ready)

	for i := 0; i < 3; i++ {
		<-done
	}

	wg.Wait()
	fmt.Println("All tasks completed")
}
`,
			},
			{
				ID:          "race-condition-fix",
				Title:       "Race Condition Fix",
				Description: "Identify and fix a race condition in concurrent code.",
				Requirements: []string{
					"Identify the race condition",
					"Protect shared data with Mutex",
					"Ensure thread safety",
					"Test with -race flag",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
)

var sharedData = 0

func modifier(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		sharedData++  // Race condition here!
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go modifier(&wg)
	}

	wg.Wait()
	fmt.Printf("Final value: %d (expected: 10000)\n", sharedData)
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
)

var (
	sharedData = 0
	mu         sync.Mutex
)

func modifier(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		mu.Lock()
		sharedData++
		mu.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go modifier(&wg)
	}

	wg.Wait()
	fmt.Printf("Final value: %d (expected: 10000)\n", sharedData)
}
`,
			},
			{
				ID:          "stopable-goroutine",
				Title:       "Stoppable Goroutine",
				Description: "Create a goroutine that can be stopped cleanly via a done channel.",
				Requirements: []string{
					"Create long-running goroutine",
					"Implement stop mechanism",
					"Ensure clean shutdown",
					"Verify goroutine exits",
				},
				InitialCode: `package main

import (
	"fmt"
	"time"
)

func worker(done <-chan struct{}) {
	// TODO: Loop until stop signal
	// TODO: Handle done channel
	// TODO: Clean up and exit
}

func main() {
	done := make(chan struct{})

	// TODO: Start worker
	// TODO: Let it run
	// TODO: Send stop signal
	// TODO: Wait for completion
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-done:
			fmt.Println("Worker stopped")
			return
		default:
			fmt.Println("Worker running...")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go worker(done, &wg)

	time.Sleep(500 * time.Millisecond)
	close(done)

	wg.Wait()
	fmt.Println("Main completed")
}
`,
			},
		},
	}
}

// Lesson 13: Channels and Communication
func (s *curriculumService) getComprehensiveLessonData13() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          13,
		Title:       "Channels and Communication",
		Description: "Master channels, buffering, select statements, and channel patterns for goroutine communication.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Building Real Applications",
		Objectives: []string{
			"Understand channel fundamentals",
			"Create and use buffered channels",
			"Send and receive on channels",
			"Close channels appropriately",
			"Use select statements for multiplexing",
			"Implement fan-in and fan-out patterns",
			"Avoid deadlocks and channel leaks",
		},
		Theory: `# Channels and Communication

## What Are Channels?

Channels provide a way for goroutines to communicate with each other and synchronize their execution. Instead of protecting shared memory with locks, Go encourages you to pass data through channels. The motto is: "Share memory by communicating, rather than communicate by sharing memory."

A channel is a typed conduit through which you can send and receive values.

## Creating Channels

Channels are created with ` + "`make`" + `:

` + "```go" + `
// Unbuffered channel - senders block until receiver is ready
messages := make(chan string)

// Buffered channel - senders can write until buffer is full
requests := make(chan int, 100)

// Channel of custom types
results := make(chan MyType, 10)
` + "```" + `

## Sending and Receiving

` + "```go" + `
package main

import (
    "fmt"
)

func main() {
    messages := make(chan string)

    go func() {
        messages <- "hello"  // Send
    }()

    msg := <-messages  // Receive
    fmt.Println(msg)
}
` + "```" + `

## Unbuffered vs Buffered Channels

### Unbuffered Channels

Unbuffered channels require both sender and receiver to be ready:

` + "```go" + `
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan int)

    go func() {
        fmt.Println("Sending 42...")
        ch <- 42  // Blocks until main receives
        fmt.Println("Sent!")
    }()

    time.Sleep(1 * time.Second)
    fmt.Println("Receiving...")
    value := <-ch
    fmt.Println("Received:", value)
}
` + "```" + `

### Buffered Channels

Buffered channels allow sending without an immediate receiver:

` + "```go" + `
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 2)

    ch <- 1
    ch <- 2

    fmt.Println(<-ch)  // 1
    fmt.Println(<-ch)  // 2
}
` + "```" + `

## Closing Channels

Only the sender should close a channel:

` + "```go" + `
package main

import (
    "fmt"
)

func main() {
    numbers := make(chan int)

    go func() {
        for i := 1; i <= 3; i++ {
            numbers <- i
        }
        close(numbers)  // Signal that no more values will be sent
    }()

    for num := range numbers {
        fmt.Println(num)
    }
}
` + "```" + `

## The Select Statement

Select allows waiting on multiple channel operations:

` + "```go" + `
package main

import (
    "fmt"
    "time"
)

func main() {
    tick := time.Tick(100 * time.Millisecond)
    boom := time.After(500 * time.Millisecond)

    for {
        select {
        case t := <-tick:
            fmt.Println("tick.", t)
        case t := <-boom:
            fmt.Println("BOOM!", t)
            return
        default:
            fmt.Println("    .")
            time.Sleep(50 * time.Millisecond)
        }
    }
}
` + "```" + `

## Fan-In Pattern

Multiple producers sending to one consumer:

` + "```go" + `
package main

import (
    "fmt"
)

func producer(id int, out chan<- string) {
    out <- fmt.Sprintf("Message from producer %d", id)
}

func main() {
    out := make(chan string)

    go producer(1, out)
    go producer(2, out)
    go producer(3, out)

    for i := 0; i < 3; i++ {
        fmt.Println(<-out)
    }
}
` + "```" + `

## Fan-Out Pattern

One producer sending to multiple consumers:

` + "```go" + `
package main

import (
    "fmt"
    "sync"
)

func main() {
    messages := make(chan string, 10)

    // Publish messages
    for i := 1; i <= 3; i++ {
        messages <- fmt.Sprintf("Message %d", i)
    }
    close(messages)

    // Multiple consumers
    var wg sync.WaitGroup
    for consumer := 1; consumer <= 2; consumer++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for msg := range messages {
                fmt.Printf("Consumer %d received: %s\n", id, msg)
            }
        }(consumer)
    }

    wg.Wait()
}
` + "```" + `

## Pipeline Pattern

Chain goroutines where output of one is input to another:

` + "```go" + `
func square(numbers <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range numbers {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    numbers := make(chan int)

    go func() {
        for i := 1; i <= 5; i++ {
            numbers <- i
        }
        close(numbers)
    }()

    squares := square(numbers)

    for s := range squares {
        fmt.Println(s)
    }
}
` + "```" + `

## Preventing Deadlocks

Common deadlock causes:

1. All goroutines blocked on channel operations
2. Receiving from closed channel
3. Sending to closed channel

` + "```go" + `
// DEADLOCK - waiting on empty channel with no goroutines to send
ch := make(chan int)
value := <-ch  // Will block forever

// FIX - use buffering or goroutine
ch := make(chan int, 1)
ch <- 1
value := <-ch  // Can receive
` + "```" + `

## Channel Best Practices

1. **Only sender closes**: The sender knows when sending is done
2. **Use send-only/receive-only types**: Make intent clear
3. **Close before empty check**: Check if receive succeeds, not if closed
4. **Avoid nil channels**: Sending on nil channel panics
5. **Use context for cancellation**: Better than custom done channels
6. **Prefer channels over shared memory**: Use for synchronization
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "basic-channel",
				Title:       "Basic Channel Communication",
				Description: "Implement simple send and receive on a channel.",
				Requirements: []string{
					"Create a channel of strings",
					"Send value from goroutine",
					"Receive value in main",
					"Print received message",
				},
				InitialCode: `package main

import "fmt"

func main() {
	// TODO: Create channel
	// TODO: Send value in goroutine
	// TODO: Receive value
	// TODO: Print result
}
`,
				Solution: `package main

import "fmt"

func main() {
	messages := make(chan string)

	go func() {
		messages <- "Hello from goroutine"
	}()

	msg := <-messages
	fmt.Println(msg)
}
`,
			},
			{
				ID:          "buffered-channel",
				Title:       "Buffered Channels",
				Description: "Work with buffered channels and understand their behavior.",
				Requirements: []string{
					"Create buffered channel with capacity 3",
					"Send values without blocking",
					"Receive values in order",
					"Demonstrate buffering advantage",
				},
				InitialCode: `package main

import "fmt"

func main() {
	// TODO: Create buffered channel with capacity 3
	// TODO: Send 3 values
	// TODO: Receive values
	// TODO: Verify order
}
`,
				Solution: `package main

import "fmt"

func main() {
	ch := make(chan int, 3)

	ch <- 1
	ch <- 2
	ch <- 3

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
`,
			},
			{
				ID:          "channel-range",
				Title:       "Iterating Over Channels",
				Description: "Use range to iterate over channel values until closed.",
				Requirements: []string{
					"Create channel",
					"Send multiple values from goroutine",
					"Close channel when done",
					"Use range to receive all values",
				},
				InitialCode: `package main

import "fmt"

func main() {
	numbers := make(chan int)

	// TODO: Send numbers 1-5 in goroutine
	// TODO: Close channel
	// TODO: Range over channel
	// TODO: Print received numbers
}
`,
				Solution: `package main

import "fmt"

func main() {
	numbers := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers)
	}()

	for num := range numbers {
		fmt.Println(num)
	}
}
`,
			},
			{
				ID:          "select-statement",
				Title:       "Select Statement",
				Description: "Use select to wait on multiple channel operations.",
				Requirements: []string{
					"Create multiple channels",
					"Use select to multiplex channels",
					"Handle different channel cases",
					"Implement timeout behavior",
				},
				InitialCode: `package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// TODO: Create goroutines sending to channels
	// TODO: Use select to receive
	// TODO: Handle timeout
	// TODO: Print results
}
`,
				Solution: `package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "one"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
	}
}
`,
			},
			{
				ID:          "fan-in",
				Title:       "Fan-In Pattern",
				Description: "Implement fan-in: multiple producers to one consumer.",
				Requirements: []string{
					"Create multiple producer goroutines",
					"All send to same channel",
					"Main receives from all",
					"Use select or range to receive",
				},
				InitialCode: `package main

import "fmt"

func producer(id int, out chan<- string) {
	// TODO: Send message to channel
}

func main() {
	out := make(chan string)

	// TODO: Start multiple producers
	// TODO: Receive messages
	// TODO: Print results
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(id int, out chan<- string) {
	out <- fmt.Sprintf("Message from producer %d", id)
}

func main() {
	out := make(chan string)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			producer(id, out)
		}(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for msg := range out {
		fmt.Println(msg)
	}
}
`,
			},
			{
				ID:          "fan-out",
				Title:       "Fan-Out Pattern",
				Description: "Implement fan-out: one producer to multiple consumers.",
				Requirements: []string{
					"Create channel with multiple consumers",
					"All consumers read same messages",
					"Proper synchronization",
					"Clean shutdown",
				},
				InitialCode: `package main

import (
	"fmt"
	"sync"
)

func main() {
	messages := make(chan string)
	var wg sync.WaitGroup

	// TODO: Create consumers
	// TODO: Send messages
	// TODO: Wait for completion
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
)

func main() {
	messages := make(chan string, 10)
	var wg sync.WaitGroup

	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for msg := range messages {
				fmt.Printf("Consumer %d: %s\n", id, msg)
			}
		}(i)
	}

	for j := 1; j <= 5; j++ {
		messages <- fmt.Sprintf("Message %d", j)
	}

	close(messages)
	wg.Wait()
}
`,
			},
			{
				ID:          "pipeline",
				Title:       "Pipeline Pattern",
				Description: "Create a pipeline: stage1 -> stage2 -> stage3.",
				Requirements: []string{
					"Create three pipeline stages",
					"Stage 1 generates numbers",
					"Stage 2 squares numbers",
					"Stage 3 doubles squared numbers",
					"Print final results",
				},
				InitialCode: `package main

import "fmt"

func generate(max int) <-chan int {
	// TODO: Create channel
	// TODO: Send numbers 1 to max
	// TODO: Close and return
}

func square(numbers <-chan int) <-chan int {
	// TODO: Create channel
	// TODO: Square each number
	// TODO: Close and return
}

func main() {
	// TODO: Connect pipeline
	// TODO: Print results
}
`,
				Solution: `package main

import "fmt"

func generate(max int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 1; i <= max; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func square(numbers <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range numbers {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	for n := range square(generate(5)) {
		fmt.Println(n)
	}
}
`,
			},
			{
				ID:          "channel-direction",
				Title:       "Channel Direction Types",
				Description: "Use send-only and receive-only channel types.",
				Requirements: []string{
					"Create send-only channel parameter",
					"Create receive-only channel parameter",
					"Implement producer with send-only",
					"Implement consumer with receive-only",
				},
				InitialCode: `package main

import "fmt"

func producer(out chan<- int) {
	// TODO: Send values to channel
}

func consumer(in <-chan int) {
	// TODO: Receive values from channel
}

func main() {
	ch := make(chan int)

	// TODO: Start producer and consumer
	// TODO: Ensure proper synchronization
}
`,
				Solution: `package main

import (
	"fmt"
	"sync"
)

func producer(out chan<- int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

func consumer(in <-chan int) {
	for value := range in {
		fmt.Println("Received:", value)
	}
}

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer(ch)
	}()

	producer(ch)
	wg.Wait()
}
`,
			},
		},
	}
}

// Lesson 14: Working with JSON and APIs
func (s *curriculumService) getComprehensiveLessonData14() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          14,
		Title:       "Working with JSON and APIs",
		Description: "Master JSON handling, HTTP clients, REST API consumption, and error handling in API operations.",
		Duration:    "6-7 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Building Real Applications",
		Objectives: []string{
			"Marshal and unmarshal JSON",
			"Use struct tags for JSON mapping",
			"Create HTTP clients",
			"Consume REST APIs",
			"Handle API responses and errors",
			"Implement request/response validation",
			"Work with API authentication",
		},
		Theory: `# Working with JSON and APIs

## JSON Encoding and Decoding

JSON is the de facto standard for APIs in Go. The ` + "`encoding/json`" + ` package makes working with JSON straightforward.

### Marshaling (Go to JSON)

Converting Go data structures to JSON:

` + "```go" + `
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string
    Age  int
    Email string
}

func main() {
    p := Person{
        Name: "Alice",
        Age: 30,
        Email: "alice@example.com",
    }

    // Compact JSON
    data, _ := json.Marshal(p)
    fmt.Println(string(data))

    // Pretty JSON
    pretty, _ := json.MarshalIndent(p, "", "  ")
    fmt.Println(string(pretty))
}
` + "```" + `

### Unmarshaling (JSON to Go)

Converting JSON to Go data structures:

` + "```go" + `
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    jsonData := []byte(` + "`" + `{
        "Name": "Bob",
        "Age": 25
    }` + "`" + `)

    var person Person
    json.Unmarshal(jsonData, &person)
    fmt.Printf("%+v\n", person)
}
` + "```" + `

## Struct Tags

Struct tags control how JSON is encoded/decoded:

` + "```go" + `
type User struct {
    ID        int    ` + "`json:\"id\"`" + `
    Name      string ` + "`json:\"name\"`" + `
    Email     string ` + "`json:\"email\"`" + `
    Password  string ` + "`json:\"-\"`" + `  // Never marshal
    Active    bool   ` + "`json:\"active,omitempty\"`" + `  // Omit if zero
    CreatedAt string ` + "`json:\"created_at\"`" + `
}
` + "```" + `

Common tag options:
- ` + "`json:\"name\"`" + ` - Map to JSON field name
- ` + "`json:\"-\"`" + ` - Skip this field
- ` + "`json:\"name,omitempty\"`" + ` - Omit if empty
- ` + "`json:\"name,string\"`" + ` - Convert to/from string

## Custom JSON Marshaling

Implement custom marshaling:

` + "```go" + `
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

type Event struct {
    Name      string
    Timestamp time.Time
}

func (e Event) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Name      string ` + "`json:\"name\"`" + `
        Timestamp string ` + "`json:\"timestamp\"`" + `
    }{
        Name:      e.Name,
        Timestamp: e.Timestamp.Format(time.RFC3339),
    })
}

func main() {
    event := Event{
        Name: "Meeting",
        Timestamp: time.Now(),
    }

    data, _ := json.MarshalIndent(event, "", "  ")
    fmt.Println(string(data))
}
` + "```" + `

## HTTP Clients

### Making GET Requests

` + "```go" + `
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    resp, err := http.Get("https://api.example.com/users/1")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
` + "```" + `

### Parsing API Responses

` + "```go" + `
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type User struct {
    ID   int    ` + "`json:\"id\"`" + `
    Name string ` + "`json:\"name\"`" + `
    Email string ` + "`json:\"email\"`" + `
}

func main() {
    resp, _ := http.Get("https://api.example.com/users/1")
    defer resp.Body.Close()

    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    fmt.Printf("User: %+v\n", user)
}
` + "```" + `

### POST Requests

` + "```go" + `
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    data := map[string]string{
        "name": "Charlie",
        "email": "charlie@example.com",
    }

    jsonData, _ := json.Marshal(data)
    resp, _ := http.Post(
        "https://api.example.com/users",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)
}
` + "```" + `

## Request Configuration

For more control, use ` + "`http.Client`" + ` and ` + "`http.Request`" + `:

` + "```go" + `
package main

import (
    "fmt"
    "net/http"
    "time"
)

func main() {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)
    req.Header.Set("Authorization", "Bearer token")
    req.Header.Set("User-Agent", "MyApp/1.0")

    resp, _ := client.Do(req)
    defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)
}
` + "```" + `

## Error Handling

Always handle errors appropriately:

` + "```go" + `
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type ErrorResponse struct {
    Error   string ` + "`json:\"error\"`" + `
    Message string ` + "`json:\"message\"`" + `
}

func fetchUser(id int) error {
    client := &http.Client{Timeout: 5 * time.Second}

    resp, err := client.Get(fmt.Sprintf("https://api.example.com/users/%d", id))
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        var errResp ErrorResponse
        json.Unmarshal(body, &errResp)
        return fmt.Errorf("API error: %s", errResp.Message)
    }

    return nil
}

func main() {
    if err := fetchUser(1); err != nil {
        fmt.Println("Error:", err)
    }
}
` + "```" + `

## Best Practices

1. **Always set timeouts**: Prevent hanging requests
2. **Close response bodies**: Avoid resource leaks
3. **Check status codes**: Don't assume success
4. **Handle errors gracefully**: Provide meaningful error messages
5. **Validate input**: Check API requirements before sending
6. **Use appropriate methods**: GET for reads, POST for creates
7. **Include authentication**: Use headers for API keys
8. **Document API contracts**: Define expected structures
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "json-marshal",
				Title:       "JSON Marshaling",
				Description: "Convert Go structs to JSON with proper struct tags.",
				Requirements: []string{
					"Define struct with JSON tags",
					"Marshal struct to JSON",
					"Pretty-print JSON",
					"Handle multiple data types",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
)

// TODO: Define struct with JSON tags

func main() {
	// TODO: Create instance
	// TODO: Marshal to JSON
	// TODO: Pretty-print
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string ` + "`json:\"name\"`" + `
	Age   int    ` + "`json:\"age\"`" + `
	Email string ` + "`json:\"email\"`" + `
}

func main() {
	person := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	// Pretty JSON
	data, _ := json.MarshalIndent(person, "", "  ")
	fmt.Println(string(data))
}
`,
			},
			{
				ID:          "json-unmarshal",
				Title:       "JSON Unmarshaling",
				Description: "Parse JSON into Go structs with error handling.",
				Requirements: []string{
					"Define struct with JSON tags",
					"Parse JSON string",
					"Handle parsing errors",
					"Display parsed data",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
)

// TODO: Define struct

func main() {
	jsonData := []byte(` + "`{\"name\": \"Bob\", \"age\": 25}`" + `)

	// TODO: Unmarshal JSON
	// TODO: Handle errors
	// TODO: Display result
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Person struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func main() {
	jsonData := []byte(` + "`{\"name\": \"Bob\", \"age\": 25}`" + `)

	var person Person
	err := json.Unmarshal(jsonData, &person)
	if err != nil {
		log.Fatalf("JSON error: %v", err)
	}

	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
}
`,
			},
			{
				ID:          "http-get",
				Title:       "HTTP GET Request",
				Description: "Make HTTP GET request and parse JSON response.",
				Requirements: []string{
					"Create HTTP GET request",
					"Parse JSON response",
					"Handle HTTP errors",
					"Display results",
				},
				InitialCode: `package main

import (
	"fmt"
	"net/http"
)

// TODO: Define response struct

func main() {
	// TODO: Make GET request
	// TODO: Parse response
	// TODO: Handle errors
	// TODO: Print data
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	UserID int    ` + "`json:\"userId\"`" + `
	ID     int    ` + "`json:\"id\"`" + `
	Title  string ` + "`json:\"title\"`" + `
	Body   string ` + "`json:\"body\"`" + `
}

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		log.Fatalf("Request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status: %d", resp.StatusCode)
	}

	var post Post
	json.NewDecoder(resp.Body).Decode(&post)
	fmt.Printf("Title: %s\nBody: %s\n", post.Title, post.Body)
}
`,
			},
			{
				ID:          "http-post",
				Title:       "HTTP POST Request",
				Description: "Send JSON data with POST request.",
				Requirements: []string{
					"Create POST request with JSON body",
					"Set appropriate headers",
					"Parse response",
					"Handle errors",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: Define request/response structs

func main() {
	// TODO: Create request data
	// TODO: Marshal to JSON
	// TODO: Send POST request
	// TODO: Parse response
}
`,
				Solution: `package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CreateRequest struct {
	Title  string ` + "`json:\"title\"`" + `
	Body   string ` + "`json:\"body\"`" + `
	UserID int    ` + "`json:\"userId\"`" + `
}

func main() {
	data := CreateRequest{
		Title:  "New Post",
		Body:   "This is a test post",
		UserID: 1,
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(
		"https://jsonplaceholder.typicode.com/posts",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatalf("Request error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}
`,
			},
			{
				ID:          "api-client",
				Title:       "API Client with Error Handling",
				Description: "Create a reusable API client with proper error handling.",
				Requirements: []string{
					"Create APIClient struct",
					"Implement GET method",
					"Implement POST method",
					"Handle errors and status codes",
					"Support authentication headers",
				},
				InitialCode: `package main

import (
	"fmt"
	"net/http"
)

type APIClient struct {
	// TODO: Define fields
}

// TODO: Implement NewAPIClient

// TODO: Implement Get method

// TODO: Implement Post method

func main() {
	// TODO: Create client
	// TODO: Make requests
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewAPIClient(baseURL, token string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClient) Get(endpoint string, result interface{}) error {
	req, _ := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func main() {
	client := NewAPIClient("https://jsonplaceholder.typicode.com", "")

	type Post struct {
		ID    int    ` + "`json:\"id\"`" + `
		Title string ` + "`json:\"title\"`" + `
	}

	var post Post
	err := client.Get("/posts/1", &post)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Post: %+v\n", post)
}
`,
			},
			{
				ID:          "json-arrays",
				Title:       "JSON Arrays and Slices",
				Description: "Work with JSON arrays and Go slices.",
				Requirements: []string{
					"Define struct for array elements",
					"Unmarshal JSON array",
					"Filter/process elements",
					"Display results",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
)

// TODO: Define element struct

func main() {
	jsonData := []byte(` + "`[{\"id\":1,\"name\":\"Alice\"},{\"id\":2,\"name\":\"Bob\"}]`" + `)

	// TODO: Unmarshal array
	// TODO: Process elements
	// TODO: Print results
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Person struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

func main() {
	jsonData := []byte(` + "`[{\"id\":1,\"name\":\"Alice\"},{\"id\":2,\"name\":\"Bob\"}]`" + `)

	var people []Person
	err := json.Unmarshal(jsonData, &people)
	if err != nil {
		log.Fatalf("JSON error: %v", err)
	}

	for _, p := range people {
		fmt.Printf("ID: %d, Name: %s\n", p.ID, p.Name)
	}
}
`,
			},
			{
				ID:          "custom-json",
				Title:       "Custom JSON Marshaling",
				Description: "Implement custom MarshalJSON/UnmarshalJSON methods.",
				Requirements: []string{
					"Implement MarshalJSON method",
					"Implement UnmarshalJSON method",
					"Handle custom formatting",
					"Test serialization/deserialization",
				},
				InitialCode: `package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Name string
	Time time.Time
}

// TODO: Implement MarshalJSON

// TODO: Implement UnmarshalJSON

func main() {
	// TODO: Create event
	// TODO: Marshal to JSON
	// TODO: Unmarshal back
	// TODO: Display
}
`,
				Solution: `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Event struct {
	Name string
	Time time.Time
}

func (e Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	return json.Marshal(struct {
		Name      string ` + "`json:\"name\"`" + `
		Timestamp string ` + "`json:\"timestamp\"`" + `
	}{
		Name:      e.Name,
		Timestamp: e.Time.Format(time.RFC3339),
	})
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := struct {
		Name      string ` + "`json:\"name\"`" + `
		Timestamp string ` + "`json:\"timestamp\"`" + `
	}{}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	t, _ := time.Parse(time.RFC3339, aux.Timestamp)
	e.Name = aux.Name
	e.Time = t
	return nil
}

func main() {
	event := Event{
		Name: "Meeting",
		Time: time.Now(),
	}

	data, _ := json.Marshal(event)
	fmt.Println("Marshaled:", string(data))

	var event2 Event
	json.Unmarshal(data, &event2)
	fmt.Printf("Unmarshaled: %+v\n", event2)
}
`,
			},
		},
	}
}

// Lesson 15: Database Fundamentals
func (s *curriculumService) getComprehensiveLessonData15() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          15,
		Title:       "Database Fundamentals",
		Description: "Master database/sql package, CRUD operations, prepared statements, transactions, and connection pooling.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Building Real Applications",
		Objectives: []string{
			"Connect to databases using database/sql",
			"Execute queries and handle results",
			"Perform CRUD operations",
			"Use prepared statements safely",
			"Implement transactions",
			"Handle connection pooling",
			"Implement proper error handling",
			"Avoid SQL injection",
		},
		Theory: `# Database Fundamentals

## Overview of database/sql

Go's ` + "`database/sql`" + ` package provides a generic interface for SQL databases. It doesn't provide a database driver itself—you must import a driver for your specific database (e.g., postgres, mysql, sqlite3).

## Setting Up Drivers

To use a database, import both ` + "`database/sql`" + ` and a driver:

` + "```go" + `
import (
    "database/sql"
    _ "github.com/lib/pq"  // PostgreSQL
)

// For MySQL
// _ "github.com/go-sql-driver/mysql"

// For SQLite
// _ "github.com/mattn/go-sqlite3"
` + "```" + `

The underscore import ensures the driver registers itself with database/sql.

## Opening Database Connections

` + "```go" + `
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    dsn := "user=postgres password=secret dbname=mydb host=localhost port=5432 sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Test connection
    err = db.Ping()
    if err != nil {
        panic(err)
    }

    fmt.Println("Connected!")
}
` + "```" + `

## Querying Single Rows

Use ` + "`QueryRow()`" + ` for queries that return a single row:

` + "```go" + `
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

type User struct {
    ID    int
    Name  string
    Email string
}

func main() {
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()

    var user User
    err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", 1).
        Scan(&user.ID, &user.Name, &user.Email)

    if err == sql.ErrNoRows {
        fmt.Println("No user found")
    } else if err != nil {
        panic(err)
    }

    fmt.Printf("User: %+v\n", user)
}
` + "```" + `

## Querying Multiple Rows

Use ` + "`Query()`" + ` to get multiple rows:

` + "```go" + `
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

type User struct {
    ID    int
    Name  string
    Email string
}

func main() {
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()

    rows, err := db.Query("SELECT id, name, email FROM users WHERE active = $1", true)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            panic(err)
        }
        users = append(users, user)
    }

    if err = rows.Err(); err != nil {
        panic(err)
    }

    for _, u := range users {
        fmt.Printf("User: %+v\n", u)
    }
}
` + "```" + `

## Executing Queries (INSERT, UPDATE, DELETE)

Use ` + "`Exec()`" + ` for statements that don't return rows:

` + "```go" + `
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()

    // INSERT
    result, err := db.Exec(
        "INSERT INTO users (name, email) VALUES ($1, $2)",
        "Charlie",
        "charlie@example.com",
    )
    if err != nil {
        panic(err)
    }

    lastID, _ := result.LastInsertId()
    rowsAffected, _ := result.RowsAffected()

    fmt.Printf("Inserted ID: %d, Rows affected: %d\n", lastID, rowsAffected)
}
` + "```" + `

## Prepared Statements

Prepared statements protect against SQL injection and improve performance:

` + "```go" + `
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()

    // Prepare statement
    stmt, err := db.Prepare("SELECT id, name FROM users WHERE id = $1")
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    // Reuse statement multiple times
    var id int
    var name string
    stmt.QueryRow(1).Scan(&id, &name)
    stmt.QueryRow(2).Scan(&id, &name)
    stmt.QueryRow(3).Scan(&id, &name)
}
` + "```" + `

## Transactions

Transactions group multiple statements into atomic operations:

` + "```go" + `
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func transfer(db *sql.DB, from, to int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    // Deduct from source
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount, from,
    )
    if err != nil {
        tx.Rollback()
        return err
    }

    // Add to destination
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount, to,
    )
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

func main() {
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()

    err := transfer(db, 1, 2, 100.00)
    if err != nil {
        fmt.Println("Transfer failed:", err)
    } else {
        fmt.Println("Transfer successful")
    }
}
` + "```" + `

## Connection Pooling

The ` + "`sql.DB`" + ` type manages a pool of connections automatically:

` + "```go" + `
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "time"
)

func main() {
    db, _ := sql.Open("postgres", dsn)

    // Configure connection pool
    db.SetMaxOpenConns(25)          // Max concurrent connections
    db.SetMaxIdleConns(5)           // Max idle connections
    db.SetConnMaxLifetime(5 * time.Minute)  // Max lifetime per connection

    defer db.Close()
}
` + "```" + `

## Error Handling

Common errors when working with databases:

` + "```go" + `
import "errors"

// Check for specific errors
if err == sql.ErrNoRows {
    // No rows in result set
}

// Check for connection errors
if err != nil {
    // Could be network error, auth error, syntax error, etc.
}

// Use errors.Is for wrapped errors
var nf *sql.ErrNoRows
if errors.As(err, &nf) {
    // No rows found
}
` + "```" + `

## Best Practices

1. **Always use prepared statements**: Prevents SQL injection
2. **Always defer Close()**: Ensures resources are cleaned up
3. **Configure connection pool**: Set appropriate limits
4. **Handle errors properly**: Check all error returns
5. **Use transactions for multi-step operations**: Ensures atomicity
6. **Scan into pointers for NULLs**: Handle database NULL properly
7. **Close result sets**: Call rows.Close() after iterating
8. **Test database operations**: Use mock databases for testing
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "db-connection",
				Title:       "Database Connection",
				Description: "Connect to a database and verify the connection.",
				Requirements: []string{
					"Create database connection",
					"Test connection with Ping()",
					"Handle errors properly",
					"Close connection gracefully",
				},
				InitialCode: `package main

import (
	"database/sql"
	"fmt"
)

func main() {
	// TODO: Import database driver
	// TODO: Create connection string
	// TODO: Open connection
	// TODO: Test connection
	// TODO: Handle errors
}
`,
				Solution: `package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Using SQLite for simplicity
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database")
}
`,
			},
			{
				ID:          "create-table",
				Title:       "Create Table",
				Description: "Create a table in the database.",
				Requirements: []string{
					"Define table schema",
					"Execute CREATE TABLE statement",
					"Handle errors",
					"Verify table creation",
				},
				InitialCode: `package main

import (
	"database/sql"
	"log"
)

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	// TODO: Create table schema
	// TODO: Execute CREATE TABLE
	// TODO: Handle errors
}
`,
				Solution: `package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	schema := ` + "`" + `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER
	)` + "`" + `

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Table created successfully")
}
`,
			},
			{
				ID:          "insert-data",
				Title:       "Insert Data",
				Description: "Insert records into the database.",
				Requirements: []string{
					"Insert single row",
					"Get last inserted ID",
					"Handle errors",
					"Verify insertion",
				},
				InitialCode: `package main

import (
	"database/sql"
	"log"
)

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	// TODO: Insert data
	// TODO: Get last ID
	// TODO: Handle errors
	// TODO: Display result
}
`,
				Solution: `package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	result, err := db.Exec(
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		"Alice", "alice@example.com", 30,
	)
	if err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}

	id, _ := result.LastInsertId()
	rows, _ := result.RowsAffected()
	fmt.Printf("Inserted ID: %d, Rows: %d\n", id, rows)
}
`,
			},
			{
				ID:          "query-single",
				Title:       "Query Single Row",
				Description: "Retrieve a single record from the database.",
				Requirements: []string{
					"Query single row using QueryRow",
					"Scan values into variables",
					"Handle no rows error",
					"Display result",
				},
				InitialCode: `package main

import (
	"database/sql"
	"log"
)

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	// TODO: Query single row
	// TODO: Scan values
	// TODO: Handle errors
	// TODO: Display
}
`,
				Solution: `package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var user User
	err = db.QueryRow("SELECT id, name, email, age FROM users WHERE id = ?", 1).
		Scan(&user.ID, &user.Name, &user.Email, &user.Age)

	if err == sql.ErrNoRows {
		fmt.Println("No user found")
		return
	} else if err != nil {
		log.Fatalf("Query error: %v", err)
	}

	fmt.Printf("User: ID=%d, Name=%s, Email=%s, Age=%d\n",
		user.ID, user.Name, user.Email, user.Age)
}
`,
			},
			{
				ID:          "query-multiple",
				Title:       "Query Multiple Rows",
				Description: "Retrieve multiple records and iterate over results.",
				Requirements: []string{
					"Query multiple rows",
					"Use rows.Next() to iterate",
					"Scan each row",
					"Handle errors",
					"Display all results",
				},
				InitialCode: `package main

import (
	"database/sql"
	"log"
)

type User struct {
	ID    int
	Name  string
}

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	// TODO: Query multiple rows
	// TODO: Iterate with rows.Next()
	// TODO: Scan each row
	// TODO: Handle errors
	// TODO: Display results
}
`,
				Solution: `package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Fatalf("Scan error: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Rows error: %v", err)
	}

	for _, u := range users {
		fmt.Printf("ID: %d, Name: %s\n", u.ID, u.Name)
	}
}
`,
			},
			{
				ID:          "prepared-statement",
				Title:       "Prepared Statements",
				Description: "Use prepared statements for safe and efficient queries.",
				Requirements: []string{
					"Prepare statement",
					"Execute prepared statement multiple times",
					"Bind parameters safely",
					"Close statement",
				},
				InitialCode: `package main

import (
	"database/sql"
	"log"
)

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	// TODO: Prepare statement
	// TODO: Execute multiple times
	// TODO: Close statement
	// TODO: Display results
}
`,
				Solution: `package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT name, email FROM users WHERE id = ?")
	if err != nil {
		log.Fatalf("Prepare error: %v", err)
	}
	defer stmt.Close()

	for id := 1; id <= 3; id++ {
		var name, email string
		err := stmt.QueryRow(id).Scan(&name, &email)
		if err == sql.ErrNoRows {
			fmt.Printf("No user with ID %d\n", id)
		} else if err != nil {
			log.Fatalf("Query error: %v", err)
		} else {
			fmt.Printf("ID %d: %s (%s)\n", id, name, email)
		}
	}
}
`,
			},
			{
				ID:          "transaction",
				Title:       "Transactions",
				Description: "Implement atomic transactions with rollback on error.",
				Requirements: []string{
					"Begin transaction",
					"Execute multiple statements",
					"Rollback on error",
					"Commit on success",
				},
				InitialCode: `package main

import (
	"database/sql"
	"log"
)

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	// TODO: Begin transaction
	// TODO: Execute statements
	// TODO: Handle errors with rollback
	// TODO: Commit on success
}
`,
				Solution: `package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Begin error: %v", err)
	}

	_, err = tx.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		"Bob", "bob@example.com", 25)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Insert error: %v", err)
	}

	_, err = tx.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		"Charlie", "charlie@example.com", 35)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Insert error: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Commit error: %v", err)
	}

	fmt.Println("Transaction completed successfully")
}
`,
			},
			{
				ID:          "update-delete",
				Title:       "Update and Delete Operations",
				Description: "Update and delete records from the database.",
				Requirements: []string{
					"Update existing records",
					"Delete records by ID",
					"Use WHERE clause",
					"Handle errors",
					"Display affected rows",
				},
				InitialCode: `package main

import (
	"database/sql"
	"log"
)

func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()

	// TODO: Update record
	// TODO: Delete record
	// TODO: Display affected rows
	// TODO: Handle errors
}
`,
				Solution: `package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	result, err := db.Exec(
		"UPDATE users SET age = ? WHERE id = ?",
		31, 1,
	)
	if err != nil {
		log.Fatalf("Update error: %v", err)
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("Updated %d row(s)\n", rows)

	result, err = db.Exec(
		"DELETE FROM users WHERE id = ?",
		2,
	)
	if err != nil {
		log.Fatalf("Delete error: %v", err)
	}
	rows, _ = result.RowsAffected()
	fmt.Printf("Deleted %d row(s)\n", rows)
}
`,
			},
		},
	}
}

// GetLessonData returns lesson data by lesson ID (11-15)
func (s *curriculumService) getLessonData11To15(lessonID int) *domain.LessonDetail {
	switch lessonID {
	case 11:
		return s.getComprehensiveLessonData11()
	case 12:
		return s.getComprehensiveLessonData12()
	case 13:
		return s.getComprehensiveLessonData13()
	case 14:
		return s.getComprehensiveLessonData14()
	case 15:
		return s.getComprehensiveLessonData15()
	default:
		return nil
	}
}
