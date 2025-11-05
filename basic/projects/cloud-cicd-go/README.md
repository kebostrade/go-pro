# Tutorial 22: Cloud Platforms & CI/CD with Go

A comprehensive guide to deploying Go applications on Google Cloud Platform (GCP) and Amazon Web Services (AWS), with production-grade CI/CD pipelines.

## 📚 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Google Cloud Platform](#google-cloud-platform)
- [Amazon Web Services](#amazon-web-services)
- [CI/CD Pipelines](#cicd-pipelines)
- [Infrastructure as Code](#infrastructure-as-code)
- [Deployment Guides](#deployment-guides)
- [Best Practices](#best-practices)

## 🎯 Overview

This tutorial covers:

- **Google Cloud Platform**: Cloud Run, Cloud Functions, Cloud Storage, Pub/Sub, Firestore, GKE
- **Amazon Web Services**: Lambda, S3, DynamoDB, SQS, SNS, ECS, API Gateway
- **CI/CD**: GitHub Actions, GitLab CI, CircleCI, Jenkins
- **Infrastructure as Code**: Terraform for GCP and AWS
- **Containerization**: Docker multi-stage builds
- **Orchestration**: Kubernetes deployments

## 📋 Prerequisites

### Required Tools

- **Go 1.22+**
- **Docker** and Docker Compose
- **Git**
- **Make**

### Cloud Accounts

- **GCP Account** with billing enabled
- **AWS Account** with billing enabled

### CLI Tools

```bash
# Google Cloud SDK
curl https://sdk.cloud.google.com | bash
gcloud init

# AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
aws configure

# Terraform
wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
unzip terraform_1.6.0_linux_amd64.zip
sudo mv terraform /usr/local/bin/
```

## 🚀 Quick Start

### 1. Setup

```bash
# Clone and navigate
cd basic/projects/cloud-cicd-go

# Install dependencies
make deps

# Setup environment
make setup
# Edit .env with your credentials
```

### 2. Run Examples Locally

```bash
# GCP Examples
make run-cloud-storage
make run-pubsub
make run-firestore

# AWS Examples
make run-s3
make run-dynamodb
make run-sqs
make run-sns

# Cloud Run locally
make run-cloud-run
```

### 3. Deploy to Cloud

```bash
# Deploy to GCP Cloud Run
make deploy-cloud-run

# Deploy to AWS Lambda
make build-lambda
make deploy-lambda
```

## 📁 Project Structure

```
cloud-cicd-go/
├── gcp/                    # Google Cloud Platform examples
│   ├── cloud-run/         # Serverless containers
│   ├── cloud-functions/   # Serverless functions
│   ├── cloud-storage/     # Object storage
│   ├── pubsub/            # Messaging
│   ├── firestore/         # NoSQL database
│   └── gke/               # Kubernetes Engine
├── aws/                    # Amazon Web Services examples
│   ├── lambda/            # Serverless functions
│   ├── s3/                # Object storage
│   ├── dynamodb/          # NoSQL database
│   ├── sqs/               # Message queuing
│   ├── sns/               # Notifications
│   ├── ecs/               # Container service
│   └── api-gateway/       # API management
├── cicd/                   # CI/CD configurations
│   ├── github-actions/    # GitHub Actions workflows
│   ├── gitlab-ci/         # GitLab CI pipelines
│   ├── circleci/          # CircleCI config
│   └── jenkins/           # Jenkins pipelines
├── terraform/              # Infrastructure as Code
│   ├── gcp/               # GCP Terraform
│   └── aws/               # AWS Terraform
├── docker/                 # Docker configurations
├── kubernetes/             # Kubernetes manifests
└── scripts/                # Deployment scripts
```

## ☁️ Google Cloud Platform

### Cloud Run

Deploy containerized applications without managing servers.

**Features:**
- Automatic scaling (0 to N)
- Pay per use
- HTTPS endpoints
- Custom domains
- Traffic splitting

**Example:**

```bash
# Run locally
cd gcp/cloud-run
go run main.go

# Deploy
gcloud run deploy cloud-cicd-app \
  --source . \
  --region us-central1 \
  --allow-unauthenticated
```

### Cloud Storage

Object storage for any amount of data.

**Features:**
- Upload/download files
- Signed URLs
- Lifecycle management
- Versioning
- Public/private access

**Example:**

```bash
cd gcp/cloud-storage
export GCP_BUCKET_NAME=your-bucket
export GCP_PROJECT_ID=your-project
go run main.go
```

### Pub/Sub

Asynchronous messaging service.

**Features:**
- Publish/subscribe pattern
- At-least-once delivery
- Message ordering
- Dead letter queues
- Push/pull subscriptions

### Firestore

NoSQL document database.

**Features:**
- Real-time updates
- Offline support
- Automatic scaling
- ACID transactions
- Rich queries

## 🌐 Amazon Web Services

### Lambda

Serverless compute service.

**Features:**
- Event-driven execution
- Automatic scaling
- Pay per invocation
- Multiple triggers
- Function URLs

**Example:**

```bash
cd aws/lambda

# Build
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip function.zip bootstrap

# Deploy
aws lambda create-function \
  --function-name cloud-cicd-app \
  --runtime provided.al2 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::ACCOUNT_ID:role/lambda-role
```

### S3

Object storage service.

**Features:**
- Unlimited storage
- Versioning
- Lifecycle policies
- Presigned URLs
- Static website hosting

**Example:**

```bash
cd aws/s3
export AWS_BUCKET_NAME=your-bucket
go run main.go
```

### DynamoDB

NoSQL key-value database.

**Features:**
- Single-digit millisecond latency
- Automatic scaling
- Global tables
- Streams
- Transactions

### SQS

Message queuing service.

**Features:**
- Standard and FIFO queues
- Dead letter queues
- Message retention
- Visibility timeout
- Long polling

## 🔄 CI/CD Pipelines

### GitHub Actions

**Workflow:** `.github/workflows/deploy.yml`

**Features:**
- Automated testing
- Multi-platform builds
- Docker image building
- GCP and AWS deployment
- Slack notifications

**Triggers:**
- Push to main/develop
- Pull requests
- Manual dispatch

**Jobs:**
1. Test - Run unit tests
2. Lint - Code quality checks
3. Build - Compile binaries
4. Docker - Build and push images
5. Deploy GCP - Deploy to Cloud Run
6. Deploy AWS - Deploy to Lambda
7. Integration Test - Verify deployment
8. Notify - Send notifications

### GitLab CI

**Configuration:** `.gitlab-ci.yml`

**Stages:**
- Build
- Test
- Deploy to staging
- Deploy to production

### CircleCI

**Configuration:** `.circleci/config.yml`

**Workflows:**
- Build and test
- Deploy to cloud platforms

### Jenkins

**Pipeline:** `Jenkinsfile`

**Stages:**
- Checkout
- Build
- Test
- Deploy

## 🏗️ Infrastructure as Code

### Terraform - GCP

**Resources:**
- Cloud Run service
- Cloud Storage bucket
- Pub/Sub topic/subscription
- Firestore database
- GKE cluster
- Service accounts
- IAM bindings

**Usage:**

```bash
cd terraform/gcp

# Initialize
terraform init

# Plan
terraform plan -var="project_id=your-project"

# Apply
terraform apply -var="project_id=your-project"
```

### Terraform - AWS

**Resources:**
- Lambda function
- S3 bucket
- DynamoDB table
- SQS queue
- SNS topic
- API Gateway
- IAM roles/policies

**Usage:**

```bash
cd terraform/aws

# Initialize
terraform init

# Plan
terraform plan

# Apply
terraform apply
```

## 📖 Deployment Guides

### Deploy to GCP Cloud Run

```bash
# 1. Build Docker image
docker build -t gcr.io/PROJECT_ID/cloud-cicd-app:latest -f gcp/cloud-run/Dockerfile .

# 2. Push to Container Registry
docker push gcr.io/PROJECT_ID/cloud-cicd-app:latest

# 3. Deploy
gcloud run deploy cloud-cicd-app \
  --image gcr.io/PROJECT_ID/cloud-cicd-app:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated
```

### Deploy to AWS Lambda

```bash
# 1. Build function
cd aws/lambda
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# 2. Create deployment package
zip function.zip bootstrap

# 3. Deploy
aws lambda update-function-code \
  --function-name cloud-cicd-app \
  --zip-file fileb://function.zip
```

### Deploy with Terraform

```bash
# GCP
cd terraform/gcp
terraform init
terraform apply -var="project_id=your-project"

# AWS
cd terraform/aws
terraform init
terraform apply
```

## ✅ Best Practices

### Security
- Use service accounts with minimal permissions
- Store secrets in Secret Manager/Parameter Store
- Enable VPC for private resources
- Use HTTPS/TLS everywhere
- Implement authentication and authorization

### Performance
- Use connection pooling
- Implement caching strategies
- Optimize cold starts
- Use CDN for static assets
- Monitor and profile applications

### Cost Optimization
- Use auto-scaling
- Implement lifecycle policies
- Use spot/preemptible instances
- Monitor usage and set budgets
- Clean up unused resources

### Reliability
- Implement health checks
- Use multiple availability zones
- Set up monitoring and alerting
- Implement retry logic
- Use circuit breakers

### Observability
- Structured logging
- Distributed tracing
- Metrics collection
- Error tracking
- Performance monitoring

## 🧪 Testing

```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# Load testing
# Use tools like Apache Bench, wrk, or k6
```

## 📚 Additional Resources

- [GCP Documentation](https://cloud.google.com/docs)
- [AWS Documentation](https://docs.aws.amazon.com/)
- [Terraform Documentation](https://www.terraform.io/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

## 🎓 Learning Outcomes

After completing this tutorial, you will be able to:

✅ Deploy Go applications to GCP Cloud Run and AWS Lambda  
✅ Use cloud storage services (Cloud Storage, S3)  
✅ Implement messaging with Pub/Sub and SQS  
✅ Work with NoSQL databases (Firestore, DynamoDB)  
✅ Create CI/CD pipelines with GitHub Actions  
✅ Manage infrastructure with Terraform  
✅ Containerize applications with Docker  
✅ Deploy to Kubernetes (GKE)  
✅ Implement cloud best practices  
✅ Monitor and troubleshoot cloud applications  

## 📝 License

MIT License - see LICENSE file for details

