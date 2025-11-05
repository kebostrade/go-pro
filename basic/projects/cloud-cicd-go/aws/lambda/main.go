package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Request represents the Lambda function input
type Request struct {
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// Response represents the Lambda function output
type Response struct {
	StatusCode int                    `json:"statusCode"`
	Body       string                 `json:"body"`
	Headers    map[string]string      `json:"headers"`
	Success    bool                   `json:"success"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// Handler handles Lambda invocations
type Handler struct {
	// Add any dependencies here (DB clients, etc.)
}

// NewHandler creates a new handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// HandleRequest is the main Lambda handler function
func (h *Handler) HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s %s", event.HTTPMethod, event.Path)
	log.Printf("Headers: %v", event.Headers)
	log.Printf("Query params: %v", event.QueryStringParameters)
	
	// Route based on HTTP method and path
	switch event.HTTPMethod {
	case "GET":
		return h.handleGet(ctx, event)
	case "POST":
		return h.handlePost(ctx, event)
	case "PUT":
		return h.handlePut(ctx, event)
	case "DELETE":
		return h.handleDelete(ctx, event)
	default:
		return h.errorResponse(405, "Method not allowed"), nil
	}
}

// handleGet handles GET requests
func (h *Handler) handleGet(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Check path
	if event.Path == "/health" {
		return h.successResponse(200, "Service is healthy", map[string]interface{}{
			"status":    "UP",
			"timestamp": time.Now(),
			"service":   "lambda-demo",
		}), nil
	}
	
	if event.Path == "/users" {
		// Return sample users
		users := []map[string]interface{}{
			{"id": "1", "name": "Alice", "email": "alice@example.com"},
			{"id": "2", "name": "Bob", "email": "bob@example.com"},
		}
		return h.successResponse(200, "Users retrieved", map[string]interface{}{
			"users": users,
			"count": len(users),
		}), nil
	}
	
	// Get user by ID from path parameter
	if id, ok := event.PathParameters["id"]; ok {
		user := map[string]interface{}{
			"id":    id,
			"name":  "Sample User",
			"email": "user@example.com",
		}
		return h.successResponse(200, "User retrieved", user), nil
	}
	
	return h.errorResponse(404, "Not found"), nil
}

// handlePost handles POST requests
func (h *Handler) handlePost(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return h.errorResponse(400, "Invalid request body"), nil
	}
	
	// Validate request
	if req.Name == "" || req.Email == "" {
		return h.errorResponse(400, "Name and email are required"), nil
	}
	
	// Process request
	result := map[string]interface{}{
		"id":         "123",
		"name":       req.Name,
		"email":      req.Email,
		"created_at": time.Now(),
		"action":     req.Action,
	}
	
	return h.successResponse(201, "Resource created successfully", result), nil
}

// handlePut handles PUT requests
func (h *Handler) handlePut(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return h.errorResponse(400, "ID parameter is required"), nil
	}
	
	var req Request
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return h.errorResponse(400, "Invalid request body"), nil
	}
	
	result := map[string]interface{}{
		"id":         id,
		"name":       req.Name,
		"email":      req.Email,
		"updated_at": time.Now(),
	}
	
	return h.successResponse(200, "Resource updated successfully", result), nil
}

// handleDelete handles DELETE requests
func (h *Handler) handleDelete(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return h.errorResponse(400, "ID parameter is required"), nil
	}
	
	result := map[string]interface{}{
		"id":         id,
		"deleted_at": time.Now(),
	}
	
	return h.successResponse(200, "Resource deleted successfully", result), nil
}

// successResponse creates a success response
func (h *Handler) successResponse(statusCode int, message string, data map[string]interface{}) events.APIGatewayProxyResponse {
	response := Response{
		StatusCode: statusCode,
		Success:    true,
		Message:    message,
		Data:       data,
	}
	
	body, _ := json.Marshal(response)
	
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}
}

// errorResponse creates an error response
func (h *Handler) errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	response := Response{
		StatusCode: statusCode,
		Success:    false,
		Message:    message,
	}
	
	body, _ := json.Marshal(response)
	
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}
}

// HandleS3Event handles S3 events
func HandleS3Event(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		log.Printf("Processing S3 event:")
		log.Printf("  Bucket: %s", record.S3.Bucket.Name)
		log.Printf("  Key: %s", record.S3.Object.Key)
		log.Printf("  Size: %d bytes", record.S3.Object.Size)
		log.Printf("  Event: %s", record.EventName)
		
		// Process the file here
		// For example: download, transform, upload to another bucket
	}
	
	return nil
}

// HandleSQSEvent handles SQS events
func HandleSQSEvent(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		log.Printf("Processing SQS message:")
		log.Printf("  MessageId: %s", record.MessageId)
		log.Printf("  Body: %s", record.Body)
		
		// Process the message here
		var message map[string]interface{}
		if err := json.Unmarshal([]byte(record.Body), &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}
		
		log.Printf("  Parsed message: %v", message)
	}
	
	return nil
}

// HandleScheduledEvent handles CloudWatch Events (scheduled)
func HandleScheduledEvent(ctx context.Context, event events.CloudWatchEvent) error {
	log.Printf("Scheduled event triggered:")
	log.Printf("  Source: %s", event.Source)
	log.Printf("  Detail Type: %s", event.DetailType)
	log.Printf("  Time: %s", event.Time)
	
	// Perform scheduled task here
	// For example: cleanup, data aggregation, report generation
	
	return nil
}

func main() {
	handler := NewHandler()
	
	// Start Lambda handler
	// For API Gateway events
	lambda.Start(handler.HandleRequest)
	
	// For S3 events, use:
	// lambda.Start(HandleS3Event)
	
	// For SQS events, use:
	// lambda.Start(HandleSQSEvent)
	
	// For scheduled events, use:
	// lambda.Start(HandleScheduledEvent)
}

// Example usage:
//
// Deploy:
//   GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
//   zip function.zip bootstrap
//   aws lambda create-function \
//     --function-name my-function \
//     --runtime provided.al2 \
//     --handler bootstrap \
//     --zip-file fileb://function.zip \
//     --role arn:aws:iam::ACCOUNT_ID:role/lambda-role
//
// Invoke:
//   aws lambda invoke \
//     --function-name my-function \
//     --payload '{"name":"John","email":"john@example.com"}' \
//     response.json
//
// Update:
//   aws lambda update-function-code \
//     --function-name my-function \
//     --zip-file fileb://function.zip

