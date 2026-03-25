package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Response is a standard API response
type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Router handles different HTTP methods and paths
func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s", req.HTTPMethod, req.Path)

	switch req.Path {
	case "/users":
		return handleUsers(req)
	case "/users/{id}":
		return handleUserByID(req)
	case "/health":
		return handleHealth()
	default:
		return notFound()
	}
}

// handleUsers handles /users endpoint
func handleUsers(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return getUsers()
	case "POST":
		return createUser(req)
	default:
		return methodNotAllowed()
	}
}

// handleUserByID handles /users/{id} endpoint
func handleUserByID(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract ID from path parameters
	userID := req.PathParameters["id"]

	switch req.HTTPMethod {
	case "GET":
		return getUserByID(userID)
	case "PUT":
		return updateUser(userID, req)
	case "DELETE":
		return deleteUser(userID)
	default:
		return methodNotAllowed()
	}
}

// getUsers returns all users
func getUsers() (events.APIGatewayProxyResponse, error) {
	users := []User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}

	return jsonMarshal(200, Response{
		Data: users,
	})
}

// createUser creates a new user
func createUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return badRequest("Invalid request body")
	}

	// Validate input
	if user.Name == "" || user.Email == "" {
		return badRequest("Name and email are required")
	}

	// Generate ID (in real app, use database)
	user.ID = fmt.Sprintf("%d", len(user.Name)+len(user.Email))

	return jsonMarshal(201, Response{
		Message: "User created successfully",
		Data:    user,
	})
}

// getUserByID returns a specific user
func getUserByID(id string) (events.APIGatewayProxyResponse, error) {
	// Mock database lookup
	user := User{
		ID:    id,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	return jsonMarshal(200, Response{
		Data: user,
	})
}

// updateUser updates a user
func updateUser(id string, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var updates User
	if err := json.Unmarshal([]byte(req.Body), &updates); err != nil {
		return badRequest("Invalid request body")
	}

	// Mock update
	user := User{
		ID:    id,
		Name:  updates.Name,
		Email: updates.Email,
	}

	return jsonMarshal(200, Response{
		Message: "User updated successfully",
		Data:    user,
	})
}

// deleteUser deletes a user
func deleteUser(id string) (events.APIGatewayProxyResponse, error) {
	return jsonMarshal(200, Response{
		Message: fmt.Sprintf("User %s deleted successfully", id),
	})
}

// handleHealth returns health status
func handleHealth() (events.APIGatewayProxyResponse, error) {
	return jsonMarshal(200, Response{
		Message: "OK",
		Data: map[string]string{
			"status": "healthy",
		},
	})
}

// Helper functions for HTTP responses

func jsonMarshal(statusCode int, data interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Failed to marshal response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods":     "GET,OPTIONS,POST,PUT,DELETE",
			"Access-Control-Max-Age":           "86400",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(body),
	}, nil
}

func badRequest(message string) (events.APIGatewayProxyResponse, error) {
	return jsonMarshal(400, Response{
		Error: message,
	})
}

func notFound() (events.APIGatewayProxyResponse, error) {
	return jsonMarshal(404, Response{
		Error: "Not found",
	})
}

func methodNotAllowed() (events.APIGatewayProxyResponse, error) {
	return jsonMarshal(405, Response{
		Error: "Method not allowed",
	})
}

func main() {
	lambda.Start(router)
}

/*
REST API Endpoints:

GET    /users        - Get all users
POST   /users        - Create user
GET    /users/{id}   - Get user by ID
PUT    /users/{id}   - Update user
DELETE /users/{id}   - Delete user
GET    /health       - Health check

Testing with curl:

# Get all users
curl http://localhost:3000/users

# Create user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie","email":"charlie@example.com"}'

# Get user by ID
curl http://localhost:3000/users/1

# Update user
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","email":"alice.updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:3000/users/1

# Health check
curl http://localhost:3000/health

API Gateway Integration:

1. In AWS Console:
   - Create REST API
   - Create resources: /users, /users/{id}, /health
   - Create methods: GET, POST, PUT, DELETE
   - Set Lambda proxy integration
   - Deploy API to stage

2. Or use AWS SAM (see template.yaml):
   sam deploy

3. Test with SAM CLI:
   sam local start-api
   # Test endpoints above
*/
