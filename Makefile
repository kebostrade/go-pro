# ╔══════════════════════════════════════════════════════════════════════════════╗
# ║                      GO-PRO Learning Platform                               ║
# ╚══════════════════════════════════════════════════════════════════════════════╝
.DEFAULT_GOAL := help
.PHONY: help build test clean dev docker docker-dev docker-prod lint security coverage setup install-tools start-dev test-integration

# ─────────────────────────────────────────────────────────────────────────────────
# Variables
# ─────────────────────────────────────────────────────────────────────────────────
APP_NAME := go-pro-backend
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := 1.23

# Docker variables
DOCKER_REGISTRY := ghcr.io
DOCKER_ORG := your-org
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(APP_NAME)

# Build flags
LDFLAGS := -ldflags="-w -s -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.commit=$(GIT_COMMIT)"

# ─────────────────────────────────────────────────────────────────────────────────
# Colors & Formatting
# ─────────────────────────────────────────────────────────────────────────────────
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
BOLD := \033[1m
DIM := \033[2m
NC := \033[0m

# Icons
ICON_OK := ✓
ICON_FAIL := ✗
ICON_ARROW := →
ICON_ROCKET := 🚀
ICON_GEAR := ⚙
ICON_PACKAGE := 📦
ICON_TEST := 🧪
ICON_LINT := 🔍
ICON_SECURITY := 🔒
ICON_DOCKER := 🐳
ICON_DB := 🗄
ICON_DEPLOY := 🌐
ICON_CLEAN := 🧹
ICON_INFO := ℹ

# Box drawing helpers
define print_header
	@printf "\n"
	@printf "$(CYAN)╔══════════════════════════════════════════════════════════════╗$(NC)\n"
	@printf "$(CYAN)║$(NC) $(BOLD)%-60s$(NC) $(CYAN)║$(NC)\n" "$(1)"
	@printf "$(CYAN)╚══════════════════════════════════════════════════════════════╝$(NC)\n"
endef

define print_section
	@printf "\n$(BLUE)┌─$(NC) $(BOLD)$(1)$(NC)\n"
endef

define print_step
	@printf "$(DIM)│$(NC)  $(YELLOW)$(ICON_ARROW)$(NC) $(1)\n"
endef

define print_success
	@printf "$(DIM)│$(NC)  $(GREEN)$(ICON_OK)$(NC) $(1)\n"
endef

define print_error
	@printf "$(DIM)│$(NC)  $(RED)$(ICON_FAIL)$(NC) $(1)\n"
endef

define print_done
	@printf "$(BLUE)└─$(NC) $(GREEN)$(ICON_OK) Done!$(NC)\n\n"
endef

define print_info
	@printf "   $(DIM)$(1)$(NC)\n"
endef

##@ General Commands

help: ## Display this help message
	@printf "\n"
	@printf "$(CYAN)╔══════════════════════════════════════════════════════════════════════════╗$(NC)\n"
	@printf "$(CYAN)║$(NC)                                                                          $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)   $(BOLD)$(WHITE) ██████╗  ██████╗       ██████╗ ██████╗  ██████╗ $(NC)                  $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)   $(BOLD)$(WHITE)██╔════╝ ██╔═══██╗      ██╔══██╗██╔══██╗██╔═══██╗$(NC)                  $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)   $(BOLD)$(WHITE)██║  ███╗██║   ██║█████╗██████╔╝██████╔╝██║   ██║$(NC)                  $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)   $(BOLD)$(WHITE)██║   ██║██║   ██║╚════╝██╔═══╝ ██╔══██╗██║   ██║$(NC)                  $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)   $(BOLD)$(WHITE)╚██████╔╝╚██████╔╝      ██║     ██║  ██║╚██████╔╝$(NC)                  $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)   $(BOLD)$(WHITE) ╚═════╝  ╚═════╝       ╚═╝     ╚═╝  ╚═╝ ╚═════╝ $(NC)                  $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)                                                                          $(CYAN)║$(NC)\n"
	@printf "$(CYAN)║$(NC)   $(DIM)Go Learning Platform $(NC)$(PURPLE)v$(VERSION)$(NC)                                       $(CYAN)║$(NC)\n"
	@printf "$(CYAN)╚══════════════════════════════════════════════════════════════════════════╝$(NC)\n"
	@printf "\n"
	@printf "  $(BOLD)Usage:$(NC) make $(CYAN)<target>$(NC)\n\n"
	@awk 'BEGIN {FS = ":.*##"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "    $(CYAN)%-18s$(NC) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n  $(BOLD)$(YELLOW)%s$(NC)\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@printf "\n"

version: ## Display version information
	$(call print_header,Version Information)
	@printf "$(DIM)│$(NC)  $(BOLD)Version:$(NC)     $(GREEN)$(VERSION)$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Build Date:$(NC)  $(YELLOW)$(BUILD_DATE)$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Git Commit:$(NC)  $(PURPLE)$(GIT_COMMIT)$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Go Version:$(NC)  $(CYAN)$(GO_VERSION)$(NC)\n"
	$(call print_done)

##@ Development

setup: install-tools init-git ## Setup development environment
	$(call print_header,$(ICON_GEAR) Setup Complete)
	$(call print_success,Development environment ready!)
	$(call print_done)

start-dev: ## Start backend and frontend development servers
	$(call print_header,$(ICON_ROCKET) Starting Dev Servers)
	@./scripts/start-dev.sh

install-tools: ## Install development tools
	$(call print_header,$(ICON_GEAR) Installing Dev Tools)
	$(call print_step,Installing golangci-lint...)
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(call print_step,Installing gosec...)
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(call print_step,Installing air (hot reload)...)
	@go install github.com/air-verse/air@latest
	$(call print_step,Installing goimports...)
	@go install golang.org/x/tools/cmd/goimports@latest
	$(call print_step,Installing govulncheck...)
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	$(call print_step,Installing pre-commit...)
	@pip install pre-commit --quiet
	$(call print_success,All tools installed)
	$(call print_done)

init-git: ## Initialize git hooks
	$(call print_section,$(ICON_GEAR) Git Hooks)
	$(call print_step,Installing pre-commit hooks...)
	@pre-commit install >/dev/null 2>&1
	@pre-commit install --hook-type commit-msg >/dev/null 2>&1
	$(call print_success,Git hooks configured)
	$(call print_done)

dev: ## Start backend + frontend dev servers (parallel with hot reload)
	$(call print_header,$(ICON_ROCKET) Development Environment)
	@printf "$(DIM)│$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Local Services:$(NC)\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Backend    $(DIM)→$(NC) http://localhost:8080 (air hot reload)\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Frontend   $(DIM)→$(NC) http://localhost:3000 (bun + turbopack)\n"
	@printf "$(DIM)│$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Docker Services (run 'make docker-dev'):$(NC)\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) PostgreSQL    $(DIM)→$(NC) localhost:5432\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) Redis        $(DIM)→$(NC) localhost:6379\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) Prometheus   $(DIM)→$(NC) http://localhost:9090\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) Grafana      $(DIM)→$(NC) http://localhost:3001 (admin/admin)\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) Jaeger       $(DIM)→$(NC) http://localhost:16686\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) Kafka       $(DIM)→$(NC) localhost:9092\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) Adminer      $(DIM)→$(NC) http://localhost:8081\n"
	@printf "$(DIM)│$(NC)    $(CYAN)●$(NC) Mailhog     $(DIM)→$(NC) http://localhost:8025\n"
	@printf "$(DIM)│$(NC)\n"
	@printf "$(BLUE)└─$(NC) $(YELLOW)Press Ctrl+C to stop$(NC)\n\n"

	@AIR_BIN=$(HOME)/go/bin/air; \
	trap 'kill 0' INT; \
		if [ -x "$$AIR_BIN" ] || command -v air >/dev/null 2>&1; then \
			cd backend && $$AIR_BIN -c .air.toml & \
		else \
			cd backend && go run ./cmd/server & \
		fi; \
		cd frontend && bun run dev & \
		wait
		wait

dev-backend: ## Start backend only with hot reload
	$(call print_header,$(ICON_ROCKET) Backend Server)
	@printf "$(DIM)│$(NC)  $(CYAN)●$(NC) API $(DIM)→$(NC) $(GREEN)http://localhost:8080$(NC)\n"
	@printf "$(BLUE)└─$(NC) $(YELLOW)Hot reload enabled (air)$(NC)\n\n"
	@cd backend && air -c .air.toml

dev-docker: ## Start development environment with Docker
	$(call print_header,$(ICON_DOCKER) Docker Dev Environment)
	$(call print_step,Building and starting containers...)
	@docker compose -f docker-compose.dev.yml up --build

##@ Build Commands

deps: ## Download and verify dependencies
	$(call print_section,$(ICON_PACKAGE) Dependencies)
	$(call print_step,Downloading modules...)
	@cd backend && go mod download
	$(call print_step,Verifying checksums...)
	@cd backend && go mod verify
	$(call print_success,Dependencies ready)
	$(call print_done)

build: deps ## Build the application
	$(call print_header,$(ICON_PACKAGE) Building $(APP_NAME))
	$(call print_step,Compiling for linux/amd64...)
	@cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/server
	$(call print_success,Binary: backend/bin/$(APP_NAME))
	$(call print_done)

build-all: ## Build for multiple platforms
	$(call print_header,$(ICON_PACKAGE) Multi-Platform Build)
	@cd backend && mkdir -p bin
	$(call print_step,linux/amd64...)
	@cd backend && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 ./cmd/server
	$(call print_step,linux/arm64...)
	@cd backend && GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux-arm64 ./cmd/server
	$(call print_step,darwin/amd64...)
	@cd backend && GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 ./cmd/server
	$(call print_step,darwin/arm64...)
	@cd backend && GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 ./cmd/server
	$(call print_step,windows/amd64...)
	@cd backend && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-windows-amd64.exe ./cmd/server
	@printf "$(DIM)│$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Binaries:$(NC)\n"
	@printf "$(DIM)│$(NC)    $(DIM)├──$(NC) $(APP_NAME)-linux-amd64\n"
	@printf "$(DIM)│$(NC)    $(DIM)├──$(NC) $(APP_NAME)-linux-arm64\n"
	@printf "$(DIM)│$(NC)    $(DIM)├──$(NC) $(APP_NAME)-darwin-amd64\n"
	@printf "$(DIM)│$(NC)    $(DIM)├──$(NC) $(APP_NAME)-darwin-arm64\n"
	@printf "$(DIM)│$(NC)    $(DIM)└──$(NC) $(APP_NAME)-windows-amd64.exe\n"
	$(call print_done)

##@ Testing

test: ## Run unit tests
	$(call print_header,$(ICON_TEST) Running Tests)
	$(call print_step,Unit tests with race detection...)
	@cd backend && go test -v -race ./... && \
		printf "$(DIM)│$(NC)  $(GREEN)$(ICON_OK)$(NC) All tests passed\n" || \
		printf "$(DIM)│$(NC)  $(RED)$(ICON_FAIL)$(NC) Some tests failed\n"
	$(call print_done)

test-coverage: ## Run tests with coverage
	$(call print_header,$(ICON_TEST) Test Coverage)
	$(call print_step,Running tests with coverage profiling...)
	@cd backend && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(call print_step,Generating HTML report...)
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	$(call print_success,Report: backend/coverage.html)
	$(call print_done)

test-integration: ## Run integration tests
	$(call print_header,$(ICON_TEST) Integration Tests)
	$(call print_step,Running integration test suite...)
	@cd backend && go test -v -tags=integration ./...
	$(call print_done)

test-load: ## Run load tests
	$(call print_header,$(ICON_TEST) Load Tests)
	$(call print_step,Executing k6 load test...)
	@cd backend && k6 run scripts/load-test.js
	$(call print_done)

##@ Code Quality

lint: ## Run linter
	$(call print_header,$(ICON_LINT) Linting Code)
	$(call print_step,Running golangci-lint...)
	@cd backend && golangci-lint run --timeout=5m && \
		printf "$(DIM)│$(NC)  $(GREEN)$(ICON_OK)$(NC) No issues found\n" || \
		printf "$(DIM)│$(NC)  $(YELLOW)⚠$(NC)  Issues detected\n"
	$(call print_done)

lint-fix: ## Run linter and fix issues
	$(call print_header,$(ICON_LINT) Auto-fixing Lint Issues)
	$(call print_step,Running golangci-lint with auto-fix...)
	@cd backend && golangci-lint run --fix --timeout=5m
	$(call print_success,Auto-fix complete)
	$(call print_done)

fmt: ## Format Go code
	$(call print_section,$(ICON_LINT) Formatting)
	$(call print_step,Running gofmt...)
	@cd backend && gofmt -s -w .
	$(call print_step,Running goimports...)
	@cd backend && goimports -w .
	$(call print_success,Code formatted)
	$(call print_done)

vet: ## Run go vet
	$(call print_section,$(ICON_LINT) Go Vet)
	$(call print_step,Analyzing code...)
	@cd backend && go vet ./... && \
		printf "$(DIM)│$(NC)  $(GREEN)$(ICON_OK)$(NC) No issues\n" || \
		printf "$(DIM)│$(NC)  $(RED)$(ICON_FAIL)$(NC) Issues found\n"
	$(call print_done)

security: ## Run security scan
	$(call print_header,$(ICON_SECURITY) Security Scan)
	$(call print_step,Running gosec analysis...)
	@cd backend && gosec -quiet ./... && \
		printf "$(DIM)│$(NC)  $(GREEN)$(ICON_OK)$(NC) No vulnerabilities found\n" || \
		printf "$(DIM)│$(NC)  $(RED)$(ICON_FAIL)$(NC) Security issues detected\n"
	$(call print_done)

vuln-check: ## Check for known vulnerabilities
	$(call print_header,$(ICON_SECURITY) Vulnerability Check)
	$(call print_step,Scanning dependencies with govulncheck...)
	@cd backend && govulncheck ./...
	$(call print_done)

pre-commit: ## Run pre-commit hooks manually
	$(call print_section,$(ICON_GEAR) Pre-commit Hooks)
	$(call print_step,Running all hooks...)
	@pre-commit run --all-files
	$(call print_done)

quality: deps lint vet security test ## Run all quality checks
	$(call print_header,$(ICON_OK) Quality Gate Passed)

##@ Docker Commands

docker-build: ## Build Docker image
	$(call print_header,$(ICON_DOCKER) Building Docker Image)
	$(call print_step,Building $(DOCKER_IMAGE):$(VERSION)...)
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg VCS_REF=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		backend/
	$(call print_success,Image: $(DOCKER_IMAGE):$(VERSION))
	$(call print_done)

docker-build-multi: ## Build multi-platform Docker image
	$(call print_header,$(ICON_DOCKER) Multi-Platform Build)
	$(call print_step,Setting up buildx builder...)
	@docker buildx create --use --name multi-arch-builder --driver docker-container 2>/dev/null || true
	$(call print_step,Building linux/amd64 + linux/arm64...)
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg VCS_REF=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		--push \
		backend/
	$(call print_success,Multi-platform image pushed!)
	$(call print_done)

docker-run: docker-build ## Run Docker container locally
	$(call print_header,$(ICON_DOCKER) Running Container)
	@printf "$(DIM)│$(NC)  $(CYAN)●$(NC) API $(DIM)→$(NC) $(GREEN)http://localhost:8080$(NC)\n"
	@printf "$(BLUE)└─$(NC) $(YELLOW)Press Ctrl+C to stop$(NC)\n\n"
	@docker run --rm -p 8080:8080 --name $(APP_NAME) $(DOCKER_IMAGE):$(VERSION)

docker-dev: ## Start full development environment
	$(call print_header,$(ICON_DOCKER) Docker Dev Environment)
	$(call print_step,Starting containers...)
	@docker compose -f docker/docker-compose.dev.yml up --build -d
	@printf "$(DIM)│$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Services:$(NC)\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Backend API      $(DIM)→$(NC) http://localhost:8080\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Frontend        $(DIM)→$(NC) http://localhost:3000\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Adminer          $(DIM)→$(NC) http://localhost:8081\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Redis Commander  $(DIM)→$(NC) http://localhost:8082\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Prometheus       $(DIM)→$(NC) http://localhost:9090\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Grafana          $(DIM)→$(NC) http://localhost:3001\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Jaeger           $(DIM)→$(NC) http://localhost:16686\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Kafka UI         $(DIM)→$(NC) http://localhost:8083\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Elasticsearch    $(DIM)→$(NC) http://localhost:9200\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Kibana           $(DIM)→$(NC) http://localhost:5601\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) MinIO            $(DIM)→$(NC) http://localhost:9000 (console:9001)\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) RabbitMQ         $(DIM)→$(NC) http://localhost:15672\n"
	@printf "$(DIM)│$(NC)    $(GREEN)●$(NC) Mailhog          $(DIM)→$(NC) http://localhost:8025\n"
	$(call print_done)

docker-prod: ## Start production environment
	$(call print_header,$(ICON_DOCKER) Production Environment)
	$(call print_step,Starting production containers...)
	@docker compose -f docker-compose.prod.yml up -d
	$(call print_success,Production environment running)
	$(call print_done)

docker-stop: ## Stop all Docker containers
	$(call print_section,$(ICON_DOCKER) Stopping Containers)
	$(call print_step,Stopping dev environment...)
	@docker compose -f docker-compose.dev.yml down 2>/dev/null || true
	$(call print_step,Stopping prod environment...)
	@docker compose -f docker-compose.prod.yml down 2>/dev/null || true
	$(call print_success,All containers stopped)
	$(call print_done)

docker-clean: ## Clean up Docker resources
	$(call print_section,$(ICON_CLEAN) Docker Cleanup)
	$(call print_step,Pruning unused containers...)
	@docker system prune -f >/dev/null 2>&1
	$(call print_step,Pruning unused volumes...)
	@docker volume prune -f >/dev/null 2>&1
	$(call print_success,Docker resources cleaned)
	$(call print_done)

##@ Database

db-migrate: ## Run database migrations
	$(call print_header,$(ICON_DB) Database Migrations)
	$(call print_step,Running migrations...)
	@printf "$(DIM)│$(NC)  $(YELLOW)⚠$(NC)  Not implemented yet\n"
	$(call print_done)

db-seed: ## Seed database with test data
	$(call print_section,$(ICON_DB) Database Seeding)
	$(call print_step,Seeding test data...)
	@printf "$(DIM)│$(NC)  $(YELLOW)⚠$(NC)  Not implemented yet\n"
	$(call print_done)

##@ Deployment

deploy-staging: ## Deploy to staging environment
	$(call print_header,$(ICON_DEPLOY) Staging Deployment)
	$(call print_step,Deploying to staging...)
	@printf "$(DIM)│$(NC)  $(YELLOW)⚠$(NC)  Not configured yet\n"
	$(call print_done)

deploy-prod: ## Deploy to production environment
	$(call print_header,$(ICON_DEPLOY) Production Deployment)
	@printf "$(DIM)│$(NC)  $(RED)⚠  PRODUCTION DEPLOYMENT$(NC)\n"
	$(call print_step,Deploying to production...)
	@printf "$(DIM)│$(NC)  $(YELLOW)⚠$(NC)  Not configured yet\n"
	$(call print_done)

##@ Maintenance

clean: ## Clean build artifacts and cache
	$(call print_header,$(ICON_CLEAN) Cleaning Up)
	$(call print_step,Removing build artifacts...)
	@cd backend && rm -rf bin/ tmp/ coverage.out coverage.html 2>/dev/null || true
	$(call print_step,Clearing Go cache...)
	@cd backend && go clean -cache -testcache -modcache 2>/dev/null || true
	$(call print_success,Workspace cleaned)
	$(call print_done)

logs: ## Show application logs
	$(call print_section,$(ICON_INFO) Dev Logs)
	@docker compose -f docker-compose.dev.yml logs -f go-pro-backend

logs-prod: ## Show production logs
	$(call print_section,$(ICON_INFO) Production Logs)
	@docker compose -f docker-compose.prod.yml logs -f go-pro-backend

backup: ## Create backup of important data
	$(call print_header,$(ICON_DB) Backup)
	$(call print_step,Creating backup...)
	@printf "$(DIM)│$(NC)  $(YELLOW)⚠$(NC)  Not implemented yet\n"
	$(call print_done)

health: ## Check application health
	$(call print_section,$(ICON_INFO) Health Check)
	@printf "$(DIM)│$(NC)  $(BOLD)Endpoint:$(NC) http://localhost:8080/api/v1/health\n"
	@printf "$(DIM)│$(NC)\n"
	@curl -s http://localhost:8080/api/v1/health | jq '.' 2>/dev/null && \
		printf "$(DIM)│$(NC)  $(GREEN)$(ICON_OK)$(NC) Application healthy\n" || \
		printf "$(DIM)│$(NC)  $(RED)$(ICON_FAIL)$(NC) Application not running or unhealthy\n"
	$(call print_done)

##@ CI/CD

ci-setup: install-tools ## Setup CI environment
	$(call print_header,$(ICON_GEAR) CI Setup Complete)

ci-test: deps lint vet security test-coverage ## Run CI test pipeline
	$(call print_header,$(ICON_TEST) CI Test Pipeline Complete)

ci-build: ci-test build docker-build ## Run CI build pipeline
	$(call print_header,$(ICON_PACKAGE) CI Build Pipeline Complete)

release: ## Create a new release
	$(call print_header,$(ICON_ROCKET) Creating Release)
	$(call print_step,Tagging v$(VERSION)...)
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	$(call print_step,Pushing tag to origin...)
	@git push origin v$(VERSION)
	@printf "$(DIM)│$(NC)\n"
	@printf "$(DIM)│$(NC)  $(BOLD)Release:$(NC) $(GREEN)v$(VERSION)$(NC)\n"
	$(call print_done)