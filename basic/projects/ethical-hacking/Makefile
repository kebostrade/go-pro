.PHONY: help build clean test run-* deps

help: ## Show this help message
	@echo "Ethical Hacking with Go - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	go mod download
	go mod tidy

build: ## Build all tools
	@echo "Building all tools..."
	@mkdir -p bin
	go build -o bin/port-scanner ./network/port-scanner
	go build -o bin/network-mapper ./network/network-mapper
	go build -o bin/packet-sniffer ./network/packet-sniffer
	go build -o bin/web-scanner ./web/web-scanner
	go build -o bin/ssl-checker ./web/ssl-checker
	go build -o bin/hash-cracker ./crypto/hash-cracker
	go build -o bin/password-analyzer ./crypto/password-analyzer
	@echo "✅ All tools built successfully!"

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf results/
	rm -rf reports/

test: ## Run tests
	go test -v -race ./...

# Network Tools
run-port-scanner: ## Run port scanner
	go run ./network/port-scanner/main.go

run-network-mapper: ## Run network mapper
	go run ./network/network-mapper/main.go

run-packet-sniffer: ## Run packet sniffer (requires root/admin)
	sudo go run ./network/packet-sniffer/main.go

# Web Security Tools
run-web-scanner: ## Run web vulnerability scanner
	go run ./web/web-scanner/main.go

run-ssl-checker: ## Run SSL/TLS checker
	go run ./web/ssl-checker/main.go

# Cryptography Tools
run-hash-cracker: ## Run hash cracker
	go run ./crypto/hash-cracker/main.go

run-password-analyzer: ## Run password analyzer
	go run ./crypto/password-analyzer/main.go

# Examples
example-scan: ## Example: Scan localhost ports
	go run ./network/port-scanner/main.go -host localhost -ports 1-1000

example-web: ## Example: Scan example.com
	go run ./web/web-scanner/main.go -url https://example.com

example-ssl: ## Example: Check SSL certificate
	go run ./web/ssl-checker/main.go -host example.com

# Security
lint: ## Run security linter
	@which gosec > /dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format code
	go fmt ./...
	@which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

quality: fmt vet lint test ## Run all quality checks

