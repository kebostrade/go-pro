.PHONY: help build run docker-build docker-run k8s-deploy terraform-init terraform-apply

help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║                                                              ║"
	@echo "║        🐳 DevOps with Go: Docker, K8s, Terraform            ║"
	@echo "║                                                              ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available targets:"
	@echo "  make build              - Build Go application"
	@echo "  make run                - Run Go application locally"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run with Docker Compose"
	@echo "  make docker-stop        - Stop Docker Compose"
	@echo "  make k8s-deploy         - Deploy to Kubernetes"
	@echo "  make k8s-delete         - Delete from Kubernetes"
	@echo "  make terraform-init     - Initialize Terraform"
	@echo "  make terraform-plan     - Plan Terraform changes"
	@echo "  make terraform-apply    - Apply Terraform changes"
	@echo "  make terraform-destroy  - Destroy Terraform resources"

# Build Go application
build:
	@echo "🔨 Building application..."
	@go build -o bin/app ./app
	@echo "✅ Build complete"

# Run locally
run:
	@echo "🚀 Starting application..."
	@go run ./app

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t devops-go-app:latest -f docker/Dockerfile .
	@echo "✅ Docker image built"

docker-run:
	@echo "🐳 Starting Docker Compose..."
	@docker-compose -f docker/docker-compose.yml up -d
	@echo "✅ Services started"
	@echo "📍 App: http://localhost:8080"
	@echo "📊 Prometheus: http://localhost:9090"
	@echo "📈 Grafana: http://localhost:3000"

docker-stop:
	@echo "🛑 Stopping Docker Compose..."
	@docker-compose -f docker/docker-compose.yml down
	@echo "✅ Services stopped"

# Kubernetes commands
k8s-deploy:
	@echo "☸️  Deploying to Kubernetes..."
	@kubectl apply -f kubernetes/namespace.yaml
	@kubectl apply -f kubernetes/configmap.yaml
	@kubectl apply -f kubernetes/secret.yaml
	@kubectl apply -f kubernetes/deployment.yaml
	@kubectl apply -f kubernetes/service.yaml
	@kubectl apply -f kubernetes/ingress.yaml
	@echo "✅ Deployed to Kubernetes"
	@echo "📍 Check status: kubectl get pods -n devops-demo"

k8s-delete:
	@echo "🗑️  Deleting from Kubernetes..."
	@kubectl delete -f kubernetes/
	@echo "✅ Deleted from Kubernetes"

# Terraform commands
terraform-init:
	@echo "🏗️  Initializing Terraform..."
	@cd terraform && terraform init
	@echo "✅ Terraform initialized"

terraform-plan:
	@echo "📋 Planning Terraform changes..."
	@cd terraform && terraform plan
	@echo "✅ Plan complete"

terraform-apply:
	@echo "🚀 Applying Terraform changes..."
	@cd terraform && terraform apply -auto-approve
	@echo "✅ Infrastructure deployed"

terraform-destroy:
	@echo "💥 Destroying Terraform resources..."
	@cd terraform && terraform destroy -auto-approve
	@echo "✅ Infrastructure destroyed"

