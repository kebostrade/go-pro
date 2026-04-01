.PHONY: help test build run-nats run-publisher run-subscriber run-worker clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod

help: ## Display this help message
	@echo "NATS Events Makefile"
	@echo ""
	@echo "Usage:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-20s %s\n", $$1, $$2}'

test: ## Run tests
	$(GOTEST) -short -race ./...

build: ## Build all binaries
	$(GOBUILD) -o publisher ./cmd/publisher
	$(GOBUILD) -o subscriber ./cmd/subscriber
	$(GOBUILD) -o worker ./cmd/worker

vet: ## Run go vet
	$(GOVET) ./...

mod-tidy: ## Tidy go modules
	$(GOMOD) tidy

# Docker commands
docker-build: ## Build Docker image
	docker build -t nats-events:latest .

run-nats: ## Start NATS server with Docker Compose
	docker-compose up -d nats

stop-nats: ## Stop NATS server
	docker-compose down

# Run commands (requires NATS running)
run-publisher: ## Run publisher
	$(GOCMD) run ./cmd/publisher

run-subscriber: ## Run subscriber
	$(GOCMD) run ./cmd/subscriber

run-worker: ## Run worker
	$(GOCMD) run ./cmd/worker

# Demo: run all components
demo: run-nats
	@echo "Starting demo..."
	@echo "1. Run 'make run-subscriber' in another terminal"
	@echo "2. Run 'make run-publisher' in another terminal"
	@echo "3. Press Ctrl+C to stop"

clean: ## Clean build artifacts
	rm -f publisher subscriber worker
	docker rmi nats-events:latest 2>/dev/null || true
