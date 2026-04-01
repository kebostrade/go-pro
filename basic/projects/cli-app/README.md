# CLI Application Template with Go and Cobra

A production-ready CLI application project template using Go and the Cobra framework.

## Features

- **Cobra v1.8.0**: Industry-standard CLI framework
- **Subcommands**: Root, greet, and serve commands
- **Configuration**: YAML-based configuration with flag overrides
- **Thread-safe greeting service**: Using sync.Mutex
- **Docker Support**: Multi-stage build for small images
- **CI/CD**: GitHub Actions with test and lint

## Project Structure

```
basic/projects/cli-app/
├── cmd/cli/main.go              # Entry point
├── internal/
│   ├── commands/
│   │   ├── root.go            # Root command
│   │   ├── greet.go           # Greet subcommand
│   │   └── serve.go          # Serve subcommand
│   └── config/
│       └── config.go          # Configuration loading
├── pkg/greeting/
│   └── greeting.go           # Core greeting logic
├── config.yaml                 # Sample configuration
├── Dockerfile
├── .github/workflows/ci.yml
├── go.mod
└── Makefile
```

## Prerequisites

- Go 1.23 or later
- Docker (optional)
- Make

## Quick Start

### Local Development

```bash
# Build and run
go build -o cli-app ./cmd/cli
./cli-app --help

# Or run directly
go run ./cmd/cli --help
```

### Docker

```bash
# Build and run with Docker
docker build -t cli-app .
docker run cli-app --help
docker run cli-app greet --name Alice --times 2
```

## Commands

### Root Command

```bash
app --help
```

### Greet Command

```bash
# Default greeting
app greet

# Greet someone specific
app greet --name Alice

# Greet multiple times
app greet --name Bob --times 3

# With custom config
app greet --name Charlie --config myconfig.yaml
```

### Serve Command

```bash
# Start server on default port 8080
app serve

# Start server on custom port
app serve --port 9090

# Bind to specific host
app serve --host 0.0.0.0 --port 3000
```

## Configuration

Configuration can be provided via YAML file:

```yaml
greeting:
  default_name: "World"
  default_times: 1
server:
  port: "8080"
  host: "localhost"
```

CLI flags override config file values.

## Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Make Commands

```bash
make build    # Build the binary
make run      # Run greet command
make test     # Run tests
make lint     # Run golangci-lint
make docker   # Build Docker image
make clean    # Clean up
```

## License

MIT
