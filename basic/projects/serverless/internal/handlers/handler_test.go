package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleHealth(t *testing.T) {
	result, err := HandleHealth(context.Background())
	if err != nil {
		t.Fatalf("HandleHealth returned error: %v", err)
	}

	var resp HealthResponse
	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", resp.Status)
	}

	if resp.Version == "" {
		t.Error("Expected non-empty version")
	}

	if resp.Time == "" {
		t.Error("Expected non-empty time")
	}
}

func TestHandleRequest_HealthPath(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/health",
		HTTPMethod: "GET",
	}

	resp, err := HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleRequest returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", resp.Headers["Content-Type"])
	}
}

func TestHandleRequest_EventPath(t *testing.T) {
	eventBody := `{"type":"lesson.completed","timestamp":"2024-01-01T00:00:00Z","payload":{"user_id":"user123","lesson_id":"lesson456","score":95,"completed":true}}`

	req := events.APIGatewayProxyRequest{
		Path:       "/event",
		HTTPMethod: "POST",
		Body:       eventBody,
	}

	resp, err := HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleRequest returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestHandleRequest_UnknownPath(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/unknown",
		HTTPMethod: "GET",
	}

	resp, err := HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleRequest returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestHandleRequest_InvalidEventBody(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/event",
		HTTPMethod: "POST",
		Body:       "invalid json",
	}

	resp, err := HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleRequest returned error: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

// HealthResponse mirrors the internal models
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}

func TestGetVersion(t *testing.T) {
	// Test default version
	version := getVersion()
	if version == "" {
		t.Error("Expected non-empty version")
	}
}
