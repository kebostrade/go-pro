.PHONY: help build run test clean deps dev

help: ## Show this help message
	@echo '╔══════════════════════════════════════════════════════════════╗'
	@echo '║                                                              ║'
	@echo '║           💬 WebSocket Chat Application                     ║'
	@echo '║                                                              ║'
	@echo '╚══════════════════════════════════════════════════════════════╝'
	@echo ''
	@echo 'Available commands:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

build: deps ## Build the application
	@echo "🔨 Building application..."
	@go build -o bin/chat-server ./cmd/main.go
	@echo "✓ Build complete: bin/chat-server"

run: ## Run the server
	@echo "🚀 Starting WebSocket Chat Server..."
	@go run ./cmd/main.go

dev: ## Run with hot reload (requires air)
	@echo "🔥 Starting development server with hot reload..."
	@air

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "📊 Generating coverage report..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "✨ Formatting code..."
	@go fmt ./...
	@goimports -w .

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

demo: build ## Run demo with example clients
	@echo "🎬 Starting demo..."
	@./bin/chat-server &
	@sleep 2
	@echo "Opening example client..."
	@open http://localhost:8080 || xdg-open http://localhost:8080 || echo "Open http://localhost:8080 in your browser"

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t websocket-chat:latest .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 websocket-chat:latest

.DEFAULT_GOAL := help

