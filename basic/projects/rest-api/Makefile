.PHONY: run build test lint docker clean

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
	docker build -t rest-api .

# Run Docker compose
docker-up:
	docker-compose up --build

# Stop Docker compose
docker-down:
	docker-compose down

# Clean up
clean:
	rm -f coverage.out
	rm -f rest-api

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Tidy dependencies
tidy:
	go mod tidy

# Verify dependencies
verify:
	go mod verify
