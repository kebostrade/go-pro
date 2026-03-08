package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ============================================
// GraphQL Authentication Examples
// ============================================
// This file demonstrates authentication and authorization patterns

// Example 1: Basic Authentication Flow
func demoBasicAuthFlow() {
	fmt.Println("\n=== Demo 1: Basic Auth Flow ===")

	fmt.Println("1. Client sends credentials:")
	fmt.Println("   POST /auth/login")
	fmt.Println("   { username: \"alice\", password: \"password123\" }")

	fmt.Println("\n2. Server validates and returns token:")
	fmt.Println("   { token: \"eyJhbGciOiJIUzI1NiIs...\", user: { id: \"2\", username: \"alice\" } }")

	fmt.Println("\n3. Client includes token in subsequent requests:")
	fmt.Println("   Authorization: Bearer eyJhbGciOiJIUzI1NiIs...")

	fmt.Println("\n4. Server validates token and adds user to context")

	fmt.Println("\nUse case: Simple JWT-based authentication")
}

// Example 2: Authenticated Query
func demoAuthenticatedQuery() {
	fmt.Println("\n=== Demo 2: Authenticated Query ===")

	query := `
		query GetCurrentUser {
			me {
				id
				username
				email
				role
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)
	fmt.Println("\nHeaders:")
	fmt.Println("Authorization: Bearer <jwt-token>")

	fmt.Println("\nResolver implementation:")
	fmt.Println(`
func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	// Get user from context (set by auth middleware)
	user := UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("not authenticated")
	}
	return user, nil
}
`)

	fmt.Println("Use case: Get current user's profile")
}

// Example 3: Authenticated Mutation
func demoAuthenticatedMutation() {
	fmt.Println("\n=== Demo 3: Authenticated Mutation ===")

	mutation := `
		mutation CreatePost($input: CreatePostInput!) {
			createPost(input: $input) {
				id
				title
				author {
					username
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":     "My Post",
			"content":   "Post content",
			"published": true,
		},
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nHeaders:")
	fmt.Println("Authorization: Bearer <jwt-token>")

	fmt.Println("\nResolver implementation:")
	fmt.Println(`
func (r *mutationResolver) CreatePost(ctx context.Context, input CreatePostInput) (*Post, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Create post with authenticated user as author
	return r.db.CreatePost(authUser.ID, input)
}
`)

	fmt.Println("Use case: Mutations that require authentication")
}

// Example 4: Role-Based Authorization
func demoRoleBasedAuth() {
	fmt.Println("\n=== Demo 4: Role-Based Authorization ===")

	mutation := `
		mutation CreateUser($input: CreateUserInput!) {
			createUser(input: $input) {
				id
				username
				role
			}
		}
	`

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Println("\nHeaders:")
	fmt.Println("Authorization: Bearer <admin-jwt-token>")

	fmt.Println("\nResolver implementation:")
	fmt.Println(`
func (r *mutationResolver) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Check if user is admin
	if authUser.Role != RoleAdmin {
		return nil, fmt.Errorf("forbidden: only admins can create users")
	}

	// Create user
	return r.db.CreateUser(input)
}
`)

	fmt.Println("\nUse case: Admin-only operations")
}

// Example 5: Resource-Level Authorization
func demoResourceAuth() {
	fmt.Println("\n=== Demo 5: Resource-Level Authorization ===")

	mutation := `
		mutation UpdatePost($id: ID!, $input: UpdatePostInput!) {
			updatePost(id: $id, input: $input) {
				id
				title
			}
		}
	`

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Println("\nResolver implementation:")
	fmt.Println(`
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*Post, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get post
	post, err := r.db.PostByID(id)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check permission (admin or author)
	if authUser.Role != RoleAdmin && post.AuthorID != authUser.ID {
		return nil, fmt.Errorf("forbidden: you can only edit your own posts")
	}

	// Update post
	return r.db.UpdatePost(id, input)
}
`)

	fmt.Println("\nUse case: Edit permissions for resource owners")
}

// Example 6: Auth Middleware (HTTP)
func demoAuthMiddleware() {
	fmt.Println("\n=== Demo 6: Auth Middleware (HTTP) ===")

	fmt.Println("Middleware implementation:")
	fmt.Println(`
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		token := r.Header.Get("Authorization")
		if token == "" {
			// Allow anonymous access for queries
			next.ServeHTTP(w, r)
			return
		}

		// Validate token (in production, verify JWT)
		if !isValidToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract user from token
		userID := extractUserIDFromToken(token)
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
`)

	fmt.Println("\nUse case: Centralized authentication for HTTP layer")
}

// Example 7: GraphQL Auth Middleware
func demoGraphQLAuthMiddleware() {
	fmt.Println("\n=== Demo 7: GraphQL Auth Middleware ===")

	fmt.Println("Middleware implementation:")
	fmt.Println(`
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
`)

	fmt.Println("\nUse case: Fine-grained auth at resolver level")
}

// Example 8: Token Generation and Validation
func demoTokenHandling() {
	fmt.Println("\n=== Demo 8: Token Generation and Validation ===")

	fmt.Println("Generate token (on login):")
	fmt.Println(`
func generateToken(user *User) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
`)

	fmt.Println("\nValidate token (on request):")
	fmt.Println(`
func validateToken(tokenString string) (*User, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Get user from database
	userID := claims["user_id"].(string)
	return db.UserByID(userID)
}
`)

	fmt.Println("\nUse case: JWT-based authentication")
}

// Example 9: Refresh Token Flow
func demoRefreshToken() {
	fmt.Println("\n=== Demo 9: Refresh Token Flow ===")

	fmt.Println("1. Initial login returns both tokens:")
	fmt.Println("   {")
	fmt.Println("     accessToken: \"eyJhbGciOi...\",  (expires: 15min)")
	fmt.Println("     refreshToken: \"eyJhbGciOi...\",  (expires: 7days)")
	fmt.Println("   }")

	fmt.Println("\n2. When access token expires:")
	fmt.Println("   POST /auth/refresh")
	fmt.Println("   { refreshToken: \"eyJhbGciOi...\" }")

	fmt.Println("\n3. Server validates and returns new tokens:")
	fmt.Println("   {")
	fmt.Println("     accessToken: \"new-access-token\",")
	fmt.Println("     refreshToken: \"new-refresh-token\"")
	fmt.Println("   }")

	fmt.Println("\nUse case: Seamless token renewal without re-login")
}

// Example 10: OAuth2 Integration
func demoOAuth2() {
	fmt.Println("\n=== Demo 10: OAuth2 Integration ===")

	fmt.Println("1. Redirect to OAuth provider:")
	fmt.Println("   GET /auth/oauth?provider=github")
	fmt.Println("   → Redirects to https://github.com/login/oauth/authorize")

	fmt.Println("\n2. User authorizes, provider redirects back:")
	fmt.Println("   GET /auth/oauth/callback?code=...")

	fmt.Println("\n3. Exchange code for access token:")
	fmt.Println("   POST https://github.com/login/oauth/access_token")
	fmt.Println("   { code: \"...\", client_id: \"...\", client_secret: \"...\" }")

	fmt.Println("\n4. Get user info from provider:")
	fmt.Println("   GET https://api.github.com/user")
	fmt.Println("   Authorization: token <access_token>")

	fmt.Println("\n5. Create or update user, generate JWT:")
	fmt.Println("   { token: \"jwt-token\", user: {...} }")

	fmt.Println("\nUse case: Social login (GitHub, Google, etc.)")
}

// ============================================
// HTTP Client Examples
// ============================================

func executeAuthenticatedRequest(query, token string) {
	fmt.Println("\nExample HTTP request:")
	fmt.Printf(`
POST /query
Authorization: Bearer %s
Content-Type: application/json

{
  "query": %s
}
`, token, query)
}

// ============================================
// Usage Examples
// ============================================

func ExampleAuthentication() {
	fmt.Println("🔐 GraphQL Authentication Examples")
	fmt.Println("==================================")

	demoBasicAuthFlow()
	demoAuthenticatedQuery()
	demoAuthenticatedMutation()
	demoRoleBasedAuth()
	demoResourceAuth()
	demoAuthMiddleware()
	demoGraphQLAuthMiddleware()
	demoTokenHandling()
	demoRefreshToken()
	demoOAuth2()

	fmt.Println("\n✅ Authentication examples completed")
}

/*
Authentication Patterns:

1. JWT Token:
   - Login → Generate JWT
   - Include token in Authorization header
   - Validate on each request
   - Refresh token before expiry

2. OAuth2:
   - Redirect to provider
   - Get authorization code
   - Exchange for access token
   - Create user account

3. API Key:
   - Simple key in header
   - Less secure, for service accounts
   - No user context

Authorization Patterns:

1. Role-Based:
   - Admin, User, Guest roles
   - Check role in resolver
   - Define permissions per role

2. Resource-Based:
   - Owner can edit own resources
   - Admin can edit anything
   - Check resource ownership

3. Attribute-Based:
   - Fine-grained permissions
   - User attributes + resource attributes
   - Complex policy evaluation

Security Best Practices:

1. Always use HTTPS in production
2. Validate tokens on every request
3. Use short-lived access tokens (15min)
4. Implement refresh token rotation
5. Store tokens securely (httpOnly cookies)
6. Implement rate limiting
7. Log authentication events
8. Monitor for suspicious activity

Error Handling:

1. 401 Unauthorized: Missing/invalid token
2. 403 Forbidden: Valid token, insufficient permissions
3. Clear error messages for users
4. Generic messages for security issues
5. Log detailed errors server-side
*/
