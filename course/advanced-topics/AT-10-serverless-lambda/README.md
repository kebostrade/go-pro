# Building Serverless Applications with Go and AWS Lambda

Create serverless applications using Go and AWS Lambda.

## Learning Objectives

- Write Lambda functions in Go
- Handle different event sources
- Use AWS SDK with Lambda
- Implement cold start optimization
- Deploy with SAM/Serverless Framework
- Debug Lambda functions

## Theory

### Basic Lambda Handler

```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
    Name string `json:"name"`
}

type Response struct {
    Message string `json:"message"`
}

func handleRequest(ctx context.Context, event Event) (Response, error) {
    return Response{
        Message: fmt.Sprintf("Hello, %s!", event.Name),
    }, nil
}

func main() {
    lambda.Start(handleRequest)
}
```

### API Gateway Handler

```go
package main

import (
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    switch req.HTTPMethod {
    case "GET":
        return handleGet(ctx, req)
    case "POST":
        return handlePost(ctx, req)
    default:
        return events.APIGatewayProxyResponse{
            StatusCode: 405,
            Body:       `{"error":"method not allowed"}`,
        }, nil
    }
}

func handleGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    id := req.PathParameters["id"]
    
    user, err := getUser(ctx, id)
    if err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: 404,
            Body:       `{"error":"not found"}`,
        }, nil
    }

    body, _ := json.Marshal(user)
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Headers:    map[string]string{"Content-Type": "application/json"},
        Body:       string(body),
    }, nil
}

func handlePost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var user User
    if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: 400,
            Body:       `{"error":"invalid request"}`,
        }, nil
    }

    if err := createUser(ctx, &user); err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: 500,
            Body:       `{"error":"internal error"}`,
        }, nil
    }

    body, _ := json.Marshal(user)
    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Body:       string(body),
    }, nil
}

func main() {
    lambda.Start(handleRequest)
}
```

### DynamoDB Stream Handler

```go
func handleDynamoDBStream(ctx context.Context, event events.DynamoDBEvent) error {
    for _, record := range event.Records {
        switch record.EventName {
        case "INSERT":
            if err := handleInsert(ctx, record.Change); err != nil {
                log.Printf("insert error: %v", err)
                return err
            }
        case "MODIFY":
            if err := handleModify(ctx, record.Change); err != nil {
                log.Printf("modify error: %v", err)
                return err
            }
        case "REMOVE":
            if err := handleRemove(ctx, record.Change); err != nil {
                log.Printf("remove error: %v", err)
                return err
            }
        }
    }
    return nil
}

func handleInsert(ctx context.Context, item events.DynamoDBStreamRecord) error {
    id := item.Keys["id"].S()
    name := item.NewImage["name"].S()
    log.Printf("New item: id=%s, name=%s", id, name)
    return nil
}
```

### S3 Event Handler

```go
func handleS3Event(ctx context.Context, event events.S3Event) error {
    for _, record := range event.Records {
        bucket := record.S3.Bucket.Name
        key := record.S3.Object.Key

        log.Printf("Processing s3://%s/%s", bucket, key)

        if err := processS3Object(ctx, bucket, key); err != nil {
            return fmt.Errorf("process %s/%s: %w", bucket, key, err)
        }
    }
    return nil
}
```

### Cold Start Optimization

```go
var (
    db     *sql.DB
    s3Svc  *s3.Client
    once   sync.Once
    initErr error
)

func initClients() error {
    once.Do(func() {
        cfg, err := config.LoadDefaultConfig(context.Background())
        if err != nil {
            initErr = err
            return
        }
        s3Svc = s3.NewFromConfig(cfg)
    })
    return initErr
}

func handleRequest(ctx context.Context, event Event) (Response, error) {
    if err := initClients(); err != nil {
        return Response{}, err
    }
}
```

### SAM Template

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Globals:
  Function:
    Runtime: provided.al2
    Handler: bootstrap
    Timeout: 30
    MemorySize: 256
    Environment:
      Variables:
        TABLE_NAME: !Ref UsersTable

Resources:
  UsersTable:
    Type: AWS::DynamoDB::Table
    Properties:
      BillingMode: PAY_PER_REQUEST
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH

  GetUserFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: bin/
      Events:
        GetUser:
          Type: Api
          Properties:
            Path: /users/{id}
            Method: get

  CreateUserFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: bin/
      Events:
        CreateUser:
          Type: Api
          Properties:
            Path: /users
            Method: post

Outputs:
  ApiUrl:
    Value: !Sub "https://${ServerlessRestApi}.execute-api.${AWS::Region}.amazonaws.com/Prod"
```

## Security Considerations

```go
func validateInput(input string) error {
    if len(input) > 1000 {
        return errors.New("input too long")
    }
    if strings.ContainsAny(input, "<>") {
        return errors.New("invalid characters")
    }
    return nil
}

func getSecret(secretName string) (string, error) {
    cfg, _ := config.LoadDefaultConfig(context.Background())
    svc := secretsmanager.NewFromConfig(cfg)
    
    resp, err := svc.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
        SecretId: &secretName,
    })
    if err != nil {
        return "", err
    }
    return *resp.SecretString, nil
}
```

## Performance Tips

```go
var warmContainer bool

func handleRequest(ctx context.Context, event Event) (Response, error) {
    if !warmContainer {
        log.Println("Cold start!")
        warmContainer = true
    }

    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
}
```

## Exercises

1. Build a REST API with Lambda + DynamoDB
2. Process S3 file uploads
3. Handle SQS messages
4. Implement Lambda layers

## Validation

```bash
cd exercises
GOOS=linux GOARCH=amd64 go build -o bootstrap
sam local invoke -e events/test.json
sam local start-api
```

## Key Takeaways

- Initialize clients outside handler
- Keep functions small and focused
- Handle timeouts gracefully
- Use environment variables for config
- Monitor cold starts

## Next Steps

**[AT-11: ML Gorgonia](../AT-11-ml-gorgonia/README.md)**

---

Serverless: scale to zero, pay per use. ☁️
