package handlers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/gopro/basic/projects/serverless/internal/models"
)

// HandleHealth handles health check requests
func HandleHealth(ctx context.Context) (string, error) {
	resp := models.HealthResponse{
		Status:  "healthy",
		Version: getVersion(),
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// HandleRequest handles incoming Lambda URL requests
// Lambda URLs send APIGatewayProxyRequest events
func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	log.Printf("Received request for Lambda function: %s, path: %s, method: %s",
		functionName, req.Path, req.HTTPMethod)

	// Route based on path
	switch req.Path {
	case "/health":
		return handleHealth(ctx)
	case "/event":
		return handleEvent(ctx, req)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error":"not found"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}

// handleHealth returns health check response
func handleHealth(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	resp := models.HealthResponse{
		Status:  "healthy",
		Version: getVersion(),
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// handleEvent processes incoming events
func handleEvent(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event models.Event

	if err := json.Unmarshal([]byte(req.Body), &event); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"invalid event format"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	log.Printf("Processing event of type: %s", event.Type)

	// Process event based on type
	var response map[string]string
	switch event.Type {
	case "lesson.completed":
		response = map[string]string{
			"status":  "processed",
			"type":    event.Type,
			"message": "Lesson completion recorded",
		}
	case "course.enrolled":
		response = map[string]string{
			"status":  "processed",
			"type":    event.Type,
			"message": "Course enrollment recorded",
		}
	default:
		response = map[string]string{
			"status":  "ignored",
			"type":    event.Type,
			"message": "Unknown event type",
		}
	}

	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// getVersion returns the application version
func getVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "1.0.0"
	}
	return version
}
