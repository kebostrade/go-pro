// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package executor

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"go-pro-backend/internal/service"
)

func TestValidateCode(t *testing.T) {
	executor := &LocalExecutor{timeout: defaultTimeout}

	tests := []struct {
		name    string
		code    string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid code",
			code: `package main
import "fmt"
func main() {
	fmt.Println("Hello")
}`,
			wantErr: false,
		},
		{
			name:    "missing package main",
			code:    `import "fmt"\nfunc main() {}`,
			wantErr: true,
			errMsg:  "must contain 'package main'",
		},
		{
			name:    "missing func main",
			code:    `package main\nimport "fmt"`,
			wantErr: true,
			errMsg:  "must contain 'func main()'",
		},
		{
			name: "dangerous import - os",
			code: `package main
import "os"
func main() {}`,
			wantErr: true,
			errMsg:  "dangerous imports",
		},
		{
			name: "dangerous import - net",
			code: `package main
import "net"
func main() {}`,
			wantErr: true,
			errMsg:  "dangerous imports",
		},
		{
			name: "dangerous import - multiline",
			code: `package main
import (
	"fmt"
	"os"
)
func main() {}`,
			wantErr: true,
			errMsg:  "dangerous imports",
		},
		{
			name:    "code too large",
			code:    strings.Repeat("a", maxCodeSize+1),
			wantErr: true,
			errMsg:  "code too large",
		},
		{
			name: "safe imports",
			code: `package main
import (
	"fmt"
	"strings"
	"time"
)
func main() {
	fmt.Println("Hello")
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.validateCode(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateCode() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestExecuteCode_SimpleOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()
	ctx := context.Background()

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

	result, err := executor.ExecuteCode(ctx, req)
	if err != nil {
		t.Fatalf("ExecuteCode() error = %v", err)
	}

	if !result.Passed {
		t.Errorf("Expected test to pass, but it failed")
	}

	if result.Score != 100 {
		t.Errorf("Expected score 100, got %d", result.Score)
	}

	if len(result.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result.Results))
	}

	if !result.Results[0].Passed {
		t.Errorf("Test case failed: %s", result.Results[0].Error)
		t.Errorf("Expected: %q, Got: %q", result.Results[0].Expected, result.Results[0].Actual)
	}
}

func TestExecuteCode_WithInput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()
	ctx := context.Background()

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
				Name:     "with input",
				Input:    "Alice",
				Expected: "Hello, Alice!",
			},
			{
				Name:     "with different input",
				Input:    "Bob",
				Expected: "Hello, Bob!",
			},
		},
		Timeout: 5 * time.Second,
	}

	result, err := executor.ExecuteCode(ctx, req)
	if err != nil {
		t.Fatalf("ExecuteCode() error = %v", err)
	}

	if !result.Passed {
		t.Errorf("Expected all tests to pass")
		for i, r := range result.Results {
			if !r.Passed {
				t.Errorf("Test %d failed: %s (expected %q, got %q)", i, r.Error, r.Expected, r.Actual)
			}
		}
	}

	if result.Score != 100 {
		t.Errorf("Expected score 100, got %d", result.Score)
	}
}

func TestExecuteCode_CompilationError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()
	ctx := context.Background()

	code := `package main
import "fmt"
func main() {
	fmt.Println("Hello"
}` // Missing closing parenthesis

	req := &service.ExecuteRequest{
		Code:     code,
		Language: "go",
		TestCases: []service.TestCase{
			{
				Name:     "test",
				Input:    "",
				Expected: "Hello",
			},
		},
		Timeout: 5 * time.Second,
	}

	result, err := executor.ExecuteCode(ctx, req)
	if err != nil {
		t.Fatalf("ExecuteCode() error = %v", err)
	}

	if result.Passed {
		t.Errorf("Expected test to fail due to compilation error")
	}

	if result.Score != 0 {
		t.Errorf("Expected score 0, got %d", result.Score)
	}

	if len(result.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result.Results))
	}

	if result.Results[0].Error == "" {
		t.Errorf("Expected error message for compilation failure")
	}
}

func TestExecuteCode_MultipleTestCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()
	ctx := context.Background()

	code := `package main
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	a, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	b, _ := strconv.Atoi(scanner.Text())
	fmt.Println(a + b)
}`

	req := &service.ExecuteRequest{
		Code:     code,
		Language: "go",
		TestCases: []service.TestCase{
			{
				Name:     "test 1+2",
				Input:    "1\n2",
				Expected: "3",
			},
			{
				Name:     "test 5+7",
				Input:    "5\n7",
				Expected: "12",
			},
			{
				Name:     "test 0+0",
				Input:    "0\n0",
				Expected: "0",
			},
		},
		Timeout: 5 * time.Second,
	}

	result, err := executor.ExecuteCode(ctx, req)
	if err != nil {
		t.Fatalf("ExecuteCode() error = %v", err)
	}

	if !result.Passed {
		t.Errorf("Expected all tests to pass")
		for i, r := range result.Results {
			if !r.Passed {
				t.Errorf("Test %d (%s) failed: %s (expected %q, got %q)", i, r.TestName, r.Error, r.Expected, r.Actual)
			}
		}
	}

	if result.Score != 100 {
		t.Errorf("Expected score 100, got %d", result.Score)
	}
}

func TestExecuteCode_DangerousImports(t *testing.T) {
	executor := NewDockerExecutor()
	ctx := context.Background()

	code := `package main
import (
	"fmt"
	"os"
)
func main() {
	fmt.Println(os.Getenv("HOME"))
}`

	req := &service.ExecuteRequest{
		Code:     code,
		Language: "go",
		TestCases: []service.TestCase{
			{
				Name:     "test",
				Input:    "",
				Expected: "something",
			},
		},
		Timeout: 5 * time.Second,
	}

	result, err := executor.ExecuteCode(ctx, req)
	if err != nil {
		t.Fatalf("ExecuteCode() error = %v", err)
	}

	// Should fail validation
	if result.Error == nil {
		t.Errorf("Expected validation error for dangerous imports")
	}

	if !strings.Contains(result.Error.Error(), "dangerous imports") {
		t.Errorf("Expected 'dangerous imports' error, got: %v", result.Error)
	}
}

func TestExtractErrorMessage(t *testing.T) {
	executor := &LocalExecutor{timeout: defaultTimeout}

	tests := []struct {
		name   string
		stderr string
		want   string
	}{
		{
			name:   "syntax error",
			stderr: "./main.go:4:20: syntax error: unexpected newline, expecting comma or }",
			want:   "./main.go:4:20: syntax error: unexpected newline, expecting comma or }",
		},
		{
			name:   "undefined variable",
			stderr: "./main.go:5:2: undefined: x",
			want:   "./main.go:5:2: undefined: x",
		},
		{
			name: "with download messages",
			stderr: `go: downloading github.com/pkg/errors v0.9.1
./main.go:5:2: undefined: x`,
			want: "./main.go:5:2: undefined: x",
		},
		{
			name:   "panic message",
			stderr: "panic: runtime error: index out of range [1] with length 0",
			want:   "panic: runtime error: index out of range [1] with length 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executor.extractErrorMessage(tt.stderr)
			if !strings.Contains(got, tt.want) {
				t.Errorf("extractErrorMessage() = %v, want containing %v", got, tt.want)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	executor := &LocalExecutor{timeout: defaultTimeout}

	tests := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{
			name:    "timeout error",
			err:     fmt.Errorf("context deadline exceeded"),
			wantMsg: "timed out",
		},
		{
			name:    "syntax error",
			err:     fmt.Errorf("syntax error: unexpected newline"),
			wantMsg: "Compilation error",
		},
		{
			name:    "panic error",
			err:     fmt.Errorf("panic: runtime error"),
			wantMsg: "Runtime error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executor.formatError(tt.err)
			if !strings.Contains(got, tt.wantMsg) {
				t.Errorf("formatError() = %v, want containing %v", got, tt.wantMsg)
			}
		})
	}
}
