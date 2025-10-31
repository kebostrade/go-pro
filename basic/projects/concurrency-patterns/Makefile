.PHONY: help deps test test-verbose test-coverage bench clean run-all fmt lint vet

# Default target
help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║                                                              ║"
	@echo "║           🔄 Concurrency Patterns in Go                     ║"
	@echo "║                                                              ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available targets:"
	@echo "  make deps              - Download dependencies"
	@echo "  make test              - Run all tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make bench             - Run benchmarks"
	@echo "  make run-all           - Run all examples"
	@echo "  make run-goroutines    - Run goroutine basics"
	@echo "  make run-channels      - Run channel examples"
	@echo "  make run-sync          - Run sync primitives"
	@echo "  make run-patterns      - Run concurrency patterns"
	@echo "  make run-worker-pool   - Run worker pool example"
	@echo "  make run-pipeline      - Run pipeline example"
	@echo "  make run-scraper       - Run web scraper example"
	@echo "  make fmt               - Format code"
	@echo "  make lint              - Run linter"
	@echo "  make vet               - Run go vet"
	@echo "  make clean             - Clean build artifacts"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Run all tests
test:
	@echo "🧪 Running tests..."
	@go test -race -v ./...

# Run tests with verbose output
test-verbose:
	@echo "🧪 Running tests (verbose)..."
	@go test -race -v -count=1 ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Run all examples
run-all:
	@echo "🚀 Running all examples..."
	@go run examples/all_patterns.go

# Run goroutine basics
run-goroutines:
	@echo "🔄 Running goroutine basics..."
	@go run basics/goroutines.go

# Run channel examples
run-channels:
	@echo "📡 Running channel examples..."
	@go run basics/channels.go

# Run sync primitives
run-sync:
	@echo "🔒 Running sync primitives..."
	@go run basics/sync_primitives.go

# Run concurrency patterns
run-patterns:
	@echo "🎯 Running concurrency patterns..."
	@go run patterns/worker_pool.go

# Run worker pool
run-worker-pool:
	@echo "👷 Running worker pool..."
	@go run examples/worker_pool_example.go

# Run pipeline
run-pipeline:
	@echo "⚙️  Running pipeline..."
	@go run examples/pipeline_example.go

# Run web scraper
run-scraper:
	@echo "🕷️  Running web scraper..."
	@go run examples/web_scraper.go

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...
	@echo "✅ Linting complete"

# Run go vet
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f coverage.out coverage.html
	@rm -rf bin/
	@echo "✅ Clean complete"

