# GO-PRO Learning Platform - Makefile
.DEFAULT_GOAL := help
.PHONY: help build test clean dev docker docker-dev docker-prod lint security coverage setup install-tools start-dev test-integration

# Variables
APP_NAME := go-pro-backend
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := 1.21

# Docker variables
DOCKER_REGISTRY := ghcr.io
DOCKER_ORG := your-org
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(APP_NAME)

# Build flags
LDFLAGS := -ldflags="-w -s -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.commit=$(GIT_COMMIT)"

# Color output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

##@ General Commands

help: ## Display this help message
	@echo "$(CYAN)GO-PRO Learning Platform Backend$(NC)"
	@echo "$(YELLOW)Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

version: ## Display version information
	@echo "$(GREEN)Version:$(NC) $(VERSION)"
	@echo "$(GREEN)Build Date:$(NC) $(BUILD_DATE)"
	@echo "$(GREEN)Git Commit:$(NC) $(GIT_COMMIT)"
	@echo "$(GREEN)Go Version:$(NC) $(GO_VERSION)"

##@ Development

setup: install-tools init-git ## Setup development environment
	@echo "$(GREEN)Development environment setup complete!$(NC)"

start-dev: ## Start backend and frontend development servers
	@echo "$(YELLOW)Starting development environment...$(NC)"
	@./scripts/start-dev.sh

test-integration: ## Run integration tests
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@./scripts/test-integration.sh

install-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/cosmtrek/air@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@pip install pre-commit
	@echo "$(GREEN)Development tools installed!$(NC)"

init-git: ## Initialize git hooks
	@echo "$(YELLOW)Setting up git hooks...$(NC)"
	@pre-commit install
	@pre-commit install --hook-type commit-msg
	@echo "$(GREEN)Git hooks installed!$(NC)"

dev: ## Start development server with hot reload
	@echo "$(YELLOW)Starting development server...$(NC)"
	@cd backend && air -c .air.toml

dev-docker: ## Start development environment with Docker
	@echo "$(YELLOW)Starting development environment with Docker...$(NC)"
	@docker-compose -f docker-compose.dev.yml up --build

##@ Build Commands

deps: ## Download and verify dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@cd backend && go mod download
	@cd backend && go mod verify

build: deps ## Build the application
	@echo "$(YELLOW)Building $(APP_NAME)...$(NC)"
	@cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/server
	@echo "$(GREEN)Build complete: backend/bin/$(APP_NAME)$(NC)"

build-all: ## Build for multiple platforms
	@echo "$(YELLOW)Building for multiple platforms...$(NC)"
	@cd backend && mkdir -p bin
	@cd backend && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 ./cmd/server
	@cd backend && GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux-arm64 ./cmd/server
	@cd backend && GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 ./cmd/server
	@cd backend && GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 ./cmd/server
	@cd backend && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-windows-amd64.exe ./cmd/server
	@echo "$(GREEN)Multi-platform build complete!$(NC)"

##@ Testing

test: ## Run unit tests
	@echo "$(YELLOW)Running unit tests...$(NC)"
	@cd backend && go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@cd backend && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: backend/coverage.html$(NC)"

test-integration: ## Run integration tests
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@cd backend && go test -v -tags=integration ./...

test-load: ## Run load tests
	@echo "$(YELLOW)Running load tests...$(NC)"
	@cd backend && k6 run scripts/load-test.js

##@ Code Quality

lint: ## Run linter
	@echo "$(YELLOW)Running golangci-lint...$(NC)"
	@cd backend && golangci-lint run --timeout=5m

lint-fix: ## Run linter and fix issues
	@echo "$(YELLOW)Running golangci-lint with fixes...$(NC)"
	@cd backend && golangci-lint run --fix --timeout=5m

fmt: ## Format Go code
	@echo "$(YELLOW)Formatting Go code...$(NC)"
	@cd backend && gofmt -s -w .
	@cd backend && goimports -w .

vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	@cd backend && go vet ./...

security: ## Run security scan
	@echo "$(YELLOW)Running security scan with gosec...$(NC)"
	@cd backend && gosec -quiet ./...

vuln-check: ## Check for known vulnerabilities
	@echo "$(YELLOW)Checking for vulnerabilities...$(NC)"
	@cd backend && govulncheck ./...

