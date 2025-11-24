# Docker Executor

Secure Go code executor using Docker containers with comprehensive sandboxing.

## Features

- **Sandboxed Execution**: Docker-based isolation with security constraints
- **Resource Limits**: 128MB memory, 0.5 CPU, 5-second timeout
- **Security**: Read-only filesystem, no network, non-root user
- **Test Automation**: Execute code against multiple test cases
- **Error Handling**: Clear error messages for compilation, runtime, and timeout errors

## Security Constraints

The executor enforces strict security measures:

- **Execution Timeout**: 5 seconds maximum
- **Memory Limit**: 128MB
- **CPU Limit**: 0.5 cores
- **Filesystem**: Read-only (except `/tmp` with noexec)
- **Network**: Completely disabled
- **User**: Non-root execution (UID 1000)
- **Code Validation**: Blocks dangerous imports (os, net, syscall, unsafe, runtime/debug)

## Usage

### Basic Example

```go
package main

import (
    "context"
    "time"
    "backend/internal/executor"
    "backend/internal/service"
)

func main() {
    // Create executor
    exec := executor.NewDockerExecutor()

    // Prepare request
    req := &service.ExecuteRequest{
        Code: `package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
}`,
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
        panic(err)
    }

    // Check results
    if result.Passed {
        fmt.Printf("All tests passed! Score: %d%%\n", result.Score)
    } else {
        fmt.Printf("Some tests failed. Score: %d%%\n", result.Score)
        for _, tr := range result.Results {
            if !tr.Passed {
                fmt.Printf("  Failed: %s - %s\n", tr.TestName, tr.Error)
            }
        }
    }
}
```

### With Input/Output Testing

```go
req := &service.ExecuteRequest{
    Code: `package main
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
}`,
    Language: "go",
    TestCases: []service.TestCase{
        {
            Name:     "Alice",
            Input:    "Alice",
            Expected: "Hello, Alice!",
        },
        {
            Name:     "Bob",
            Input:    "Bob",
            Expected: "Hello, Bob!",
        },
    },
    Timeout: 5 * time.Second,
}
```

## Test Case Format

Test cases compare expected vs actual output:

```go
service.TestCase{
    Name:     "test name",       // Descriptive test name
    Input:    "input data",      // Optional stdin input
    Expected: "expected output", // Expected stdout output
}
```

Output comparison:
- Whitespace is trimmed from both expected and actual
- Exact string matching (case-sensitive)
- Newlines are preserved

## Result Structure

```go
type ExecuteResult struct {
    Passed        bool          // True if all tests passed
    Score         int           // Percentage (0-100)
    Results       []TestResult  // Individual test results
    ExecutionTime time.Duration // Total execution time
    Error         error         // Validation or execution error
}

type TestResult struct {
    TestName string // Test identifier
    Passed   bool   // Test passed/failed
    Expected string // Expected output
    Actual   string // Actual output
    Error    string // Error message if failed
}
```

## Error Types

### Validation Errors

Code validation failures (before execution):

```go
result.Error = "dangerous imports detected: os, net, syscall, unsafe, and runtime/debug are not allowed"
result.Error = "code too large: max 65536 bytes allowed"
result.Error = "code must contain 'package main'"
result.Error = "code must contain 'func main()'"
```

### Compilation Errors

Go compilation failures:

```go
testResult.Error = "Compilation error: syntax error: unexpected newline"
testResult.Error = "Compilation error: undefined: variableName"
```

### Runtime Errors

Code execution failures:

```go
testResult.Error = "Runtime error: panic: index out of range"
testResult.Error = "Runtime error: panic: runtime error: nil pointer dereference"
```

### Timeout Errors

Execution exceeded time limit:

```go
testResult.Error = "Code execution timed out (5s limit)"
```

## Docker Requirements

The executor requires Docker to be installed and accessible:

```bash
# Check Docker installation
docker --version

# Pull required image
docker pull golang:1.23-alpine

# Test Docker access
docker run --rm golang:1.23-alpine go version
```

## Testing

Run tests with Docker available:

```bash
# Run all tests
cd backend/internal/executor
go test -v

# Skip integration tests (no Docker required)
go test -v -short

# Run specific test
go test -v -run TestExecuteCode_SimpleOutput
```

