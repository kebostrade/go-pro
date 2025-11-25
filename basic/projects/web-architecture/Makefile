.PHONY: help build run test clean docker-build docker-run migrate-up migrate-down

# Default target
help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║                                                              ║"
	@echo "║           🏗️  Web Architecture with Go                      ║"
	@echo "║                                                              ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available targets:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback database migrations"

# Build the application
build:
	@echo "🔨 Building application..."
	@go build -o bin/server ./cmd/server
	@echo "✅ Build complete: bin/server"

# Run the application
run:
	@echo "🚀 Starting server..."
	@go run ./cmd/server

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t web-architecture:latest .
	@echo "✅ Docker image built"

# Run with Docker Compose
docker-run:
	@echo "🐳 Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "✅ Services started"

# Database migrations
migrate-up:
	@echo "📊 Running migrations..."
	@psql $(DATABASE_URL) -f migrations/001_init.sql
	@echo "✅ Migrations complete"

migrate-down:
	@echo "📊 Rolling back migrations..."
	@psql $(DATABASE_URL) -f migrations/001_init_down.sql
	@echo "✅ Rollback complete"

