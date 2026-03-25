// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package executor_test

import (
	"context"
	"os/exec"
	"testing"
	"time"

	executor "go-pro-backend/internal/executor"
	"go-pro-backend/internal/service"
)

func TestDockerExecutor_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker integration test in short mode")
	}

	exec := executor.NewDockerExecutor()

	code := `package main
import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`

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

	result, err := exec.ExecuteCode(context.Background(), req)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	t.Logf("Execution completed: Passed=%v, Score=%d%%", result.Passed, result.Score)
}

func TestDockerExecutor_WithInput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker integration test in short mode")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available")
	}

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

	result, err := exec.ExecuteCode(context.Background(), req)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	t.Logf("Passed: %v, Score: %d%%", result.Passed, result.Score)
}

func TestDockerExecutor_ValidationError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker integration test in short mode")
	}

	exec := executor.NewDockerExecutor()

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

	result, err := exec.ExecuteCode(context.Background(), req)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}

	if result.Error != nil {
		t.Logf("Validation failed: %v", result.Error)
	}
}
