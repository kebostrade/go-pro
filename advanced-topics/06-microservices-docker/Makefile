# Makefile for Microservices with Docker

.PHONY: help build up down restart logs clean test health examples

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Microservices with Docker - Available Commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

build: ## Build all Docker images
	@echo "$(BLUE)Building all services...$(NC)"
	docker-compose build
	@echo "$(GREEN)✓ Build complete$(NC)"

up: ## Start all services
	@echo "$(BLUE)Starting all services...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✓ Services started$(NC)"
	@echo ""
	@echo "API Gateway: http://localhost:8080"
	@echo "User Service: http://localhost:8001"
	@echo "Order Service: http://localhost:8002"
	@echo ""
	@echo "Check health with: make health"

down: ## Stop all services
	@echo "$(BLUE)Stopping all services...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

restart: ## Restart all services
	@echo "$(BLUE)Restarting all services...$(NC)"
	docker-compose restart
	@echo "$(GREEN)✓ Services restarted$(NC)"

logs: ## Show logs from all services (follow mode)
	docker-compose logs -f

logs-api: ## Show API Gateway logs
	docker-compose logs -f api-gateway

logs-users: ## Show User Service logs
	docker-compose logs -f service-a

logs-orders: ## Show Order Service logs
	docker-compose logs -f service-b

clean: ## Stop services and remove volumes
	@echo "$(RED)Stopping services and removing volumes...$(NC)"
	docker-compose down -v
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

rebuild: ## Rebuild and restart all services
	@echo "$(BLUE)Rebuilding all services...$(NC)"
	docker-compose up -d --build
	@echo "$(GREEN)✓ Services rebuilt and restarted$(NC)"

health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo ""
	@echo "API Gateway:"
	@curl -s http://localhost:8080/health | jq '.' 2>/dev/null || echo "  Unavailable"
	@echo ""
	@echo "User Service:"
	@curl -s http://localhost:8001/health | jq '.' 2>/dev/null || echo "  Unavailable"
	@echo ""
	@echo "Order Service:"
	@curl -s http://localhost:8002/health | jq '.' 2>/dev/null || echo "  Unavailable"

ps: ## Show running containers
	docker-compose ps

test: health ## Run health check tests
	@echo "$(BLUE)Running health tests...$(NC)"
	@curl -s http://localhost:8080/health | grep -q "healthy" && echo "$(GREEN)✓ All services healthy$(NC)" || echo "$(RED)✗ Some services unhealthy$(NC)"

examples: ## Run example script
	@echo "$(BLUE)Running examples...$(NC)"
	./examples.sh

# Service-specific commands
start-users: ## Start only User Service
	docker-compose up -d service-a

start-orders: ## Start only Order Service
	docker-compose up -d service-b

start-gateway: ## Start only API Gateway
	docker-compose up -d api-gateway

# Database commands
db-users: ## Connect to users database
	docker-compose exec users-db psql -U postgres -d users_db

db-orders: ## Connect to orders database
	docker-compose exec orders-db psql -U postgres -d orders_db

db-redis: ## Connect to Redis
	docker-compose exec redis redis-cli

# Development commands
dev: ## Start services with hot reload (requires air)
	@echo "$(BLUE)Starting development environment...$(NC)"
	@echo "Install air first: go install github.com/cosmtrek/air@latest"
	@for service in service-a service-b api-gateway; do \
		cd $$service && air & \
		cd ..; \
	done

fmt: ## Format all Go code
	@echo "$(BLUE)Formatting code...$(NC)"
	@for service in service-a service-b api-gateway; do \
		echo "  Formatting $$service..."; \
		cd $$service && go fmt ./... && cd ..; \
	done
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet on all services
	@echo "$(BLUE)Vetting code...$(NC)"
	@for service in service-a service-b api-gateway; do \
		echo "  Vetted $$service"; \
		cd $$service && go vet ./... && cd ..; \
	done
	@echo "$(GREEN)✓ Vet complete$(NC)"

