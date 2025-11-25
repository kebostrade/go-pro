.PHONY: help deps build run test clean docker-build docker-up docker-down proto

# Variables
SERVICES := api-gateway user-service product-service order-service
DOCKER_COMPOSE := docker-compose -f deployments/docker-compose.yml

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy

build: ## Build all services
	@echo "🔨 Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd services/$$service && go build -o ../../bin/$$service ./cmd/main.go && cd ../..; \
	done
	@echo "✓ Build complete"

build-%: ## Build specific service (e.g., make build-user-service)
	@echo "🔨 Building $*..."
	@cd services/$* && go build -o ../../bin/$* ./cmd/main.go
	@echo "✓ $* built successfully"

run-gateway: ## Run API Gateway
	@echo "🚀 Starting API Gateway..."
	@cd services/api-gateway && go run ./cmd/main.go

run-user: ## Run User Service
	@echo "🚀 Starting User Service..."
	@cd services/user-service && go run ./cmd/main.go

run-product: ## Run Product Service
	@echo "🚀 Starting Product Service..."
	@cd services/product-service && go run ./cmd/main.go

run-order: ## Run Order Service
	@echo "🚀 Starting Order Service..."
	@cd services/order-service && go run ./cmd/main.go

test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test -v -race -cover ./...

test-integration: ## Run integration tests
	@echo "🧪 Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...

bench: ## Run benchmarks
	@echo "📊 Running benchmarks..."
	@go test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "✨ Formatting code..."
	@gofmt -s -w .
	@goimports -w .

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -rf coverage.out
	@echo "✓ Clean complete"

proto: ## Generate protobuf code
	@echo "🔧 Generating protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/*.proto
	@echo "✓ Protobuf generation complete"

docker-build: ## Build Docker images for all services
	@echo "🐳 Building Docker images..."
	@$(DOCKER_COMPOSE) build
	@echo "✓ Docker images built"

docker-up: ## Start all services with Docker Compose
	@echo "🐳 Starting services..."
	@$(DOCKER_COMPOSE) up -d
	@echo "✓ Services started"
	@echo ""
	@echo "📡 Services available at:"
	@echo "  API Gateway:     http://localhost:8080"
	@echo "  User Service:    http://localhost:8081"
	@echo "  Product Service: http://localhost:8082"
	@echo "  Order Service:   http://localhost:8083"

docker-down: ## Stop all services
	@echo "🐳 Stopping services..."
	@$(DOCKER_COMPOSE) down
	@echo "✓ Services stopped"

docker-logs: ## View logs from all services
	@$(DOCKER_COMPOSE) logs -f

docker-ps: ## Show running containers
	@$(DOCKER_COMPOSE) ps

docker-clean: ## Remove all containers, volumes, and images
	@echo "🧹 Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) down -v --rmi all
	@echo "✓ Docker cleanup complete"

dev: ## Run in development mode with hot reload
	@echo "🔥 Starting development mode..."
	@air

install-tools: ## Install development tools
	@echo "🔧 Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/cosmtrek/air@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✓ Tools installed"

demo: ## Run demo scenario
	@echo "🎬 Running demo..."
	@./scripts/demo.sh

k8s-deploy: ## Deploy to Kubernetes
	@echo "☸️  Deploying to Kubernetes..."
	@kubectl apply -f deployments/k8s/
	@echo "✓ Deployed to Kubernetes"

k8s-delete: ## Delete from Kubernetes
	@echo "☸️  Deleting from Kubernetes..."
	@kubectl delete -f deployments/k8s/
	@echo "✓ Deleted from Kubernetes"

