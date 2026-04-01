# GoPro Learning Platform - Serverless Lambda Template

A production-ready AWS Lambda template using SAM (Serverless Application Model) with Go 1.23+.

## Architecture

This template provides a serverless function accessible via Lambda Function URLs (no API Gateway required).

```
┌─────────────────────────────────────────────────────────────┐
│                        AWS Lambda                            │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              GoProFunction (Lambda)                     │  │
│  │  - Function URL: https://xxx.lambda-url.region.on.aws  │  │
│  │  - Runtime: go1.x                                     │  │
│  │  - Memory: 256MB                                      │  │
│  │  - Timeout: 30s                                       │  │
│  │  - Auth: NONE (public access)                         │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Features

- **Lambda Function URLs**: Direct HTTPS access without API Gateway
- **Go 1.23+**: Latest Go runtime with modules
- **AWS SAM**: Infrastructure as Code for serverless applications
- **Docker Support**: Local testing with SAM CLI
- **CI/CD**: GitHub Actions workflow included

## Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check endpoint |
| `/event` | POST | Event processing endpoint |

## Prerequisites

- Go 1.23 or later
- AWS account with appropriate permissions
- AWS SAM CLI
- Docker (for local testing)

### Install SAM CLI

```bash
# macOS
brew install aws-sam-cli

# Linux
pip install aws-sam-cli

# Verify installation
sam --version
```

## Project Structure

```
serverless/
├── cmd/
│   └── handler/          # Lambda entry point
│       └── main.go
├── internal/
│   ├── handlers/         # Request handlers
│   │   ├── handler.go
│   │   └── handler_test.go
│   └── models/           # Data models
│       └── events.go
├── template.yaml         # SAM template
├── Dockerfile            # Container build
├── Makefile              # Build automation
├── go.mod                # Go module definition
└── README.md
```

## Quick Start

### Local Development

```bash
# Clone the repository
cd basic/projects/serverless

# Initialize Go modules
go mod tidy

# Run tests
make test

# Build for local testing
make build

# Local invoke (requires Docker)
make local
```

### Deploy to AWS

```bash
# Configure AWS credentials
aws configure

# Deploy using SAM
make deploysam

# Or step-by-step:
make build-linux
sam deploy --config-file .aws/config
```

### Testing the Lambda URL

After deployment, get the Lambda URL from the outputs:

```bash
# Health check
curl https://<function-url>.lambda-url.us-east-1.on.aws/health

# Send event
curl -X POST https://<function-url>.lambda-url.us-east-1.on.aws/event \
  -H "Content-Type: application/json" \
  -d '{
    "type": "lesson.completed",
    "timestamp": "2024-01-01T00:00:00Z",
    "payload": {
      "user_id": "user123",
      "lesson_id": "lesson456",
      "score": 95,
      "completed": true
    }
  }'
```

## Event Types

### lesson.completed

Triggered when a learner completes a lesson.

```json
{
  "type": "lesson.completed",
  "timestamp": "2024-01-01T00:00:00Z",
  "payload": {
    "user_id": "user123",
    "lesson_id": "lesson456",
    "score": 95,
    "completed": true
  }
}
```

### course.enrolled

Triggered when a learner enrolls in a course.

```json
{
  "type": "course.enrolled",
  "timestamp": "2024-01-01T00:00:00Z",
  "payload": {
    "user_id": "user123",
    "course_id": "course789"
  }
}
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | development | Environment (development/production) |
| `APP_VERSION` | 1.0.0 | Application version |
| `LOG_LEVEL` | info | Logging level |

### SAM Template Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| Environment | development | Deployment environment |

## Docker

### Build Container Image

```bash
docker build -t gopro-lambda:latest .
```

### Run Locally

```bash
docker run -p 9000:8080 \
  -e APP_ENV=development \
  gopro-lambda:latest
```

### Test Local Endpoint

```bash
curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations \
  -d '{"path": "/health", "httpMethod": "GET"}'
```

## CI/CD

The project includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that:

1. Runs tests on every push/PR
2. Builds the Lambda package
3. Runs SAM build validation
4. Lints the code with golangci-lint

### Required Secrets

- `AWS_ACCESS_KEY_ID`: AWS access key
- `AWS_SECRET_ACCESS_KEY`: AWS secret key

## Cost Optimization

- Memory: 256MB (sufficient for most workloads)
- Timeout: 30s
- No API Gateway fees (using Lambda Function URLs)
- Pay only for actual execution time

## Security

- Lambda Function URL uses NONE auth (public) - add your own authentication if needed
- CloudWatch Logs for observability
- IAM roles with least privilege

## Troubleshooting

### SAM build fails

```bash
# Clean and rebuild
make clean
sam build --use-container
```

### Lambda URL not accessible

Check the function's resource-based policy allows public access:

```bash
aws lambda get-policy --function-name gopro-handler
```

### Cold start issues

Increase memory or use Provisioned Concurrency for production:

```yaml
ProvisionedConcurrencyConfig:
  ProvisionedConcurrency: 5
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please read CONTRIBUTING.md for details.
