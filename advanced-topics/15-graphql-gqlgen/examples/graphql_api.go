package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ============================================
// Local Types (for demonstration purposes)
// ============================================

// User represents a user in the system (local copy for examples)
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// ============================================
// Complete GraphQL API Example
// ============================================

// This example demonstrates a complete GraphQL API with:
// - Schema definition
// - Query and mutation resolvers
// - Authentication middleware
// - Error handling
// - Data loading patterns
// - Subscription support

// Example 1: Simple Query
func exampleSimpleQuery() {
	fmt.Println("\n=== Example 1: Simple Query ===")

	query := `
		query {
			users {
				id
				username
				email
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)

	// Result:
	// {
	//   "data": {
	//     "users": [
	//       { "id": "1", "username": "admin", "email": "admin@example.com" },
	//       { "id": "2", "username": "alice", "email": "alice@example.com" }
	//     ]
	//   }
	// }
}

// Example 2: Query with Arguments
func exampleQueryWithArguments() {
	fmt.Println("\n=== Example 2: Query with Arguments ===")

	query := `
		query {
			users(role: USER, active: true, page: {limit: 10, offset: 0}) {
				id
				username
				email
				role
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)

	// Result: Only active users with USER role
}

// Example 3: Nested Query
func exampleNestedQuery() {
	fmt.Println("\n=== Example 3: Nested Query ===")

	query := `
		query {
			user(id: "2") {
				id
				username
				email
				posts {
					id
					title
					content
					comments {
						id
						content
						author {
							username
						}
					}
				}
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)

	// Result: User with posts and comments
}

// Example 4: Mutation
func exampleMutation() {
	fmt.Println("\n=== Example 4: Mutation ===")

	mutation := `
		mutation CreateUser {
			createUser(input: {
				username: "bob"
				email: "bob@example.com"
				password: "password123"
				role: USER
			}) {
				id
				username
				email
				role
				createdAt
			}
		}
	`

	fmt.Printf("Mutation:\n%s\n", mutation)

	// Result: Newly created user
}

// Example 5: Mutation with Error Handling
func exampleMutationWithError() {
	fmt.Println("\n=== Example 5: Mutation with Error ===")

	mutation := `
		mutation CreateUser {
			createUser(input: {
				username: ""
				email: "invalid-email"
				password: "123"
				role: USER
			}) {
				id
				username
			}
		}
	`

	fmt.Printf("Mutation:\n%s\n", mutation)

	// Result:
	// {
	//   "data": { "createUser": null },
	//   "errors": [
	//     {
	//       "message": "username and email are required",
	//       "path": ["createUser"]
	//     }
	//   ]
	// }
}

// Example 6: Authentication
func exampleAuthenticatedMutation() {
	fmt.Println("\n=== Example 6: Authenticated Mutation ===")

	// First, authenticate and get token
	token := authenticate("admin", "password")

	// Then, use token in mutation
	mutation := `
		mutation CreateUser {
			createUser(input: {
				username: "charlie"
				email: "charlie@example.com"
				password: "password123"
				role: USER
			}) {
				id
				username
				email
			}
		}
	`

	fmt.Printf("Mutation with Auth:\n%s\n", mutation)
	fmt.Printf("Authorization: Bearer %s\n", token)

	// Request includes header:
	// Authorization: Bearer <token>
}

// Example 7: Subscription
func exampleSubscription() {
	fmt.Println("\n=== Example 7: Subscription ===")

	subscription := `
		subscription PostPublished {
			postPublished {
				id
				title
				author {
					username
				}
			}
		}
	`

	fmt.Printf("Subscription:\n%s\n", subscription)

	// Result: Stream of published posts
	// {
	//   "data": {
	//     "postPublished": {
	//       "id": "1",
	//       "title": "New Post",
	//       "author": { "username": "alice" }
	//     }
	//   }
	// }
	// ... (more events as posts are published)
}

// Example 8: Query with Fragments
func exampleQueryWithFragments() {
	fmt.Println("\n=== Example 8: Query with Fragments ===")

	query := `
		fragment UserFields on User {
			id
			username
			email
		}

		fragment PostFields on Post {
			id
			title
			content
			published
		}

		query GetUserData {
			user(id: "2") {
				...UserFields
				posts {
					...PostFields
				}
			}
		}
	`

	fmt.Printf("Query with Fragments:\n%s\n", query)

	// Result: Reusable field definitions
}

// Example 9: Query with Directives
func exampleQueryWithDirectives() {
	fmt.Println("\n=== Example 9: Query with Directives ===")

	query := `
		query GetUser($withEmail: Boolean!) {
			user(id: "2") {
				id
				username
				email @include(if: $withEmail)
				posts @skip(if: false) {
					id
					title
				}
			}
		}
	`

	fmt.Printf("Query with Directives:\n%s\n", query)

	// Variables: { "withEmail": true }
	// Result: Conditional field inclusion
}

// Example 10: Multiple Queries in One Request
func exampleBatchQueries() {
	fmt.Println("\n=== Example 10: Batch Queries ===")

	query := `
		query {
			users(limit: 5) {
				id
				username
			}
			posts(limit: 5) {
				id
				title
			}
			stats {
				users
				posts
				publishedPosts
			}
		}
	`

	fmt.Printf("Batch Query:\n%s\n", query)

	// Result: Multiple queries in single request
	// More efficient than multiple HTTP requests
}

// Example 11: Error Handling
func exampleErrorHandling() {
	fmt.Println("\n=== Example 11: Error Handling ===")

	query := `
		query {
			user(id: "999") {
				id
				username
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)

	// Result:
	// {
	//   "data": { "user": null },
	//   "errors": [
	//     {
	//       "message": "user not found: 999",
	//       "path": ["user"]
	//     }
	//   ]
	// }

	// Note: Partial success - some fields may succeed even if others fail
}

// Example 12: Introspection Query
func exampleIntrospection() {
	fmt.Println("\n=== Example 12: Introspection ===")

	query := `
		query {
			__schema {
				types {
					name
					kind
					description
					fields {
						name
					}
				}
			}
		}
	`

	fmt.Printf("Introspection Query:\n%s\n", query)

	// Result: Schema information
	// Useful for tooling and documentation
}

// ============================================
// HTTP Client Example
// ============================================

func executeGraphQLQuery(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	// Build request body
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer " + token)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// ============================================
// Authentication Example
// ============================================

type AuthResponse struct {
	Token string `json:"token"`
	User  *User `json:"user"`
}

func authenticate(username, password string) string {
	// In production, validate credentials against database
	// and generate JWT token
	return "mock-jwt-token"
}

// ============================================
// WebSocket Client Example (Subscriptions)
// ============================================

func exampleSubscriptionClient() {
	fmt.Println("\n=== WebSocket Subscription Client ===")

	// Connect to WebSocket endpoint
	// conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/subscription", nil)

	// Send subscription
	subscription := `
		subscription PostPublished {
			postPublished {
				id
				title
				author {
					username
				}
			}
		}
	`

	fmt.Printf("Subscription:\n%s\n", subscription)

	// Listen for events
	// for {
	//     var event map[string]interface{}
	//     err := conn.ReadJSON(&event)
	//     if err != nil {
	//         log.Printf("Error: %v", err)
	//         break
	//     }
	//     log.Printf("Event: %v", event)
	// }
}

// ============================================
// Main Function
// ============================================

func main() {
	fmt.Println("🚀 GraphQL API Examples")
	fmt.Println("========================")

	// Run examples
	exampleSimpleQuery()
	exampleQueryWithArguments()
	exampleNestedQuery()
	exampleMutation()
	exampleMutationWithError()
	exampleAuthenticatedMutation()
	exampleSubscription()
	exampleQueryWithFragments()
	exampleQueryWithDirectives()
	exampleBatchQueries()
	exampleErrorHandling()
	exampleIntrospection()

	fmt.Println("\n✅ All examples completed")
	fmt.Println("\n📖 Next Steps:")
	fmt.Println("1. Run the server: go run server.go")
	fmt.Println("2. Open GraphQL Playground: http://localhost:8080")
	fmt.Println("3. Try these queries interactively")
	fmt.Println("4. Explore the schema using introspection")
}

/*
Testing Examples:

1. Start the server:
   go run server.go

2. Open GraphQL Playground:
   open http://localhost:8080

3. Try queries from examples

4. Use curl for testing:
   curl -X POST http://localhost:8080/query \
     -H "Content-Type: application/json" \
     -d '{"query": "{ users { id username } }"}'

Common Queries:

# Get all users
query { users { id username email role } }

# Get user with posts
query { user(id: "2") { id username posts { title } } }

# Create user (requires auth)
mutation { createUser(input: {username: "bob", email: "bob@example.com", password: "pass", role: USER}) { id username } }

# Search posts
query { searchPosts(query: "graphql", limit: 5) { id title } }

# Get statistics
query { stats { users posts publishedPosts } }

Performance Tips:

1. Use query batching to reduce round trips
2. Implement data loaders to prevent N+1 queries
3. Add query complexity analysis
4. Use persisted queries for security
5. Implement caching at resolver level
6. Use pagination for large lists
7. Add rate limiting
8. Monitor query performance

Security Best Practices:

1. Always validate user input
2. Check authentication/authorization in resolvers
3. Use HTTPS in production
4. Implement rate limiting
5. Add query depth limiting
6. Use query cost analysis
7. Sanitize error messages
8. Keep dependencies updated
*/