deps: ## Download dependencies for all services
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@for service in service-a service-b api-gateway; do \
		echo "  Downloading $$service dependencies..."; \
		cd $$service && go mod download && cd ..; \
	done
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

tidy: ## Tidy go.mod for all services
	@echo "$(BLUE)Tidying modules...$(NC)"
	@for service in service-a service-b api-gateway; do \
		echo "  Tidying $$service..."; \
		cd $$service && go mod tidy && cd ..; \
	done
	@echo "$(GREEN)✓ Modules tidied$(NC)"

# API testing commands
test-create-user: ## Test creating a user
	@echo "$(BLUE)Testing: Create user$(NC)"
	curl -X POST http://localhost:8080/api/users \
		-H "Content-Type: application/json" \
		-d '{"name":"Test User","email":"test@example.com"}' | jq '.'

test-get-users: ## Test getting all users
	@echo "$(BLUE)Testing: Get all users$(NC)"
	curl http://localhost:8080/api/users | jq '.'

test-create-order: ## Test creating an order
	@echo "$(BLUE)Testing: Create order$(NC)"
	@echo "Replace USER_ID with actual user ID"
	curl -X POST http://localhost:8080/api/orders \
		-H "Content-Type: application/json" \
		-d '{"user_id":"USER_ID","items":[{"product":"Widget","quantity":2}],"total":99.99}' | jq '.'

test-get-orders: ## Test getting all orders
	@echo "$(BLUE)Testing: Get all orders$(NC)"
	curl http://localhost:8080/api/orders | jq '.'

# Monitoring commands
stats: ## Show container resource usage
	docker stats

top: ## Show running processes
	docker-compose top

# Backup commands
backup-db: ## Backup all databases
	@echo "$(BLUE)Backing up databases...$(NC)"
	@mkdir -p backups
	docker exec users-db pg_dump -U postgres users_db > backups/users_db_$(shell date +%Y%m%d_%H%M%S).sql
	docker exec orders-db pg_dump -U postgres orders_db > backups/orders_db_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✓ Databases backed up to backups/$(NC)"

restore-users: ## Restore users database (usage: make restore-users FILE=backup.sql)
	@echo "$(BLUE)Restoring users database...$(NC)"
	docker-compose exec -T users-db psql -U postgres users_db < backups/$(FILE)
	@echo "$(GREEN)✓ Users database restored$(NC)"

# Utilities
shell-users: ## Open shell in User Service container
	docker-compose exec service-a sh

shell-orders: ## Open shell in Order Service container
	docker-compose exec service-b sh

shell-gateway: ## Open shell in API Gateway container
	docker-compose exec api-gateway sh

# Documentation
docs: ## Generate documentation (if using godoc)
	@echo "$(BLUE)Starting godoc server...$(NC)"
	@echo "Access at: http://localhost:6060"
	godoc -http=:6060

# CI/CD simulation
ci: ## Simulate CI pipeline
	@echo "$(BLUE)Running CI pipeline...$(NC)"
	@echo "  1. Formatting..."
	make fmt
	@echo "  2. Vetting..."
	make vet
	@echo "  3. Building..."
	make build
	@echo "  4. Testing..."
	make test
	@echo "$(GREEN)✓ CI pipeline complete$(NC)"

# Info commands
info: ## Show system information
	@echo "$(BLUE)System Information:$(NC)"
	@echo ""
	@echo "Docker Version:"
	@docker --version
	@echo ""
	@echo "Docker Compose Version:"
	@docker-compose --version
	@echo ""
	@echo "Go Version:"
	@go version
	@echo ""
	@echo "Running Containers:"
	@docker-compose ps
	@echo ""
	@echo "Network:"
	@docker network inspect microservices_microservices 2>/dev/null || echo "  Network not created yet"

version: ## Show version information
	@echo "Microservices with Docker - v1.0.0"
	@echo ""
	@echo "Services:"
	@echo "  API Gateway:  1.0.0"
	@echo "  User Service: 1.0.0"
	@echo "  Order Service: 1.0.0"
