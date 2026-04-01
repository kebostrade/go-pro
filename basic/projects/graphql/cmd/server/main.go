package main

import (
	"context"
	"encoding/json"
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

	"basic/projects/graphql/internal/graph"
	"basic/projects/graphql/pkg/auth"
	"basic/projects/graphql/pkg/models"
)

func main() {
	// Create database and resolver
	db := models.NewDB()
	resolver := graph.NewResolver(db)

	// Create router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Timeout(60 * time.Second))

	// Auth middleware
	router.Use(authMiddleware)

	// GraphQL playground (development only)
	if os.Getenv("ENV") != "production" {
		router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	}

	// GraphQL endpoint
	router.Handle("/query", GraphQLHandler(resolver))

	// WebSocket endpoint for subscriptions
	router.HandleFunc("/subscriptions", subscriptionHandler(resolver))

	// Health check
	router.Get("/health", healthCheckHandler)

	// CORS middleware
	router.Use(corsMiddleware)

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

func GraphQLHandler(resolver *graph.Resolver) *handler.Server {
	// Create GraphQL server
	srv := handler.New(&executableSchema{
		resolvers: resolver,
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

		// Validate token
		if token != "" {
			// Remove "Bearer " prefix
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}

			// Validate JWT
			claims, err := auth.ValidateJWT(token)
			if err == nil {
				// Add user to context
				ctx := auth.ContextWithUser(r.Context(), claims.UserID)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
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
	if operation.Operation != nil && (operation.Operation.Operation == ast.Query || operation.Operation.Operation == ast.Subscription) {
		return next(ctx)
	}

	// Require auth for mutations
	userID := graph.UserFromContext(ctx)
	if userID == "" {
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
		operation.Operation.Operation,
		operation.OperationName,
		field.Field.Name,
	)

	// Call next resolver
	res, err := next(ctx)

	// Log duration
	log.Printf(
		"GraphQL: %s %s.%s took %v",
		operation.Operation.Operation,
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

func subscriptionHandler(resolver *graph.Resolver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		// Handle WebSocket messages for subscriptions
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}

			// Parse subscription request
			var req subscriptionRequest
			if err := json.Unmarshal(message, &req); err != nil {
				log.Printf("Failed to parse subscription request: %v", err)
				continue
			}

			// Handle subscription
			go handleSubscription(conn, resolver, req)
		}
	}
}

type subscriptionRequest struct {
	Type    string          `json:"type"`
	ID      string          `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

type subscriptionPayload struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func handleSubscription(conn *websocket.Conn, resolver *graph.Resolver, req subscriptionRequest) {
	if req.Type == "connection_init" {
		// Send connection acknowledgment
		conn.WriteJSON(map[string]string{"type": "connection_ack"})
		return
	}

	if req.Type == "subscribe" {
		var payload subscriptionPayload
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			log.Printf("Failed to parse subscription payload: %v", err)
			return
		}

		// Start subscription based on operation
		_ = context.Background() // ctx for future subscription implementation

		// Send initial result
		conn.WriteJSON(map[string]interface{}{
			"type": "next",
			"id":   req.ID,
			"payload": map[string]interface{}{
				"data": map[string]interface{}{
					"userCreated": map[string]interface{}{
						"id":       "1",
						"username": "test",
					},
				},
			},
		})
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
	resolvers *graph.Resolver
}

// Query resolvers
func (e *executableSchema) Query() *queryResolver {
	return &queryResolver{e.resolvers}
}

type queryResolver struct {
	r *graph.Resolver
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	return r.r.User(ctx, id)
}

func (r *queryResolver) Users(ctx context.Context, role *models.Role, active *bool, page *graph.PageInput) ([]*models.User, error) {
	return r.r.Users(ctx, role, active, page)
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	return r.r.Me(ctx)
}

func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	return r.r.Post(ctx, id)
}

func (r *queryResolver) Posts(ctx context.Context, authorID *string, published *bool, tags []string, search *string, page *graph.PageInput) ([]*models.Post, error) {
	return r.r.Posts(ctx, authorID, published, tags, search, page)
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*models.Comment, error) {
	return r.r.GetComment(ctx, id)
}

func (r *queryResolver) Comments(ctx context.Context, postID string) ([]*models.Comment, error) {
	return r.r.Comments(ctx, postID)
}

func (r *queryResolver) Stats(ctx context.Context) (*models.Stats, error) {
	return r.r.Stats(ctx)
}

// Mutation resolvers
func (e *executableSchema) Mutation() *mutationResolver {
	return &mutationResolver{e.resolvers}
}

type mutationResolver struct {
	r *graph.Resolver
}

func (r *mutationResolver) CreateUser(ctx context.Context, input graph.CreateUserInput) (*models.User, error) {
	return r.r.CreateUser(ctx, input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input graph.UpdateUserInput) (*models.User, error) {
	return r.r.UpdateUser(ctx, id, input)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	return r.r.DeleteUser(ctx, id)
}

func (r *mutationResolver) DeactivateUser(ctx context.Context, id string) (bool, error) {
	return r.r.DeactivateUser(ctx, id)
}

func (r *mutationResolver) CreatePost(ctx context.Context, input graph.CreatePostInput) (*models.Post, error) {
	return r.r.CreatePost(ctx, input)
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input graph.UpdatePostInput) (*models.Post, error) {
	return r.r.UpdatePost(ctx, id, input)
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	return r.r.DeletePost(ctx, id)
}

func (r *mutationResolver) PublishPost(ctx context.Context, id string) (*models.Post, error) {
	return r.r.PublishPost(ctx, id)
}

func (r *mutationResolver) UnpublishPost(ctx context.Context, id string) (*models.Post, error) {
	return r.r.UnpublishPost(ctx, id)
}

func (r *mutationResolver) CreateComment(ctx context.Context, input graph.CreateCommentInput) (*models.Comment, error) {
	return r.r.CreateComment(ctx, input)
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	return r.r.DeleteComment(ctx, id)
}

// Field resolvers
func (e *executableSchema) User() *userResolver {
	return &userResolver{e.resolvers}
}

type userResolver struct {
	r *graph.Resolver
}

func (r *userResolver) Posts(ctx context.Context, obj *models.User, limit *int, offset *int) ([]*models.Post, error) {
	return r.r.UserPosts(ctx, obj, limit, offset)
}

func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (time.Time, error) {
	return obj.CreatedAt, nil
}

func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (time.Time, error) {
	return obj.UpdatedAt, nil
}

func (e *executableSchema) Post() *postResolver {
	return &postResolver{e.resolvers}
}

type postResolver struct {
	r *graph.Resolver
}

func (r *postResolver) Author(ctx context.Context, obj *models.Post) (*models.User, error) {
	return r.r.PostAuthor(ctx, obj)
}

func (r *postResolver) Comments(ctx context.Context, obj *models.Post) ([]*models.Comment, error) {
	return r.r.PostComments(ctx, obj)
}

func (r *postResolver) CreatedAt(ctx context.Context, obj *models.Post) (time.Time, error) {
	return obj.CreatedAt, nil
}

func (r *postResolver) UpdatedAt(ctx context.Context, obj *models.Post) (time.Time, error) {
	return obj.UpdatedAt, nil
}

func (e *executableSchema) Comment() *commentResolver {
	return &commentResolver{e.resolvers}
}

type commentResolver struct {
	r *graph.Resolver
}

func (r *commentResolver) Author(ctx context.Context, obj *models.Comment) (*models.User, error) {
	return r.r.CommentAuthor(ctx, obj)
}

func (r *commentResolver) Post(ctx context.Context, obj *models.Comment) (*models.Post, error) {
	return r.r.CommentPost(ctx, obj)
}

func (r *commentResolver) CreatedAt(ctx context.Context, obj *models.Comment) (time.Time, error) {
	return obj.CreatedAt, nil
}

// Subscription resolvers
func (e *executableSchema) Subscription() *subscriptionResolver {
	return &subscriptionResolver{e.resolvers}
}

type subscriptionResolver struct {
	r *graph.Resolver
}

func (r *subscriptionResolver) UserCreated(ctx context.Context) (<-chan *models.User, error) {
	return r.r.UserCreated(ctx)
}

func (r *subscriptionResolver) PostCreated(ctx context.Context) (<-chan *models.Post, error) {
	return r.r.PostCreated(ctx)
}

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	return r.r.CommentAdded(ctx, postID)
}

// Schema returns the GraphQL schema
func (e *executableSchema) Schema() *ast.Schema {
	return nil
}

// Complexity is required by graphql.ExecutableSchema
func (e *executableSchema) Complexity(typeName string, fieldName string, childComplexity int, args map[string]interface{}) (int, bool) {
	return 1, true
}

// Exec is required by graphql.ExecutableSchema
func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	// This is a placeholder - gqlgen generates the actual implementation
	return func(ctx context.Context) *graphql.Response {
		return &graphql.Response{
			Data: []byte("{}"),
		}
	}
}