## Integration with Exercise Service

The executor is integrated into the exercise submission flow:

```go
// In ExerciseService.SubmitExercise
func (s *exerciseService) SubmitExercise(ctx context.Context, exerciseID string, req *domain.SubmitExerciseRequest) (*domain.ExerciseSubmissionResult, error) {
    // Get exercise details
    exercise, err := s.repo.FindByID(ctx, exerciseID)
    if err != nil {
        return nil, err
    }

    // Prepare executor request
    execReq := &service.ExecuteRequest{
        Code:      req.Code,
        Language:  "go",
        TestCases: convertTestCases(exercise.TestCases),
        Timeout:   5 * time.Second,
    }

    // Execute code
    result, err := s.executor.ExecuteCode(ctx, execReq)
    if err != nil {
        return nil, err
    }

    // Return submission result
    return &domain.ExerciseSubmissionResult{
        Passed:  result.Passed,
        Score:   result.Score,
        Results: convertTestResults(result.Results),
    }, nil
}
```

## Security Considerations

### Blocked Capabilities

The executor prevents:
- File system access (read-only, except `/tmp`)
- Network access (no sockets, HTTP, etc.)
- System calls (limited syscall access)
- Process forking (resource limits prevent fork bombs)
- Memory exhaustion (128MB hard limit)
- CPU exhaustion (5s timeout + 0.5 CPU limit)

### Allowed Operations

Students can use:
- Standard input/output (`fmt`, `bufio`)
- String manipulation (`strings`, `strconv`)
- Math operations (`math`, `math/rand`)
- Time operations (`time`)
- Data structures (slices, maps, structs)
- Algorithms and logic

### Known Limitations

1. **No File I/O**: Cannot read/write files (except stdin/stdout)
2. **No Network**: Cannot make HTTP requests or open sockets
3. **No OS Access**: Cannot read environment variables or execute commands
4. **Short Execution**: 5-second maximum (prevents infinite loops)
5. **Limited Memory**: 128MB maximum (prevents memory bombs)

## Performance

Typical execution times:
- Simple output: ~200-500ms
- With input: ~300-600ms
- Multiple tests: ~400-800ms
- Compilation error: ~100-200ms

Performance factors:
- Docker container startup: ~100-200ms
- Go compilation: ~100-300ms
- Code execution: ~10-100ms
- Cleanup: ~10-50ms

## Switching from Mock to Docker Executor

To enable the Docker executor in production:

```go
// In internal/service/interfaces.go

// Option 1: Environment flag
func NewServices(repos *repository.Repositories, config *Config) (*Services, error) {
    var executorService ExecutorService
    if os.Getenv("USE_DOCKER_EXECUTOR") == "true" {
        executorService = executor.NewDockerExecutor()
    } else {
        executorService = NewMockExecutorService()
    }

    return &Services{
        // ... other services
        Executor: executorService,
    }, nil
}

// Option 2: Configuration-based
func NewServices(repos *repository.Repositories, config *Config) (*Services, error) {
    return &Services{
        // ... other services
        Executor: executor.NewDockerExecutor(), // Direct usage
    }, nil
}
```

## Troubleshooting

### "docker: command not found"

Docker is not installed or not in PATH:

```bash
# Install Docker
# Ubuntu/Debian
sudo apt-get install docker.io

# macOS
brew install docker

# Add user to docker group
sudo usermod -aG docker $USER
```

### "permission denied while trying to connect to Docker"

User lacks Docker permissions:

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in, or:
newgrp docker
```

### "timeout exceeded" errors

Code has infinite loops or is too slow:

```go
// Bad: Infinite loop
for {
    // Never terminates
}

// Good: Bounded loop
for i := 0; i < 1000000; i++ {
    // Terminates
}
```

### "dangerous imports detected"

Code uses blocked packages:

```go
// Bad: Network access
import "net"

// Bad: File access
import "os"

// Good: Safe operations
import "fmt"
import "strings"
```

## Future Enhancements

Potential improvements:
- [ ] Multi-language support (Python, JavaScript)
- [ ] Custom resource limits per exercise
- [ ] Execution caching for repeated submissions
- [ ] Detailed performance metrics
- [ ] Code quality analysis
- [ ] Memory usage profiling
- [ ] Test coverage reporting
