.PHONY: help test test-verbose test-coverage bench run-examples clean

help: ## Show this help message
	@echo '📚 Design Patterns in Go - Available Commands'
	@echo ''
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies ready"

test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test ./... -v

test-verbose: ## Run tests with verbose output
	@echo "🧪 Running tests (verbose)..."
	@go test ./... -v -race

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@go test ./... -bench=. -benchmem

run-singleton: ## Run singleton example
	@echo "🔹 Running Singleton Pattern..."
	@go run examples/singleton_example.go

run-factory: ## Run factory example
	@echo "🔹 Running Factory Pattern..."
	@go run examples/factory_example.go

run-builder: ## Run builder example
	@echo "🔹 Running Builder Pattern..."
	@go run examples/builder_example.go

run-strategy: ## Run strategy example
	@echo "🔹 Running Strategy Pattern..."
	@go run examples/strategy_example.go

run-observer: ## Run observer example
	@echo "🔹 Running Observer Pattern..."
	@go run examples/observer_example.go

run-all-examples: ## Run all pattern examples
	@echo "🎯 Running all examples..."
	@make run-singleton
	@make run-factory
	@make run-builder
	@make run-strategy
	@make run-observer

fmt: ## Format code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./...
	@echo "✓ Linting complete"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean
	@echo "✓ Clean complete"

.DEFAULT_GOAL := help

