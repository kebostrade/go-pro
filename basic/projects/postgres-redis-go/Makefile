.PHONY: help deps docker-up docker-down test clean run-pgx run-gorm run-redis run-combined

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	go mod download
	go mod tidy

docker-up: ## Start PostgreSQL and Redis with Docker Compose
	docker compose up -d
	@echo "Waiting for databases to be ready..."
	@sleep 5

docker-down: ## Stop Docker containers
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f

test: ## Run tests
	go test -v -race -cover ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# PostgreSQL Examples
run-pgx: ## Run pgx examples
	go run postgres/pgx/main.go

run-gorm: ## Run GORM examples
	go run postgres/gorm/main.go

run-migrations: ## Run database migrations
	go run postgres/migrations/main.go

# Redis Examples
run-redis-basic: ## Run basic Redis examples
	go run redis/basic/main.go

run-redis-pubsub: ## Run Redis pub/sub examples
	go run redis/pubsub/main.go

run-redis-patterns: ## Run Redis patterns examples
	go run redis/patterns/main.go

# Combined Examples
run-cache: ## Run cache-aside pattern example
	go run combined/cache/main.go

run-session: ## Run session management example
	go run combined/session/main.go

run-all: ## Run all examples
	@echo "Running PostgreSQL pgx example..."
	@make run-pgx
	@echo "\nRunning GORM example..."
	@make run-gorm
	@echo "\nRunning Redis basic example..."
	@make run-redis-basic
	@echo "\nRunning cache example..."
	@make run-cache

setup: docker-up deps ## Setup project (start Docker + install deps)
	@echo "Setup complete! Run 'make run-all' to test examples"

