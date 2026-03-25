package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Example 1: Simple Handler - No payload processing
func simpleHandler() error {
	fmt.Println("Hello from Lambda!")
	return nil
}

// Example 2: Handler with return value
func helloHandler() (string, error) {
	return "Hello from Lambda!", nil
}

// Example 3: Handler with context and return value
func contextHandler(ctx context.Context) (string, error) {
	// Access context metadata
	deadline, _ := ctx.Deadline()
	fmt.Printf("Lambda will timeout at: %v\n", deadline)

	// Check if context was cancelled (timeout)
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		// Continue processing
	}

	return "Context-aware handler", nil
}

// Example 4: API Gateway Proxy Request Handler
func apiGatewayHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract request information
	fmt.Printf("HTTP Method: %s\n", req.HTTPMethod)
	fmt.Printf("Path: %s\n", req.Path)
	fmt.Printf("Query String: %v\n", req.QueryStringParameters)
	fmt.Printf("Headers: %v\n", req.Headers)
	fmt.Printf("Body: %s\n", req.Body)

	// Parse query parameters
	name := "World"
	if req.QueryStringParameters != nil {
		if n, ok := req.QueryStringParameters["name"]; ok {
			name = n
		}
	}

	// Create response
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"message": "Hello, %s!"}`, name),
	}

	return response, nil
}

// Example 5: Handler with structured input/output
type Request struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func structuredHandler(ctx context.Context, req Request) (Response, error) {
	return Response{
		Message:   fmt.Sprintf("%s, %s!", req.Message, req.Name),
		Timestamp: time.Now(),
	}, nil
}

// Example 6: Error handling
func errorHandlingHandler(ctx context.Context, req Request) (Response, error) {
	if req.Name == "" {
		return Response{}, fmt.Errorf("name is required")
	}

	return Response{
		Message:   fmt.Sprintf("Hello, %s!", req.Name),
		Timestamp: time.Now(),
	}, nil
}

// Example 7: HTTP status code handling
func httpStatusHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := req.QueryStringParameters["name"]

	if name == "" {
		// Return 400 Bad Request
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "name parameter is required"}`,
		}, nil
	}

	// Return 200 OK
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"message": "Hello, %s!"}`, name),
	}, nil
}

// Example 8: CORS handling
func corsHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle preflight OPTIONS request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Methods":     "GET,OPTIONS,POST,PUT,DELETE",
				"Access-Control-Max-Age":           "86400",
			},
		}, nil
	}

	// Handle actual request
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                 "application/json",
		},
		Body: `{"message": "CORS enabled"}`,
	}, nil
}

func main() {
	// Choose which handler to use
	// lambda.Start(simpleHandler)
	// lambda.Start(helloHandler)
	// lambda.Start(contextHandler)
	// lambda.Start(apiGatewayHandler)
	// lambda.Start(structuredHandler)
	lambda.Start(corsHandler)
}

/*
Testing Lambda Handlers Locally:

1. Using AWS SAM:
   sam local invoke LambdaHandler --event events/example.json

2. Creating test events (events/example.json):
   {
     "httpMethod": "GET",
     "path": "/hello",
     "queryStringParameters": {
       "name": "World"
     },
     "headers": {
       "Content-Type": "application/json"
     },
     "body": "{\"name\":\"World\",\"message\":\"Hello\"}"
   }

3. Testing API Gateway locally:
   sam local start-api
   curl http://localhost:3000/hello?name=World

4. Build for deployment:
   GOOS=linux GOARCH=amd64 go build -o main

Handler Signatures:
- func ()
- func () error
- func (TIn) error
- func () (TOut, error)
- func (TIn) (TOut, error)
- func (context.Context) error
- func (context.Context) (TOut, error)
- func (context.Context, TIn) error
- func (context.Context, TIn) (TOut, error)
*/
