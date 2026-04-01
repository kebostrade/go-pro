# gRPC Service Makefile

PROTO_DIR := proto
SERVER_DIR := cmd/server
CLIENT_DIR := cmd/client

# Protoc settings
PROTO_FILES := $(PROTO_DIR)/user.proto

# Colors
COLOR_RESET := \033[0m
COLOR_GREEN := \033[32m
COLOR_BLUE := \033[34m

.PHONY: all generate proto server client test lint clean help

# Default target
all: generate

## generate: Generate Go code from .proto files
generate: proto
	@echo "$(COLOR_BLUE)Generating Go code from .proto files...$(COLOR_RESET)"
	@cd $(PROTO_DIR) && \
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		user.proto
	@echo "$(COLOR_GREEN)✓ Generated Go code from .proto files$(COLOR_RESET)"

## proto: Alias for generate
proto: generate

## server: Run gRPC server
server:
	@echo "$(COLOR_BLUE)Starting gRPC server...$(COLOR_RESET)"
	@go run $(SERVER_DIR)/main.go

## client: Run gRPC client
client:
	@echo "$(COLOR_BLUE)Starting gRPC client...$(COLOR_RESET)"
	@go run $(CLIENT_DIR)/main.go

## test: Run tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@go test -v -race -cover ./...

## lint: Run linter
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@golangci-lint run ./... || true

## fmt: Format code
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@go fmt ./...
	@echo "$(COLOR_GREEN)✓ Formatted code$(COLOR_RESET)"

## clean: Clean generated files
clean:
	@echo "$(COLOR_BLUE)Cleaning generated files...$(COLOR_RESET)"
	@cd $(PROTO_DIR) && rm -f *.pb.go
	@echo "$(COLOR_GREEN)✓ Cleaned generated files$(COLOR_RESET)"

## deps: Install dependencies
deps:
	@echo "$(COLOR_BLUE)Installing dependencies...$(COLOR_RESET)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "$(COLOR_GREEN)✓ Installed dependencies$(COLOR_RESET)"

## grpcurl: Install grpcurl for debugging
grpcurl:
	@echo "$(COLOR_BLUE)Installing grpcurl...$(COLOR_RESET)"
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@echo "$(COLOR_GREEN)✓ Installed grpcurl$(COLOR_RESET)"

## reflect: List all gRPC services using reflection
reflect:
	@echo "$(COLOR_BLUE)Listing gRPC services...$(COLOR_RESET)"
	@grpcurl -plaintext localhost:50051 list

## help: Show this help message
help:
	@echo "gRPC Service - Available Commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