pre-commit: ## Run pre-commit hooks manually
	@echo "$(YELLOW)Running pre-commit hooks...$(NC)"
	@pre-commit run --all-files

quality: deps lint vet security test ## Run all quality checks

##@ Docker Commands

docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg VCS_REF=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		backend/
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(VERSION)$(NC)"

docker-build-multi: ## Build multi-platform Docker image
	@echo "$(YELLOW)Building multi-platform Docker image...$(NC)"
	@docker buildx create --use --name multi-arch-builder --driver docker-container 2>/dev/null || true
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg VCS_REF=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		--push \
		backend/
	@echo "$(GREEN)Multi-platform Docker image built and pushed!$(NC)"

docker-run: docker-build ## Run Docker container locally
	@echo "$(YELLOW)Running Docker container...$(NC)"
	@docker run --rm -p 8080:8080 --name $(APP_NAME) $(DOCKER_IMAGE):$(VERSION)

docker-dev: ## Start full development environment
	@echo "$(YELLOW)Starting development environment...$(NC)"
	@docker-compose -f docker-compose.dev.yml up --build -d
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo "$(CYAN)Services available:$(NC)"
	@echo "  - Backend API: http://localhost:8080"
	@echo "  - Adminer: http://localhost:8081"
	@echo "  - Redis Commander: http://localhost:8082"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana: http://localhost:3000"

docker-prod: ## Start production environment
	@echo "$(YELLOW)Starting production environment...$(NC)"
	@docker-compose -f docker-compose.prod.yml up -d
	@echo "$(GREEN)Production environment started!$(NC)"

docker-stop: ## Stop all Docker containers
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	@docker-compose -f docker-compose.dev.yml down 2>/dev/null || true
	@docker-compose -f docker-compose.prod.yml down 2>/dev/null || true

docker-clean: ## Clean up Docker resources
	@echo "$(YELLOW)Cleaning up Docker resources...$(NC)"
	@docker system prune -f
	@docker volume prune -f

##@ Database

db-migrate: ## Run database migrations
	@echo "$(YELLOW)Running database migrations...$(NC)"
	@echo "$(RED)Database migrations not implemented yet$(NC)"

db-seed: ## Seed database with test data
	@echo "$(YELLOW)Seeding database...$(NC)"
	@echo "$(RED)Database seeding not implemented yet$(NC)"

##@ Deployment

deploy-staging: ## Deploy to staging environment
	@echo "$(YELLOW)Deploying to staging...$(NC)"
	@echo "$(RED)Staging deployment not configured yet$(NC)"

deploy-prod: ## Deploy to production environment
	@echo "$(YELLOW)Deploying to production...$(NC)"
	@echo "$(RED)Production deployment not configured yet$(NC)"

##@ Maintenance

clean: ## Clean build artifacts and cache
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@cd backend && rm -rf bin/ tmp/ coverage.out coverage.html
	@cd backend && go clean -cache -testcache -modcache
	@echo "$(GREEN)Clean complete!$(NC)"

logs: ## Show application logs
	@echo "$(YELLOW)Showing application logs...$(NC)"
	@docker-compose -f docker-compose.dev.yml logs -f go-pro-backend

logs-prod: ## Show production logs
	@echo "$(YELLOW)Showing production logs...$(NC)"
	@docker-compose -f docker-compose.prod.yml logs -f go-pro-backend

backup: ## Create backup of important data
	@echo "$(YELLOW)Creating backup...$(NC)"
	@echo "$(RED)Backup functionality not implemented yet$(NC)"

health: ## Check application health
	@echo "$(YELLOW)Checking application health...$(NC)"
	@curl -s http://localhost:8080/api/v1/health | jq '.' || echo "$(RED)Application not running or unhealthy$(NC)"

##@ CI/CD

ci-setup: install-tools ## Setup CI environment
	@echo "$(GREEN)CI environment setup complete!$(NC)"

ci-test: deps lint vet security test-coverage ## Run CI test pipeline
	@echo "$(GREEN)CI test pipeline complete!$(NC)"

ci-build: ci-test build docker-build ## Run CI build pipeline
	@echo "$(GREEN)CI build pipeline complete!$(NC)"

release: ## Create a new release
	@echo "$(YELLOW)Creating release...$(NC)"
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git push origin v$(VERSION)
	@echo "$(GREEN)Release v$(VERSION) created!$(NC)"