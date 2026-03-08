package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handleGet(ctx, req)
	case "POST":
		return handlePost(ctx, req)
	case "PUT":
		return handlePut(ctx, req)
	case "DELETE":
		return handleDelete(ctx, req)
	default:
		return response(405, `{"error":"method not allowed"}`), nil
	}
}

func handleGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]

	if id == "" {
		users := []User{
			{ID: "1", Name: "John", Email: "john@example.com"},
			{ID: "2", Name: "Jane", Email: "jane@example.com"},
		}
		data, _ := json.Marshal(users)
		return response(200, string(data)), nil
	}

	user := User{ID: id, Name: "John Doe", Email: "john@example.com"}
	data, _ := json.Marshal(user)
	return response(200, string(data)), nil
}

func handlePost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return response(400, `{"error":"invalid JSON"}`), nil
	}

	if user.Name == "" || user.Email == "" {
		return response(400, `{"error":"name and email required"}`), nil
	}

	user.ID = "generated-id"
	data, _ := json.Marshal(user)
	return response(201, string(data)), nil
}

func handlePut(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	if id == "" {
		return response(400, `{"error":"id required"}`), nil
	}

	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return response(400, `{"error":"invalid JSON"}`), nil
	}

	user.ID = id
	data, _ := json.Marshal(user)
	return response(200, string(data)), nil
}

func handleDelete(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	if id == "" {
		return response(400, `{"error":"id required"}`), nil
	}

	return response(204, ""), nil
}

func response(status int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}
}

func main() {
	if os.Getenv("LOCAL") == "true" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			req := events.APIGatewayProxyRequest{
				HTTPMethod: r.Method,
				Path:       r.URL.Path,
				Body:       readBody(r),
			}

			resp, _ := handleRequest(r.Context(), req)
			w.WriteHeader(resp.StatusCode)
			w.Write([]byte(resp.Body))
		})
		log.Println("Running locally on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		lambda.Start(handleRequest)
	}
}

func readBody(r *http.Request) string {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	return string(body)
}
