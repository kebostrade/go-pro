# AWS Lambda Serverless with Go

Learn to build serverless applications using AWS Lambda and Go. This topic covers Lambda functions, event sources, and deployment with AWS SAM.

## Learning Objectives

- Understand AWS Lambda fundamentals
- Write Lambda handlers in Go
- Integrate with API Gateway, DynamoDB, and S3
- Use AWS SAM for deployment
- Test Lambda functions locally
- Follow serverless best practices

## Prerequisites

- Go 1.23+ installed
- AWS account (free tier works)
- AWS CLI configured
- Docker (for local testing)
- AWS SAM CLI

## Setup

```bash
# Install AWS SAM CLI
macOS:
brew tap aws/tap
brew install aws-sam-cli

Linux:
wget https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-linux-x86_64.zip
unzip aws-sam-cli-linux-x86_64.zip -d sam-installation
sudo ./sam-installation/install

Windows:
# Download from https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html

# Verify installation
sam --version

# Configure AWS CLI
aws configure
```

## Examples

### 1. Basic Lambda Handler
`examples/lambda_handler.go` - Simple "Hello World" Lambda function

### 2. API Gateway Integration
`examples/api_gateway.go` - REST API with Lambda proxy integration

### 3. DynamoDB Integration
`examples/dynamodb.go` - CRUD operations with DynamoDB

### 4. S3 Event Processing
`examples/s3_events.go` - Process S3 upload events

## Running Examples

### Local Testing with SAM

```bash
# Build the Lambda function
cd 14-serverless-lambda
sam build

# Invoke Lambda locally
sam local invoke LambdaHandler --event events/example.json

# Start local API Gateway
sam local start-api

# Test the endpoint
curl http://localhost:3000/hello
```

### Deploy to AWS

```bash
# Deploy all resources
sam deploy --guided

# First deployment will prompt for:
# - Stack name
# - AWS Region
# - Confirm changes before deploy
# - Allow SAM CLI IAM role creation
# - Save arguments to samconfig.toml

# Subsequent deployments
sam deploy

# Get the API URL
aws cloudformation describe-stacks \
  --stack-name serverless-lambda-go \
  --query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
  --output text
```

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│  API Gateway │────▶│   Lambda    │
└─────────────┘     └──────────────┘     └──────┬──────┘
                                               │
                    ┌─────────────┐            │
                    │    DynamoDB │◀───────────┘
                    └─────────────┘

┌─────────────┐     ┌──────────────┐
│     S3      │────▶│   Lambda     │────▶│ Processing │
│  (Upload)   │     │  (Events)    │     └─────────────┘
└─────────────┘     └──────────────┘
```

## Key Concepts

### Lambda Handler Patterns

1. **Simple Handler**: No payload processing
```go
func handler() error {
    fmt.Println("Hello from Lambda!")
    return nil
}
```

2. **API Gateway Request Handler**: HTTP request/response
```go
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       "Hello World",
    }, nil
}
```

3. **Event Source Handler**: S3, DynamoDB, etc.
```go
func handler(ctx context.Context, req events.S3Event) error {
    for _, record := range req.Records {
        processS3Record(record)
    }
    return nil
}
```

### Environment Variables

```go
// Access in Lambda
dbTable := os.Getenv("DYNAMODB_TABLE")
region := os.Getenv("AWS_REGION")
```

### Context Timeout

```go
func handler(ctx context.Context) error {
    // Lambda context automatically times out
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-process():
        return nil
    }
}
```

## Best Practices

### 1. Keep Handlers Stateless
- ❌ Don't store state in global variables
- ✅ Use external services (DynamoDB, S3, etc.)

### 2. Handle Context Properly
- Always check context.Done()
- Use context for all I/O operations
- Respect Lambda timeout

### 3. Optimize Cold Starts
- Minimize dependencies
- Use AWS Lambda provisioned concurrency
- Keep deployment package small

### 4. Error Handling
- Return errors for Lambda to retry
- Use structured logging
- Set appropriate retry policies

### 5. Security
- Use least-privilege IAM roles
- Never hardcode credentials
- Enable encryption at rest
- Use VPC when needed

## Cost Optimization

- **Free Tier**: 1M requests/month free
- **Memory vs CPU**: Adjust based on needs
- **Provisioned Concurrency**: For consistent performance
- **Reserved Concurrency**: Prevent runaway costs

## Monitoring

```bash
# View CloudWatch logs
aws logs tail /aws/lambda/serverless-lambda-go --follow

# Get Lambda metrics
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Invocations \
  --dimensions Name=FunctionName,Value=serverless-lambda-go \
  --start-time 2024-01-01T00:00:00Z \
  --end-time 2024-01-02T00:00:00Z \
  --period 86400 \
  --statistics Sum
```

## Troubleshooting

### Cold Start Issues
- Increase memory allocation
- Use provisioned concurrency
- Reduce package size

### Timeout Errors
- Increase timeout in AWS SAM template
- Optimize code performance
- Check external service latency

### IAM Permission Errors
- Verify IAM role policies
- Check resource ARNs
- Use CloudTrail for debugging

## Resources

- [AWS Lambda Go API](https://pkg.go.dev/github.com/aws/aws-lambda-go)
- [AWS SAM Developer Guide](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/)
- [AWS Lambda Developer Guide](https://docs.aws.amazon.com/lambda/latest/dg/)
- [Serverless Best Practices](https://serverless.com/blog/serverless-best-practices/)

## Next Steps

1. Complete all examples
2. Build a serverless REST API
3. Add authentication with Cognito
4. Implement CI/CD with AWS CodePipeline
5. Explore Step Functions for orchestration
