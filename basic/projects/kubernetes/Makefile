.PHONY: help test build docker-build docker-run helm-deploy helm-template helm-lint manifest-validate clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
BINARY_NAME=operator
BUILD_DIR=.

help: ## Display this help message
	@echo "Kubernetes Operator Makefile"
	@echo ""
	@echo "Usage:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-20s %s\n", $$1, $$2}'

test: ## Run tests
	$(GOTEST) -short -race ./...

test-e2e: ## Run e2e tests (requires k8s cluster)
	$(GOTEST) -v ./...

build: ## Build operator binary
	$(GOBUILD) -o $(BINARY_NAME) ./cmd

run: ## Run operator locally (requires k8s config)
	$(GOCMD) run ./cmd

vet: ## Run go vet
	$(GOVET) ./...

mod-tidy: ## Tidy go modules
	$(GOMOD) tidy

docker-build: ## Build Docker image
	docker build -t gopro-k8s:latest .

docker-buildx: ## Build Docker image with buildx
	docker buildx build -t gopro-k8s:latest --platform linux/amd64,linux/arm64 .

docker-run: ## Run Docker container
	docker run -p 8080:8080 -p 8081:8081 gopro-k8s:latest

helm-lint: ## Lint Helm chart
	helm lint ./helm/gopro-chart

helm-template: ## Render Helm templates
	helm template my-release ./helm/gopro-chart

helm-deploy: ## Deploy using Helm
	helm install gopro-release ./helm/gopro-chart --namespace gopro --create-namespace

helm-upgrade: ## Upgrade Helm deployment
	helm upgrade gopro-release ./helm/gopro-chart --namespace gopro

helm-uninstall: ## Uninstall Helm release
	helm uninstall gopro-release --namespace gopro

manifest-validate: ## Validate Kubernetes manifests (requires kubectl)
	kubectl apply --dry-run=client -f manifests/

manifest-apply: ## Apply Kubernetes manifests
	kubectl apply -f manifests/ --namespace gopro

manifest-delete: ## Delete Kubernetes manifests
	kubectl delete -f manifests/ --namespace gopro

# Generate DeepCopy functions (requires controller-gen)
generate: ## Generate deepcopy code
	$(GOCMD) generate ./...

# Install kubebuilder binaries for local development
install-kubebuilder: ## Install kubebuilder locally
	@echo "Installing kubebuilder..."
	os=$(shell go env GOOS) && \
	arch=$(shell go env GOARCH) && \
	curl -L -o /tmp/kubebuilder.tar.gz "https://go.kubebuilder.io/dl/$${os}/$${arch}/latest/kubebuilder" && \
	tar -xzf /tmp/kubebuilder.tar.gz -C /tmp && \
	mv /tmp/kubebuilder_*/bin/* $(GOPATH)/bin/

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	docker rmi gopro-k8s:latest 2>/dev/null || true
