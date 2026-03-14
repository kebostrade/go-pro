# gRPC Distributed System Examples Makefile

# Variables
PROTO_DIR := proto
SERVER_DIR := server
CLIENT_DIR := client
EXAMPLES_DIR := examples

# Protoc settings
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
GO_OUT := .

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m

.PHONY: all clean gen proto server client examples help test lint fmt

# Default target
all: gen

## generate: Generate Go code from .proto files
gen: proto
	@echo "$(COLOR_BLUE)Generating Go code from .proto files...$(COLOR_RESET)"
	@cd $(PROTO_DIR) && \
	protoc --go_out=$(GO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
		*.proto
	@echo "$(COLOR_GREEN)✓ Generated Go code from .proto files$(COLOR_RESET)"

## proto: Alias for gen
proto: gen

## server: Run gRPC server
server:
	@echo "$(COLOR_BLUE)Starting gRPC server...$(COLOR_RESET)"
	@cd $(SERVER_DIR) && go run main.go

## client: Run gRPC client
client:
	@echo "$(COLOR_BLUE)Starting gRPC client...$(COLOR_RESET)"
	@cd $(CLIENT_DIR) && go run main.go

## examples: Run all examples
examples:
	@echo "$(COLOR_BLUE)Running gRPC examples...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Running unary RPC example...$(COLOR_RESET)"
	@cd $(EXAMPLES_DIR) && go run 01-unary-rpc.go
	@echo "$(COLOR_YELLOW)Running server streaming example...$(COLOR_RESET)"
	@cd $(EXAMPLES_DIR) && go run 02-server-stream.go
	@echo "$(COLOR_YELLOW)Running client streaming example...$(COLOR_RESET)"
	@cd $(EXAMPLES_DIR) && go run 03-client-stream.go
	@echo "$(COLOR_YELLOW)Running bidirectional streaming example...$(COLOR_RESET)"
	@cd $(EXAMPLES_DIR) && go run 04-bidirectional-stream.go

## test: Run all tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@cd $(SERVER_DIR) && go test -v ./...
	@cd $(CLIENT_DIR) && go test -v ./...
	@cd $(EXAMPLES_DIR) && go test -v ./...

## lint: Run linter
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@cd $(SERVER_DIR) && golangci-lint run ./...
	@cd $(CLIENT_DIR) && golangci-lint run ./...
	@cd $(EXAMPLES_DIR) && golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@goimports -w $(PROTO_DIR)/*.pb.go
	@goimports -w $(SERVER_DIR)/*.go
	@goimports -w $(CLIENT_DIR)/*.go
	@goimports -w $(EXAMPLES_DIR)/*.go
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

## certs: Generate TLS certificates
certs:
	@echo "$(COLOR_BLUE)Generating TLS certificates...$(COLOR_RESET)"
	@mkdir -p certs
	@openssl genrsa -out certs/server.key 2048
	@openssl req -new -x509 -sha256 -key certs/server.key \
		-out certs/server.crt -days 3650 \
		-subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
	@echo "$(COLOR_GREEN)✓ Generated TLS certificates$(COLOR_RESET)"

## grpcurl: Install grpcurl for debugging
grpcurl:
	@echo "$(COLOR_BLUE)Installing grpcurl...$(COLOR_RESET)"
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@echo "$(COLOR_GREEN)✓ Installed grpcurl$(COLOR_RESET)"

## reflect: List all gRPC services using reflection
reflect:
	@echo "$(COLOR_BLUE)Listing gRPC services...$(COLOR_RESET)"
	@grpcurl -plaintext localhost:50051 list

## describe: Describe a service
describe:
	@echo "$(COLOR_BLUE)Describing UserService...$(COLOR_RESET)"
	@grpcurl -plaintext localhost:50051 describe UserService

## help: Show this help message
help:
	@echo "$(COLOR_BOLD)gRPC Distributed System Examples - Available Commands:$(COLOR_RESET)"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
