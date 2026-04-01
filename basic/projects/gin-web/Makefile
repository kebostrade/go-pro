.PHONY: run build test lint docker docker-up docker-down clean fmt tidy verify

# Run the server
run:
	go run ./cmd/server

# Build the binary
build:
	go build -v ./...

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run linter
lint:
	golangci-lint run ./...

# Build Docker image
docker:
	docker build -t gin-web .

# Run Docker compose
docker-up:
	docker-compose up --build

# Stop Docker compose
docker-down:
	docker-compose down

# Clean up
clean:
	rm -f coverage.out
	rm -f gin-web

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Verify dependencies
verify:
	go mod verify

# Run tests and show coverage
coverage: test
	go tool cover -func=coverage.out | grep total
