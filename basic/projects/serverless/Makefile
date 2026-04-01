.PHONY: build test vet lint deploysam clean

# Build the Lambda handler
build:
	go build -o handler ./cmd/handler

# Run tests
test:
	go test -v -cover ./...

# Run go vet
vet:
	go vet ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run --timeout 5m

# Build for Linux (for SAM deployment)
build-linux:
	GOOS=linux GOARCH=amd64 go build -o handler ./cmd/handler

# Deploy using SAM (requires AWS credentials)
deploysam: build-linux
	sam deploy --config-file .aws/config

# Local invocation testing (requires Docker)
local:
	sam local invoke GoProFunction -e events/event.json

# Clean build artifacts
clean:
	rm -f handler
	rm -rf .aws-sam
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Run all quality checks
quality: vet lint test

# Initialize the project (first time setup)
init:
	@echo "Ensure you have AWS SAM CLI installed: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html"
