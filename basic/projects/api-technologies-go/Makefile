.PHONY: help setup deps clean docker-up docker-down docker-logs \
        run-rest-basic run-rest-chi run-rest-gin \
        run-grpc-unary run-grpc-streaming \
        run-graphql run-gateway run-all \
        proto-gen graphql-gen test lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: deps ## Setup project (install dependencies)
	@echo "✅ Project setup complete!"

deps: ## Download Go dependencies
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/ dist/ tmp/
	find . -name "*.pb.go" -delete
	find . -name "*_gen.go" -delete

# Docker commands
docker-up: ## Start Docker services
	@echo "🐳 Starting Docker services..."
	docker compose up -d
	@echo "✅ Services started!"
	@echo "📊 PostgreSQL: localhost:5432"
	@echo "🔴 Redis: localhost:6379"

docker-down: ## Stop Docker services
	@echo "🛑 Stopping Docker services..."
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f

# REST API examples
run-rest-basic: ## Run basic REST API example
	@echo "🚀 Running basic REST API..."
	cd rest/basic && go run main.go

run-rest-chi: ## Run Chi router REST API example
	@echo "🚀 Running Chi REST API..."
	cd rest/chi && go run main.go

run-rest-gin: ## Run Gin framework REST API example
	@echo "🚀 Running Gin REST API..."
	cd rest/gin && go run main.go

# gRPC examples
proto-gen: ## Generate protobuf code
	@echo "🔧 Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		grpc/proto/*.proto
	@echo "✅ Protobuf code generated!"

run-grpc-unary: ## Run gRPC unary example
	@echo "🚀 Running gRPC unary example..."
	cd grpc/unary && go run main.go

run-grpc-streaming: ## Run gRPC streaming example
	@echo "🚀 Running gRPC streaming example..."
	cd grpc/streaming && go run main.go

# GraphQL examples
graphql-gen: ## Generate GraphQL code
	@echo "🔧 Generating GraphQL code..."
	cd graphql/gqlgen && go run github.com/99designs/gqlgen generate
	@echo "✅ GraphQL code generated!"

run-graphql: ## Run GraphQL server
	@echo "🚀 Running GraphQL server..."
	cd graphql/gqlgen && go run main.go

# Combined examples
run-gateway: ## Run API Gateway (REST + gRPC + GraphQL)
	@echo "🚀 Running API Gateway..."
	cd combined/gateway && go run main.go

run-comparison: ## Run performance comparison
	@echo "📊 Running performance comparison..."
	cd combined/comparison && go run main.go

# Run all examples
run-all: ## Run all examples sequentially
	@echo "🚀 Running all examples..."
	@make run-rest-basic &
	@make run-rest-chi &
	@make run-grpc-unary &
	@make run-graphql &
	@echo "✅ All examples started!"

# Testing
test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v -race -cover ./...

# Linting
lint: ## Run linter
	@echo "🔍 Running linter..."
	golangci-lint run ./...

# Build
build: ## Build all binaries
	@echo "🔨 Building binaries..."
	@mkdir -p bin
	go build -o bin/rest-basic ./rest/basic
	go build -o bin/rest-chi ./rest/chi
	go build -o bin/rest-gin ./rest/gin
	go build -o bin/grpc-unary ./grpc/unary
	go build -o bin/grpc-streaming ./grpc/streaming
	go build -o bin/graphql ./graphql/gqlgen
	go build -o bin/gateway ./combined/gateway
	@echo "✅ Build complete! Binaries in bin/"

