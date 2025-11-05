#!/bin/bash

# Deployment script for cloud-cicd-go
# Supports GCP Cloud Run, AWS Lambda, and Kubernetes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Deploy to GCP Cloud Run
deploy_gcp() {
    log_info "Deploying to GCP Cloud Run..."
    
    if ! command_exists gcloud; then
        log_error "gcloud CLI not found. Please install Google Cloud SDK."
        exit 1
    fi
    
    # Check if project is set
    PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
    if [ -z "$PROJECT_ID" ]; then
        log_error "GCP project not set. Run: gcloud config set project PROJECT_ID"
        exit 1
    fi
    
    log_info "Using GCP project: $PROJECT_ID"
    
    # Deploy
    gcloud run deploy cloud-cicd-app \
        --source ./gcp/cloud-run \
        --region us-central1 \
        --platform managed \
        --allow-unauthenticated \
        --set-env-vars "ENV=production"
    
    # Get service URL
    SERVICE_URL=$(gcloud run services describe cloud-cicd-app \
        --region us-central1 \
        --format 'value(status.url)')
    
    log_info "Deployment successful!"
    log_info "Service URL: $SERVICE_URL"
}

# Deploy to AWS Lambda
deploy_aws() {
    log_info "Deploying to AWS Lambda..."
    
    if ! command_exists aws; then
        log_error "AWS CLI not found. Please install AWS CLI."
        exit 1
    fi
    
    # Build Lambda function
    log_info "Building Lambda function..."
    cd aws/lambda
    GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
    zip function.zip bootstrap
    
    # Deploy
    log_info "Deploying to Lambda..."
    aws lambda update-function-code \
        --function-name cloud-cicd-app \
        --zip-file fileb://function.zip \
        --region us-east-1
    
    # Update environment variables
    aws lambda update-function-configuration \
        --function-name cloud-cicd-app \
        --environment "Variables={ENV=production}" \
        --region us-east-1
    
    # Get function URL
    FUNCTION_URL=$(aws lambda get-function-url-config \
        --function-name cloud-cicd-app \
        --region us-east-1 \
        --query 'FunctionUrl' \
        --output text 2>/dev/null || echo "No function URL configured")
    
    log_info "Deployment successful!"
    log_info "Function URL: $FUNCTION_URL"
    
    # Cleanup
    rm bootstrap function.zip
    cd ../..
}

# Deploy to Kubernetes
deploy_k8s() {
    log_info "Deploying to Kubernetes..."
    
    if ! command_exists kubectl; then
        log_error "kubectl not found. Please install kubectl."
        exit 1
    fi
    
    # Check if cluster is accessible
    if ! kubectl cluster-info >/dev/null 2>&1; then
        log_error "Cannot connect to Kubernetes cluster."
        exit 1
    fi
    
    # Apply manifests
    log_info "Applying Kubernetes manifests..."
    kubectl apply -f kubernetes/deployment.yaml
    
    # Wait for deployment
    log_info "Waiting for deployment to be ready..."
    kubectl wait --for=condition=available --timeout=300s \
        deployment/cloud-cicd-app -n cloud-cicd
    
    # Get service URL
    SERVICE_IP=$(kubectl get service cloud-cicd-app -n cloud-cicd \
        -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "Pending")
    
    log_info "Deployment successful!"
    log_info "Service IP: $SERVICE_IP"
}

# Build Docker image
build_docker() {
    log_info "Building Docker image..."
    
    if ! command_exists docker; then
        log_error "Docker not found. Please install Docker."
        exit 1
    fi
    
    # Build image
    docker build -t cloud-cicd-app:latest -f docker/Dockerfile .
    
    log_info "Docker image built successfully!"
}

# Push Docker image
push_docker() {
    log_info "Pushing Docker image..."
    
    REGISTRY=${1:-gcr.io}
    PROJECT_ID=${2:-$(gcloud config get-value project 2>/dev/null)}
    
    if [ -z "$PROJECT_ID" ]; then
        log_error "Project ID not provided and cannot be determined."
        exit 1
    fi
    
    # Tag image
    docker tag cloud-cicd-app:latest $REGISTRY/$PROJECT_ID/cloud-cicd-app:latest
    
    # Push image
    docker push $REGISTRY/$PROJECT_ID/cloud-cicd-app:latest
    
    log_info "Docker image pushed successfully!"
}

# Run tests
run_tests() {
    log_info "Running tests..."
    
    go test -v -race -cover ./...
    
    log_info "Tests passed!"
}

# Main script
main() {
    case "$1" in
        gcp)
            deploy_gcp
            ;;
        aws)
            deploy_aws
            ;;
        k8s|kubernetes)
            deploy_k8s
            ;;
        docker)
            build_docker
            ;;
        push)
            push_docker "$2" "$3"
            ;;
        test)
            run_tests
            ;;
        all)
            run_tests
            build_docker
            deploy_gcp
            deploy_aws
            ;;
        *)
            echo "Usage: $0 {gcp|aws|k8s|docker|push|test|all}"
            echo ""
            echo "Commands:"
            echo "  gcp       - Deploy to GCP Cloud Run"
            echo "  aws       - Deploy to AWS Lambda"
            echo "  k8s       - Deploy to Kubernetes"
            echo "  docker    - Build Docker image"
            echo "  push      - Push Docker image (usage: push [registry] [project-id])"
            echo "  test      - Run tests"
            echo "  all       - Run tests and deploy to all platforms"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"

