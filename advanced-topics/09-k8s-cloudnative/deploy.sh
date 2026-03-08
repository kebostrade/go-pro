#!/bin/bash

# Kubernetes Deployment Script for Go Application
# This script builds, deploys, and tests a Go application on Kubernetes

set -e

APP_NAME="k8s-go-sample"
IMAGE_NAME="k8s-go-sample"
IMAGE_TAG="v1"
NAMESPACE="default"

echo "=== Kubernetes Deployment Script ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_step() {
    echo -e "${GREEN}→${NC} $1"
}

print_step "Building Go application"
cd sample-app
go build -o app
cd ..
echo ""

# Check if using Minikube
if command -v minikube &> /dev/null; then
    print_step "Minikube detected, building image with Minikube Docker"
    eval $(minikube docker-env)
else
    print_step "Building Docker image"
fi

docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
echo ""

print_step "Applying Kubernetes manifests"
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
kubectl apply -f hpa.yaml
echo ""

print_step "Waiting for deployment to be ready"
kubectl rollout status deployment/${APP_NAME}
echo ""

print_step "Getting pod status"
kubectl get pods -l app=${APP_NAME}
echo ""

print_step "Getting service status"
kubectl get svc ${APP_NAME}
echo ""

print_step "Getting ingress status"
kubectl get ingress ${APP_NAME}-ingress
echo ""

print_step "Deployment successful!"
echo ""
echo "Access the application:"
echo "  1. Port-forward: kubectl port-forward svc/${APP_NAME} 8080:80"
echo "  2. Minikube tunnel (for Ingress): minikube tunnel"
echo "  3. Test: curl http://localhost:8080/"
echo ""
echo "Monitor logs:"
echo "  kubectl logs -f deployment/${APP_NAME}"
echo ""
echo "Cleanup:"
echo "  kubectl delete all -l app=${APP_NAME}"
