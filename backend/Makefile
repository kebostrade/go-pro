.PHONY: help test test-unit test-integration test-coverage test-race test-bench clean build run lint fmt vet

# Variables
BINARY_NAME=go-pro-backend
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
COVERAGE_DIR=./coverage

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Test parameters
TEST_FLAGS=-v -race
INTEGRATION_FLAGS=-tags=integration
COVERAGE_FLAGS=-coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies downloaded"

test: ## Run all tests
	@echo "Running all tests..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) $(TEST_FLAGS) ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) $(TEST_FLAGS) -short ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) $(TEST_FLAGS) $(INTEGRATION_FLAGS) ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) $(TEST_FLAGS) $(COVERAGE_FLAGS) ./...
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.out"
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "HTML coverage report: $(COVERAGE_DIR)/coverage.html"

test-coverage-view: test-coverage ## Run tests with coverage and open HTML report
	@echo "Opening coverage report..."
	@open $(COVERAGE_DIR)/coverage.html || xdg-open $(COVERAGE_DIR)/coverage.html

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GOTEST) -race ./...

test-bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GOTEST) -bench=. -benchmem ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	$(GOTEST) -v ./...

lint: ## Run linter
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

check: fmt vet lint ## Run all checks (fmt, vet, lint)
	@echo "All checks passed"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME):latest

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	$(GOCMD) run $(MAIN_PATH) migrate up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	$(GOCMD) run $(MAIN_PATH) migrate down

migrate-status: ## Check migration status
	@echo "Checking migration status..."
	$(GOCMD) run $(MAIN_PATH) migrate status

db-reset: ## Reset database (down and up)
	@echo "Resetting database..."
	$(MAKE) migrate-down
	$(MAKE) migrate-up

# Development helpers
dev: ## Run in development mode with hot reload
	@echo "Starting development server..."
	@which air > /dev/null || (echo "air not installed. Installing..." && go install github.com/air-verse/air@latest)
	air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	go install github.com/stretchr/testify@latest
	@echo "Tools installed"

# CI/CD helpers
ci-test: ## Run tests for CI
	@echo "Running CI tests..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) $(TEST_FLAGS) $(COVERAGE_FLAGS) ./...

ci-lint: ## Run linter for CI
	@echo "Running CI linter..."
	golangci-lint run --timeout=5m ./...

ci: ci-lint ci-test ## Run all CI checks

# Security
security-scan: ## Run security scan
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "gosec not installed. Installing..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec ./...

# Performance
profile-cpu: ## Run CPU profiling
	@echo "Running CPU profiling..."
	$(GOTEST) -cpuprofile=cpu.prof -bench=. ./...
	$(GOCMD) tool pprof cpu.prof

profile-mem: ## Run memory profiling
	@echo "Running memory profiling..."
	$(GOTEST) -memprofile=mem.prof -bench=. ./...
	$(GOCMD) tool pprof mem.prof

# Database
db-shell: ## Open database shell
	@echo "Opening database shell..."
	docker-compose exec postgres psql -U gopro -d gopro

db-backup: ## Backup database
	@echo "Backing up database..."
	docker-compose exec postgres pg_dump -U gopro gopro > backup_$$(date +%Y%m%d_%H%M%S).sql

db-restore: ## Restore database from backup (usage: make db-restore FILE=backup.sql)
	@echo "Restoring database from $(FILE)..."
	docker-compose exec -T postgres psql -U gopro gopro < $(FILE)

# Default target
.DEFAULT_GOAL := help
