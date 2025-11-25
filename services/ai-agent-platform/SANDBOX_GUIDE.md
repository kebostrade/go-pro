# 🔒 Code Execution Sandbox - Security Guide

## Overview

The Code Execution Sandbox provides a secure, isolated environment for running untrusted code with comprehensive resource limits and security policies.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│              Code Execution Request                  │
└────────────────────┬────────────────────────────────┘
                     │
        ┌────────────▼────────────┐
        │  Security Validator     │
        │  - Pattern Detection    │
        │  - Import Validation    │
        │  - Size Limits          │
        └────────────┬────────────┘
                     │
        ┌────────────▼────────────┐
        │  Resource Limits Mgr    │
        │  - Memory Limits        │
        │  - CPU Limits           │
        │  - Process Limits       │
        └────────────┬────────────┘
                     │
        ┌────────────▼────────────┐
        │   Docker Sandbox        │
        │  - Isolated Container   │
        │  - Read-only FS         │
        │  - No Network           │
        └────────────┬────────────┘
                     │
        ┌────────────▼────────────┐
        │   Execution Result      │
        │  - Output/Error         │
        │  - Metrics              │
        │  - Violations           │
        └─────────────────────────┘
```

## 🔐 Security Features

### 1. Container Isolation

Each code execution runs in a separate Docker container with:

- **Read-only root filesystem**: Prevents file modifications
- **No network access**: Blocks external connections
- **Dropped capabilities**: Removes all Linux capabilities
- **No new privileges**: Prevents privilege escalation
- **Resource limits**: CPU, memory, and process constraints

### 2. Code Validation

Before execution, code is validated for:

- **Dangerous patterns**: `eval()`, `exec()`, command execution
- **Blocked imports**: Network, file system, subprocess modules
- **Code size limits**: Maximum 100KB per file
- **Syntax validation**: Language-specific parsing

### 3. Resource Limits

Default limits per language:

| Language   | Memory | CPU Time | Processes |
|------------|--------|----------|-----------|
| Go         | 512MB  | 30s      | 10        |
| Python     | 512MB  | 30s      | 5         |
| JavaScript | 256MB  | 20s      | 5         |
| TypeScript | 256MB  | 20s      | 5         |
| Rust       | 1024MB | 60s      | 10        |
| Java       | 1024MB | 45s      | 10        |
| C++        | 512MB  | 30s      | 5         |
| C          | 256MB  | 20s      | 5         |

## 🛠️ Usage

### Basic Execution

```go
import (
    "context"
    "ai-agent-platform/internal/sandbox"
    "ai-agent-platform/pkg/types"
)

// Create sandbox
sandboxConfig := sandbox.SandboxConfig{
    ImagePrefix:    "coding-sandbox",
    NetworkMode:    "none",
    DefaultTimeout: 30 * time.Second,
}
dockerSandbox := sandbox.NewDockerSandbox(sandboxConfig)

// Execute code
result, err := dockerSandbox.Execute(context.Background(), types.ExecutionRequest{
    Code:     "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"Hello\") }",
    Language: "go",
    Timeout:  30,
})

if err != nil {
    log.Fatal(err)
}

fmt.Println("Output:", result.Output)
fmt.Println("Success:", result.Success)
```

### With Security Validation

```go
// Create security validator
validator := sandbox.NewSecurityValidator()
validator.RegisterPolicy("go", sandbox.GetGoPolicy())

// Validate code before execution
if err := validator.Validate("go", code); err != nil {
    log.Fatal("Security validation failed:", err)
}

// Execute if validation passes
result, err := dockerSandbox.Execute(ctx, request)
```

### With Custom Resource Limits

```go
// Create resource limits manager
limitsManager := sandbox.NewResourceLimitsManager()

// Set custom limits
customLimits := types.ResourceLimits{
    MaxMemoryMB:      1024,
    MaxCPUTime:       60,
    MaxProcesses:     20,
    NetworkAccess:    false,
    FileSystemAccess: false,
}

limitsManager.SetCustomLimits("go", customLimits)

// Get limits for execution
limits := limitsManager.GetLimits("go")

// Execute with custom limits
request.ResourceLimits = &limits
result, err := dockerSandbox.Execute(ctx, request)
```

## 🚫 Blocked Operations

### Go

- `os/exec` - Command execution
- `syscall` - System calls
- `unsafe` - Unsafe operations
- `net/http`, `net` - Network access
- `os.RemoveAll`, `os.Remove` - File deletion

### Python

- `os`, `sys` - Operating system access
- `subprocess` - Process spawning
- `socket`, `urllib`, `requests` - Network access
- `eval()`, `exec()` - Dynamic code execution
- `__import__` - Dynamic imports

### JavaScript/TypeScript

- `child_process` - Process spawning
- `fs` - File system access
- `net`, `http`, `https` - Network access
- `eval()`, `Function()` - Dynamic code execution
- `vm` - Virtual machine module

## 📊 Resource Monitoring

### Execution Metrics

```go
result := &types.ExecutionResult{
    Success:       true,
    Output:        "Hello, World!",
    ExitCode:      0,
    ExecutionTime: 1250,      // milliseconds
    MemoryUsed:    45000000,  // bytes
    CPUTime:       800,       // milliseconds
}
```

### Checking Violations

```go
limitsManager := sandbox.NewResourceLimitsManager()
limits := limitsManager.GetLimits("go")

