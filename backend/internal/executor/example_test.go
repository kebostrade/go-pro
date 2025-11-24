// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package executor_test

import (
	"context"
	"fmt"
	"time"

	"go-pro-backend/internal/executor"
	"go-pro-backend/internal/service"
)

// Example demonstrates basic usage of the Docker executor.
// Note: This example requires Docker to be installed and running.
func Example() {
	// Create executor
	exec := executor.NewDockerExecutor()

	// Prepare code to execute
	code := `package main
import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`

	// Create execution request with test cases
	req := &service.ExecuteRequest{
		Code:     code,
		Language: "go",
		TestCases: []service.TestCase{
			{
				Name:     "simple output",
				Input:    "",
				Expected: "Hello, World!",
			},
		},
		Timeout: 5 * time.Second,
	}

	// Execute code
	result, err := exec.ExecuteCode(context.Background(), req)
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	// Display results (note: actual output depends on Docker availability)
	fmt.Printf("Execution completed: Passed=%v, Score=%d%%\n", result.Passed, result.Score)

	// Unordered output:
	// Execution completed: Passed=true, Score=100%
}

// Example_withInput demonstrates code execution with stdin input.
// Note: This example requires Docker to be installed and running.
func Example_withInput() {
	exec := executor.NewDockerExecutor()

	code := `package main
import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()
	fmt.Printf("Hello, %s!", name)
}`

	req := &service.ExecuteRequest{
		Code:     code,
		Language: "go",
		TestCases: []service.TestCase{
			{
				Name:     "greet Alice",
				Input:    "Alice",
				Expected: "Hello, Alice!",
			},
			{
				Name:     "greet Bob",
				Input:    "Bob",
				Expected: "Hello, Bob!",
			},
		},
		Timeout: 5 * time.Second,
	}

	result, _ := exec.ExecuteCode(context.Background(), req)

	fmt.Printf("Passed: %v, Score: %d%%\n", result.Passed, result.Score)
	// Unordered output:
	// Passed: true, Score: 100%
}

// Example_validationError demonstrates code validation failures.
func Example_validationError() {
	exec := executor.NewDockerExecutor()

	// Code with dangerous import
	code := `package main
import "os"

func main() {
	os.Exit(0)
}`

	req := &service.ExecuteRequest{
		Code:     code,
		Language: "go",
		TestCases: []service.TestCase{
			{
				Name:     "test",
				Input:    "",
				Expected: "",
			},
		},
		Timeout: 5 * time.Second,
	}

	result, _ := exec.ExecuteCode(context.Background(), req)

	if result.Error != nil {
		fmt.Printf("Validation failed: %v\n", result.Error)
	}

	// Output:
	// Validation failed: dangerous imports detected: os, net, syscall, unsafe, and runtime/debug are not allowed
}
