.PHONY: build run test lint docker clean

# Build the binary
build:
	go build -v -o cli-app ./cmd/cli

# Run the greet command
run:
	go run ./cmd/cli greet --name "Go Developer"

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run linter
lint:
	golangci-lint run ./...

# Build Docker image
docker:
	docker build -t cli-app .

# Clean up
clean:
	rm -f coverage.out
	rm -f cli-app

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Tidy dependencies
tidy:
	go mod tidy
