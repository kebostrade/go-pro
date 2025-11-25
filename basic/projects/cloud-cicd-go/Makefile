.PHONY: help deps test build clean docker-build docker-push deploy-gcp deploy-aws

# Variables
APP_NAME := cloud-cicd-app
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GCP_PROJECT := $(shell grep GCP_PROJECT_ID .env 2>/dev/null | cut -d '=' -f2)
GCP_REGION := us-central1
AWS_REGION := us-east-1
DOCKER_REGISTRY := gcr.io
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(GCP_PROJECT)/$(APP_NAME):$(VERSION)

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download
	go mod tidy

test: ## Run tests
	go test -v -race -cover ./...

build: ## Build the application
	go build -o bin/$(APP_NAME) ./gcp/cloud-run/main.go

clean: ## Clean build artifacts
	rm -rf bin/ dist/ build/

# GCP Commands
gcp-auth: ## Authenticate with GCP
	gcloud auth login
	gcloud config set project $(GCP_PROJECT)

run-cloud-run: ## Run Cloud Run example locally
	cd gcp/cloud-run && go run main.go

run-cloud-functions: ## Run Cloud Functions example
	cd gcp/cloud-functions && go run main.go

run-cloud-storage: ## Run Cloud Storage example
	cd gcp/cloud-storage && go run main.go

run-pubsub: ## Run Pub/Sub example
	cd gcp/pubsub && go run main.go

run-firestore: ## Run Firestore example
	cd gcp/firestore && go run main.go

deploy-cloud-run: ## Deploy to Cloud Run
	gcloud run deploy $(APP_NAME) \
		--source ./gcp/cloud-run \
		--region $(GCP_REGION) \
		--platform managed \
		--allow-unauthenticated

deploy-cloud-function: ## Deploy Cloud Function
	gcloud functions deploy $(APP_NAME) \
		--runtime go122 \
		--trigger-http \
		--allow-unauthenticated \
		--entry-point HandleRequest \
		--source ./gcp/cloud-functions

# AWS Commands
aws-auth: ## Configure AWS credentials
	aws configure

run-lambda: ## Run Lambda example locally
	cd aws/lambda && go run main.go

run-s3: ## Run S3 example
	cd aws/s3 && go run main.go

run-dynamodb: ## Run DynamoDB example
	cd aws/dynamodb && go run main.go

run-sqs: ## Run SQS example
	cd aws/sqs && go run main.go

run-sns: ## Run SNS example
	cd aws/sns && go run main.go

build-lambda: ## Build Lambda function
	cd aws/lambda && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	cd aws/lambda && zip function.zip bootstrap

deploy-lambda: build-lambda ## Deploy Lambda function
	aws lambda update-function-code \
		--function-name $(APP_NAME) \
		--zip-file fileb://aws/lambda/function.zip \
		--region $(AWS_REGION)

# Docker Commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) -f docker/Dockerfile .

docker-push: docker-build ## Push Docker image to registry
	docker push $(DOCKER_IMAGE)

docker-run: ## Run Docker container locally
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Terraform Commands
terraform-init-gcp: ## Initialize Terraform for GCP
	cd terraform/gcp && terraform init

terraform-plan-gcp: ## Plan Terraform changes for GCP
	cd terraform/gcp && terraform plan

terraform-apply-gcp: ## Apply Terraform changes for GCP
	cd terraform/gcp && terraform apply

terraform-init-aws: ## Initialize Terraform for AWS
	cd terraform/aws && terraform init

terraform-plan-aws: ## Plan Terraform changes for AWS
	cd terraform/aws && terraform plan

terraform-apply-aws: ## Apply Terraform changes for AWS
	cd terraform/aws && terraform apply

# Kubernetes Commands
k8s-apply: ## Apply Kubernetes manifests
	kubectl apply -f kubernetes/

k8s-delete: ## Delete Kubernetes resources
	kubectl delete -f kubernetes/

k8s-logs: ## View application logs
	kubectl logs -f deployment/$(APP_NAME)

k8s-port-forward: ## Port forward to service
	kubectl port-forward service/$(APP_NAME) 8080:80

# CI/CD
ci-test: ## Run CI tests
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

ci-lint: ## Run linters
	golangci-lint run ./...

ci-build: ## Build for CI
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(APP_NAME) ./gcp/cloud-run/main.go

# Development
dev: ## Run in development mode
	air

fmt: ## Format code
	go fmt ./...
	goimports -w .

lint: ## Run linter
	golangci-lint run

# All examples
run-all-gcp: ## Run all GCP examples
	@echo "Running all GCP examples..."
	@make run-cloud-storage
	@make run-pubsub
	@make run-firestore

run-all-aws: ## Run all AWS examples
	@echo "Running all AWS examples..."
	@make run-s3
	@make run-dynamodb
	@make run-sqs
	@make run-sns

# Setup
setup: ## Initial setup
	@echo "Setting up project..."
	cp .env.example .env
	@echo "Please edit .env with your credentials"
	make deps
	@echo "Setup complete!"

