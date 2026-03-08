package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// User represents a user in DynamoDB
type User struct {
	ID        string    `dynamodbav:"id" json:"id"`
	Name      string    `dynamodbav:"name" json:"name"`
	Email     string    `dynamodbav:"email" json:"email"`
	CreatedAt time.Time `dynamodbav:"created_at" json:"created_at"`
	UpdatedAt time.Time `dynamodbav:"updated_at" json:"updated_at"`
}

// DynamoDBClient wraps the DynamoDB service
type DynamoDBClient struct {
	client *dynamodb.Client
	table  string
}

// NewDynamoDBClient creates a new DynamoDB client
func NewDynamoDBClient(ctx context.Context) (*DynamoDBClient, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Get table name from environment
	table := os.Getenv("DYNAMODB_TABLE")
	if table == "" {
		table = "Users" // Default table name
	}

	return &DynamoDBClient{
		client: dynamodb.NewFromConfig(cfg),
		table:  table,
	}, nil
}

// CreateUser creates a new user in DynamoDB
func (dc *DynamoDBClient) CreateUser(ctx context.Context, user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Marshal Go type to DynamoDB attribute values
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Put item in DynamoDB
	_, err = dc.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dc.table),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

// GetUser retrieves a user by ID
func (dc *DynamoDBClient) GetUser(ctx context.Context, id string) (*User, error) {
	// Get item from DynamoDB
	result, err := dc.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(dc.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	// Check if item exists
	if result.Item == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Unmarshal DynamoDB attribute values to Go type
	var user User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// ListUsers retrieves all users
func (dc *DynamoDBClient) ListUsers(ctx context.Context) ([]User, error) {
	// Scan table (for production, use Query with index)
	result, err := dc.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(dc.table),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	// Unmarshal items
	var users []User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (dc *DynamoDBClient) UpdateUser(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()

	// Build update expression
	updateExpr := "SET #name = :name, #email = :email, #updated_at = :updated_at"
 exprNames := map[string]string{
		"#name":       "name",
		"#email":      "email",
		"#updated_at": "updated_at",
	}
 exprValues := map[string]types.AttributeValue{
		":name":       &types.AttributeValueMemberS{Value: user.Name},
		":email":      &types.AttributeValueMemberS{Value: user.Email},
		":updated_at": &types.AttributeValueMemberS{Value: user.UpdatedAt.Format(time.RFC3339)},
	}

	// Update item in DynamoDB
	_, err := dc.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:    aws.String(dc.table),
		Key:          map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: user.ID}},
		UpdateExpression: aws.String(updateExpr),
		ExpressionAttributeNames: exprNames,
		ExpressionAttributeValues: exprValues,
	})

	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// DeleteUser deletes a user by ID
func (dc *DynamoDBClient) DeleteUser(ctx context.Context, id string) error {
	// Delete item from DynamoDB
	_, err := dc.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(dc.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// Lambda handler function
func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s", req.HTTPMethod, req.Path)

	// Initialize DynamoDB client
	db, err := NewDynamoDBClient(ctx)
	if err != nil {
		return errorResponse(500, fmt.Sprintf("Failed to initialize DynamoDB: %v", err)), nil
	}

	// Route request
	switch req.Path {
	case "/users":
		if req.HTTPMethod == "GET" {
			return listUsers(ctx, db)
		} else if req.HTTPMethod == "POST" {
			return createUser(ctx, db, req)
		}
	case "/users/{id}":
		if req.HTTPMethod == "GET" {
			return getUser(ctx, db, req)
		} else if req.HTTPMethod == "PUT" {
			return updateUser(ctx, db, req)
		} else if req.HTTPMethod == "DELETE" {
			return deleteUser(ctx, db, req)
		}
	}

	return notFound()
}

// Handlers
func listUsers(ctx context.Context, db *DynamoDBClient) (events.APIGatewayProxyResponse, error) {
	users, err := db.ListUsers(ctx)
	if err != nil {
		return errorResponse(500, fmt.Sprintf("Failed to list users: %v", err)), nil
	}

	return jsonResponse(200, map[string]interface{}{
		"users": users,
	})
}

func createUser(ctx context.Context, db *DynamoDBClient, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return errorResponse(400, "Invalid request body"), nil
	}

	// Generate ID (in production, use UUID or DynamoDB auto-increment)
	user.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	if err := db.CreateUser(ctx, &user); err != nil {
		return errorResponse(500, fmt.Sprintf("Failed to create user: %v", err)), nil
	}

	return jsonResponse(201, map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

func getUser(ctx context.Context, db *DynamoDBClient, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]

	user, err := db.GetUser(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			return notFound(), nil
		}
		return errorResponse(500, fmt.Sprintf("Failed to get user: %v", err)), nil
	}

	return jsonResponse(200, map[string]interface{}{
		"user": user,
	})
}

func updateUser(ctx context.Context, db *DynamoDBClient, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]

	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return errorResponse(400, "Invalid request body"), nil
	}

	user.ID = userID
	if err := db.UpdateUser(ctx, &user); err != nil {
		return errorResponse(500, fmt.Sprintf("Failed to update user: %v", err)), nil
	}

	return jsonResponse(200, map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	})
}

func deleteUser(ctx context.Context, db *DynamoDBClient, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]

	if err := db.DeleteUser(ctx, userID); err != nil {
		return errorResponse(500, fmt.Sprintf("Failed to delete user: %v", err)), nil
	}

	return jsonResponse(200, map[string]interface{}{
		"message": "User deleted successfully",
	})
}

// Helper functions
func jsonResponse(statusCode int, data interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return errorResponse(500, "Failed to marshal response"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}, nil
}

func errorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return jsonResponse(statusCode, map[string]string{
		"error": message,
	})
}

func notFound() (events.APIGatewayProxyResponse, error) {
	return errorResponse(404, "Not found"), nil
}

func main() {
	lambda.Start(handleRequest)
}

/*
DynamoDB Table Setup:

1. Create DynamoDB table:
   aws dynamodb create-table \
     --table-name Users \
     --attribute-definitions AttributeName=id,AttributeType=S \
     --key-schema AttributeName=id,KeyType=HASH \
     --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

2. Or use AWS CloudFormation/SAM (see template.yaml)

3. Set environment variable:
   export DYNAMODB_TABLE=Users

4. Test locally:
   sam local invoke DynamoDBHandler --event events/dynamodb_event.json

Testing:

# List users
curl http://localhost:3000/users

# Create user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'

# Get user
curl http://localhost:3000/users/{id}

# Update user
curl -X PUT http://localhost:3000/users/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","email":"alice.updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:3000/users/{id}

IAM Permissions Required:

- dynamodb:GetItem
- dynamodb:PutItem
- dynamodb:Scan
- dynamodb:UpdateItem
- dynamodb:DeleteItem

Best Practices:

1. Use Query instead of Scan for production
2. Add GSI (Global Secondary Index) for alternate queries
3. Use batch operations for multiple items
4. Implement exponential backoff for retries
5. Use DynamoDB Accelerator (DAX) for caching
6. Enable TTL for automatic item expiration
7. Use ON-DEMAND capacity for variable workloads
*/