violations := limitsManager.CheckLimitsExceeded(result, limits)
if len(violations) > 0 {
    for _, violation := range violations {
        log.Println("Violation:", violation)
    }
}
```

## 🐳 Docker Setup

### Build Sandbox Images

```bash
cd services/ai-agent-platform/docker/sandbox
chmod +x build-images.sh
./build-images.sh
```

This builds:
- `coding-sandbox-go:latest`
- `coding-sandbox-python:latest`
- `coding-sandbox-node:latest`

### Test Images

```bash
# Test Go sandbox
docker run --rm coding-sandbox-go:latest go version

# Test Python sandbox
docker run --rm coding-sandbox-python:latest python3 --version

# Test Node.js sandbox
docker run --rm coding-sandbox-node:latest node --version
```

### Manual Build

```bash
# Build Go sandbox
docker build -t coding-sandbox-go:latest -f Dockerfile.go .

# Build Python sandbox
docker build -t coding-sandbox-python:latest -f Dockerfile.python .

# Build Node.js sandbox
docker build -t coding-sandbox-node:latest -f Dockerfile.node .
```

## 🔧 Configuration

### Security Policies

```go
// Custom Go policy
policy := &sandbox.SecurityPolicy{
    AllowedImports: []string{
        "fmt",
        "strings",
        "math",
    },
    BlockedImports: []string{
        "os/exec",
        "net/http",
    },
    BlockedPatterns: []sandbox.CodePattern{
        {
            Pattern:        `os\.Remove`,
            Description:    "File deletion",
            Severity:       "high",
            Recommendation: "File operations not allowed",
        },
    },
    MaxCodeSize:            50 * 1024,
    AllowNetworkAccess:     false,
    AllowFileSystemAccess:  false,
    AllowExternalCommands:  false,
}

validator.RegisterPolicy("go", policy)
```

### Resource Limits

```go
// Recommended limits based on complexity
limits := limitsManager.GetRecommendedLimits(
    "go",
    codeSize,    // bytes
    complexity,  // cyclomatic complexity
)

// Estimate cost
cost := limitsManager.EstimateCost(limits)
fmt.Printf("Estimated cost: $%.4f\n", cost)
```

## 🚨 Error Handling

### Common Errors

```go
// Timeout error
if result.ExitCode == 124 {
    log.Println("Execution timeout")
}

// Memory limit exceeded
if strings.Contains(result.Error, "out of memory") {
    log.Println("Memory limit exceeded")
}

// Security violation
if err != nil {
    if toolErr, ok := err.(*types.ToolError); ok {
        if toolErr.Code == "UNSAFE_CODE_PATTERN" {
            log.Println("Security violation:", toolErr.Message)
        }
    }
}
```

## 📈 Performance

### Benchmarks

| Operation | Time |
|-----------|------|
| Container startup | ~100ms |
| Go execution | ~500ms |
| Python execution | ~300ms |
| JavaScript execution | ~200ms |
| Security validation | ~10ms |

### Optimization Tips

1. **Reuse containers**: Keep containers warm for faster execution
2. **Batch operations**: Execute multiple code snippets together
3. **Cache images**: Pre-pull Docker images
4. **Limit concurrency**: Use semaphore to control concurrent executions

## 🔍 Debugging

### Enable Verbose Logging

```go
sandboxConfig := sandbox.SandboxConfig{
    EnableLogging: true,
}
```

### Inspect Container

```bash
# List running containers
docker ps

# View container logs
docker logs <container-id>

# Inspect container
docker inspect <container-id>
```

## 🎯 Best Practices

1. **Always validate code** before execution
2. **Use appropriate limits** for each language
3. **Monitor resource usage** and adjust limits
4. **Clean up** temporary files regularly
5. **Log all executions** for audit trail
6. **Rate limit** execution requests
7. **Implement quotas** per user/session

## 🚀 Production Deployment

### Kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: code-sandbox
spec:
  containers:
  - name: sandbox
    image: coding-sandbox-go:latest
    resources:
      limits:
        memory: "512Mi"
        cpu: "1000m"
      requests:
        memory: "256Mi"
        cpu: "500m"
    securityContext:
      runAsNonRoot: true
      runAsUser: 1000
      readOnlyRootFilesystem: true
      allowPrivilegeEscalation: false
```

### Docker Compose

```yaml
version: '3.8'
services:
  sandbox:
    image: coding-sandbox-go:latest
    read_only: true
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    networks:
      - none
    tmpfs:
      - /tmp:rw,noexec,nosuid,size=100m
```

## 📞 Support

For issues or questions:
- Check the documentation
- Review security policies
- Test with simple examples first
- Enable verbose logging for debugging

---

**Built with security in mind for safe code execution** 🔒

