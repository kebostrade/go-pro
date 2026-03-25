package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/v2/ast"
)

// ============================================
// Server Setup
// ============================================

func main() {
	// Create router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(authMiddleware)

	// GraphQL playground (development only)
	if os.Getenv("ENV") != "production" {
		router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	}

	// GraphQL endpoint
	router.Handle("/query", GraphQLHandler())

	// Health check
	router.Get("/health", healthCheckHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server ready at http://localhost:%s", port)
	log.Printf("📖 GraphQL Playground at http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ============================================
// GraphQL Handler
// ============================================

func GraphQLHandler() *handler.Server {
	// Create GraphQL server
	srv := handler.New(&executableSchema{
		resolvers: &resolvers{
			db: NewMockDB(),
		},
	})

	// Add custom middleware
	srv.Use(&authMiddlewareGraphQL{})
	srv.Use(&loggingMiddlewareGraphQL{})
	srv.Use(&errorHandlingMiddlewareGraphQL{})

	return srv
}

// ============================================
// HTTP Middleware
// ============================================

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		token := r.Header.Get("Authorization")

		// Validate token (in production, verify JWT)
		if token == "" {
			// Allow anonymous access
			next.ServeHTTP(w, r)
			return
		}

		// For demo, extract user ID from token
		// In production, validate JWT and extract claims
		userID := "2" // Mock authenticated user

		// Get user from database
		db := NewMockDB()
		user, err := db.UserByID(userID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ============================================
// GraphQL Middleware
// ============================================

type authMiddlewareGraphQL struct{}

func (authMiddlewareGraphQL) ExtensionName() string {
	return "AuthMiddleware"
}

func (authMiddlewareGraphQL) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (authMiddlewareGraphQL) Middleware(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	// Get operation type
	operation := graphql.GetOperationContext(ctx)

	// Allow queries and subscriptions without auth
	if operation.Operation != nil && (operation.Operation.Operation == "query" || operation.Operation.Operation == "subscription") {
		return next(ctx)
	}

	// Require auth for mutations
	user := UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	return next(ctx)
}

type loggingMiddlewareGraphQL struct{}

func (loggingMiddlewareGraphQL) ExtensionName() string {
	return "LoggingMiddleware"
}

func (loggingMiddlewareGraphQL) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (loggingMiddlewareGraphQL) Middleware(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	start := time.Now()

	// Get operation details
	operation := graphql.GetOperationContext(ctx)
	field := graphql.GetFieldContext(ctx)

	log.Printf(
		"GraphQL: %s %s.%s",
		operation.Operation,
		operation.OperationName,
		field.Field.Name,
	)

	// Call next resolver
	res, err := next(ctx)

	// Log duration
	log.Printf(
		"GraphQL: %s %s.%s took %v",
		operation.Operation,
		operation.OperationName,
		field.Field.Name,
		time.Since(start),
	)

	return res, err
}

type errorHandlingMiddlewareGraphQL struct{}

func (errorHandlingMiddlewareGraphQL) ExtensionName() string {
	return "ErrorHandlingMiddleware"
}

func (errorHandlingMiddlewareGraphQL) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (errorHandlingMiddlewareGraphQL) Middleware(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	// Call next resolver
	res, err := next(ctx)

	// Handle errors
	if err != nil {
		// Log error
		log.Printf("GraphQL error: %v", err)

		// Return custom error extensions
		return nil, graphql.ErrorOnPath(ctx, err)
	}

	return res, nil
}

// ============================================
// WebSocket Handler for Subscriptions
// ============================================

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Process subscription message
		log.Printf("Received message: %s", message)

		// Echo back for demo
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

// ============================================
// Health Check Handler
// ============================================

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// ============================================
// CORS Middleware
// ============================================

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ============================================
// GraphQL Schema
// ============================================

type executableSchema struct {
	resolvers *resolvers
}

// Custom resolver interfaces for gqlgen
type QueryResolverI interface {
	User(ctx context.Context, id string) (*User, error)
	Users(ctx context.Context, role *Role, active *bool, page *PageInput, sort *PostSort) ([]*User, error)
	Me(ctx context.Context) (*User, error)
	Post(ctx context.Context, id string) (*Post, error)
	Posts(ctx context.Context, authorId *string, published *bool, tags []string, search *string, page *PageInput, sort *PostSort) ([]*Post, error)
	Comment(ctx context.Context, id string) (*Comment, error)
	Comments(ctx context.Context, postId string) ([]*Comment, error)
	SearchUsers(ctx context.Context, query string, limit *int) ([]*User, error)
	SearchPosts(ctx context.Context, query string, limit *int) ([]*Post, error)
	Stats(ctx context.Context) (*Stats, error)
}

type MutationResolverI interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	DeactivateUser(ctx context.Context, id string) (bool, error)
	CreatePost(ctx context.Context, input CreatePostInput) (*Post, error)
	UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*Post, error)
	DeletePost(ctx context.Context, id string) (bool, error)
	PublishPost(ctx context.Context, id string) (*Post, error)
	UnpublishPost(ctx context.Context, id string) (*Post, error)
	CreateComment(ctx context.Context, input CreateCommentInput) (*Comment, error)
	DeleteComment(ctx context.Context, id string) (bool, error)
	UploadImage(ctx context.Context, file interface{}) (string, error)
}

type UserResolverI interface {
	Posts(ctx context.Context, obj *User, limit *int, offset *int) ([]*Post, error)
}

type PostResolverI interface {
	Author(ctx context.Context, obj *Post) (*User, error)
	Comments(ctx context.Context, obj *Post) ([]*Comment, error)
}

type CommentResolverI interface {
	Author(ctx context.Context, obj *Comment) (*User, error)
	Post(ctx context.Context, obj *Comment) (*Post, error)
}

type SubscriptionResolverI interface{}

func (e *executableSchema) Query() QueryResolverI {
	return e.resolvers.Query()
}

func (e *executableSchema) Mutation() MutationResolverI {
	return e.resolvers.Mutation()
}

func (e *executableSchema) User() UserResolverI {
	return e.resolvers.User()
}

func (e *executableSchema) Post() PostResolverI {
	return e.resolvers.Post()
}

func (e *executableSchema) Comment() CommentResolverI {
	return e.resolvers.Comment()
}

func (e *executableSchema) Subscription() SubscriptionResolverI {
	return nil
}

// Schema returns the GraphQL schema (required by graphql.ExecutableSchema)
func (e *executableSchema) Schema() *ast.Schema {
	return nil
}

// Complexity is required by graphql.ExecutableSchema
func (e *executableSchema) Complexity(typeName string, fieldName string, childComplexity int, args map[string]any) (int, bool) {
	return 1, true
}

// Exec is required by graphql.ExecutableSchema - it returns a ResponseHandler for executing queries
func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	// Return a basic response handler that executes the query
	// In a real implementation, this would parse and execute the GraphQL operation
	return func(ctx context.Context) *graphql.Response {
		return &graphql.Response{
			Data: []byte("{}"),
		}
	}
}

/*
Server Configuration:

Environment Variables:
- PORT: Server port (default: 8080)
- ENV: Environment (development/production)
- JWT_SECRET: Secret for JWT validation
- DATABASE_URL: Database connection string

Running the Server:

1. Start server:
   go run server.go

2. Access GraphQL Playground:
   open http://localhost:8080

3. Health check:
   curl http://localhost:8080/health

Production Deployment:

1. Enable production mode:
   export ENV=production
   export PORT=8080
   go run server.go

2. Use TLS:
   - Put nginx/caddy in front
   - Or use Go's http.ListenAndServeTLS

3. Process management:
   - Use systemd/supervisor
   - Or Docker/Kubernetes

Testing:

# Query
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ users { id username } }"}'

# Mutation with auth
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"query":"mutation { createUser(input: {...}) { id username } }"}'
*/
