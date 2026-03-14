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
	srv.Use(authMiddlewareGraphQL)
	srv.Use(loggingMiddlewareGraphQL)
	srv.Use(errorHandlingMiddlewareGraphQL)

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

func authMiddlewareGraphQL(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	// Get operation type
	operation := graphql.GetOperationContext(ctx)

	// Allow queries and subscriptions without auth
	if operation.Operation == "query" || operation.Operation == "subscription" {
		return next(ctx)
	}

	// Require auth for mutations
	user := UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	return next(ctx)
}

func loggingMiddlewareGraphQL(ctx context.Context, next graphql.Resolver) (interface{}, error) {
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

func errorHandlingMiddlewareGraphQL(ctx context.Context, next graphql.Resolver) (interface{}, error) {
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

func (e *executableSchema) Query() graphql.QueryResolver {
	return e.resolvers.Query()
}

func (e *executableSchema) Mutation() graphql.MutationResolver {
	return e.resolvers.Mutation()
}

func (e *executableSchema) User() graphql.UserResolver {
	if userResolver, ok := e.resolvers.User().(graphql.UserResolver); ok {
		return userResolver
	}
	return nil
}

func (e *executableSchema) Post() graphql.PostResolver {
	if postResolver, ok := e.resolvers.Post().(graphql.PostResolver); ok {
		return postResolver
	}
	return nil
}

func (e *executableSchema) Comment() graphql.CommentResolver {
	if commentResolver, ok := e.resolvers.Comment().(graphql.CommentResolver); ok {
		return commentResolver
	}
	return nil
}

func (e *executableSchema) Subscription() graphql.SubscriptionResolver {
	return nil
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
