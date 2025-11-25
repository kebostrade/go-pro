.PHONY: help build run test clean docker-build docker-run deps lint fmt

# Variables
APP_NAME=finagent
VERSION?=0.1.0
BUILD_DIR=bin
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/agent-server

build-all: ## Build all services
	@echo "Building all services..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/agent-server ./cmd/agent-server
	$(GOBUILD) -o $(BUILD_DIR)/llm-gateway ./cmd/llm-gateway
	$(GOBUILD) -o $(BUILD_DIR)/workflow-engine ./cmd/workflow-engine

run: ## Run the application
	$(GOCMD) run ./cmd/agent-server/main.go

run-example: ## Run fraud detection example
	@echo "Running fraud detection example..."
	$(GOCMD) run ./examples/fraud_detection/main.go

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Format code
	$(GOFMT) ./...
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	goimports -w .

vet: ## Run go vet
	$(GOVET) ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Start all services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop all services
	docker-compose down

docker-compose-logs: ## View logs
	docker-compose logs -f

# Database commands
db-migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@which migrate > /dev/null || (echo "Install golang-migrate first: https://github.com/golang-migrate/migrate" && exit 1)
	migrate -path migrations -database "$(DATABASE_URL)" up

db-migrate-down: ## Rollback database migrations
	migrate -path migrations -database "$(DATABASE_URL)" down

db-migrate-create: ## Create a new migration (usage: make db-migrate-create NAME=migration_name)
	migrate create -ext sql -dir migrations -seq $(NAME)

# Development
dev: ## Run in development mode with hot reload
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/cosmtrek/air@latest

# Examples
example-fraud: ## Run fraud detection example
	$(GOCMD) run ./examples/fraud_detection/main.go

example-customer-support: ## Run customer support example
	$(GOCMD) run ./examples/customer_support/main.go

example-risk-assessment: ## Run risk assessment example
	$(GOCMD) run ./examples/risk_assessment/main.go

# Monitoring
prometheus: ## Start Prometheus
	docker run -d -p 9090:9090 -v $(PWD)/configs/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

jaeger: ## Start Jaeger
	docker run -d -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one:latest

grafana: ## Start Grafana
	docker run -d -p 3000:3000 grafana/grafana

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@which godoc > /dev/null || (echo "Installing godoc..." && go install golang.org/x/tools/cmd/godoc@latest)
	@echo "Documentation server running at http://localhost:6060"
	godoc -http=:6060

# Release
release: clean test build ## Build release version
	@echo "Creating release $(VERSION)..."
	@mkdir -p releases
	@tar -czf releases/$(APP_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) .

.DEFAULT_GOAL := help

