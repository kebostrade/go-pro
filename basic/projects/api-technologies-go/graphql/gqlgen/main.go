package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Models
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	Posts     []*Post   `json:"posts"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	AuthorID  string    `json:"authorId"`
	Author    *User     `json:"author"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Store
type Store struct {
	mu         sync.RWMutex
	users      map[string]*User
	posts      map[string]*Post
	nextUserID int
	nextPostID int
}

func NewStore() *Store {
	return &Store{
		users:      make(map[string]*User),
		posts:      make(map[string]*Post),
		nextUserID: 1,
		nextPostID: 1,
	}
}

func (s *Store) CreateUser(username, email, role string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := strconv.Itoa(s.nextUserID)
	s.nextUserID++

	user := &User{
		ID:        id,
		Username:  username,
		Email:     email,
		Role:      role,
		Active:    true,
		Posts:     []*Post{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.users[id] = user
	return user
}

func (s *Store) GetUser(id string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

func (s *Store) GetAllUsers() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *Store) CreatePost(title, content string, published bool, authorID string, tags []string) *Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := strconv.Itoa(s.nextPostID)
	s.nextPostID++

	post := &Post{
		ID:        id,
		Title:     title,
		Content:   content,
		Published: published,
		AuthorID:  authorID,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.posts[id] = post

	// Add to user's posts
	if user, exists := s.users[authorID]; exists {
		user.Posts = append(user.Posts, post)
	}

	return post
}

func (s *Store) GetPost(id string) *Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.posts[id]
}

func (s *Store) GetAllPosts() []*Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts
}

// GraphQL Handler
type GraphQLHandler struct {
	store *Store
}

func NewGraphQLHandler() *GraphQLHandler {
	return &GraphQLHandler{
		store: NewStore(),
	}
}

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GraphQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body")
		return
	}

	// Simple query parsing (in real app, use a proper GraphQL library)
	result := h.executeQuery(r.Context(), req.Query, req.Variables)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *GraphQLHandler) executeQuery(ctx context.Context, query string, variables map[string]interface{}) GraphQLResponse {
	// Simplified query execution for tutorial
	// In production, use gqlgen or graphql-go

	// Example: Handle basic queries
	if contains(query, "users") && !contains(query, "createUser") {
		users := h.store.GetAllUsers()
		return GraphQLResponse{Data: map[string]interface{}{"users": users}}
	}

	if contains(query, "user(") {
		// Extract ID from query (simplified)
		id := extractID(query)
		user := h.store.GetUser(id)
		if user == nil {
			return GraphQLResponse{Errors: []string{"User not found"}}
		}
		return GraphQLResponse{Data: map[string]interface{}{"user": user}}
	}

	if contains(query, "posts") {
		posts := h.store.GetAllPosts()
		return GraphQLResponse{Data: map[string]interface{}{"posts": posts}}
	}

	if contains(query, "createUser") {
		username := getVariable(variables, "username", "").(string)
		email := getVariable(variables, "email", "").(string)
		role := getVariable(variables, "role", "USER").(string)

		user := h.store.CreateUser(username, email, role)
		return GraphQLResponse{Data: map[string]interface{}{"createUser": user}}
	}

	if contains(query, "createPost") {
		title := getVariable(variables, "title", "").(string)
		content := getVariable(variables, "content", "").(string)
		published := getVariable(variables, "published", false).(bool)
		authorID := getVariable(variables, "authorId", "").(string)
		tags := []string{}

		post := h.store.CreatePost(title, content, published, authorID, tags)
		return GraphQLResponse{Data: map[string]interface{}{"createPost": post}}
	}

	return GraphQLResponse{Errors: []string{"Query not supported in this simplified example"}}
}

func (h *GraphQLHandler) sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(GraphQLResponse{
		Errors: []string{message},
	})
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != "" &&
		(s == substr || len(s) >= len(substr))
}

func extractID(query string) string {
	// Simplified ID extraction
	return "1"
}

func getVariable(vars map[string]interface{}, key string, defaultVal interface{}) interface{} {
	if val, ok := vars[key]; ok {
		return val
	}
	return defaultVal
}

func main() {
	handler := NewGraphQLHandler()

	// Seed data
	user1 := handler.store.CreateUser("alice", "alice@example.com", "ADMIN")
	user2 := handler.store.CreateUser("bob", "bob@example.com", "USER")
	handler.store.CreatePost("First Post", "Hello GraphQL!", true, user1.ID, []string{"tutorial", "graphql"})
	handler.store.CreatePost("Second Post", "Learning Go", false, user2.ID, []string{"go", "programming"})

	// Serve GraphQL endpoint
	http.Handle("/graphql", handler)

	// Serve simple playground HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, playgroundHTML)
	})

	port := ":8081"
	fmt.Printf("🚀 GraphQL server starting on http://localhost%s\n", port)
	fmt.Printf("📊 GraphQL Playground: http://localhost%s\n", port)
	fmt.Printf("📡 GraphQL Endpoint: http://localhost%s/graphql\n", port)
	fmt.Println("\n📚 Example queries (use in playground):")
	fmt.Println(exampleQueries)

	log.Fatal(http.ListenAndServe(port, nil))
}

const playgroundHTML = `<!DOCTYPE html>
<html>
<head>
    <title>GraphQL Playground</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        textarea { width: 100%; height: 200px; font-family: monospace; }
        button { padding: 10px 20px; background: #0066cc; color: white; border: none; cursor: pointer; }
        button:hover { background: #0052a3; }
        pre { background: #f4f4f4; padding: 10px; border-radius: 4px; }
    </style>
</head>
<body>
    <h1>🚀 GraphQL Playground</h1>
    <h3>Query:</h3>
    <textarea id="query">query GetUsers {
  users {
    id
    username
    email
    role
  }
}</textarea>
    <br><br>
    <button onclick="executeQuery()">Execute Query</button>
    <h3>Response:</h3>
    <pre id="response"></pre>

    <script>
        async function executeQuery() {
            const query = document.getElementById('query').value;
            const response = await fetch('/graphql', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ query })
            });
            const data = await response.json();
            document.getElementById('response').textContent = JSON.stringify(data, null, 2);
        }
    </script>
</body>
</html>`

const exampleQueries = `
# Get all users
query GetUsers {
  users {
    id
    username
    email
    role
  }
}

# Get all posts
query GetPosts {
  posts {
    id
    title
    content
    published
  }
}

# Create user (use variables)
mutation CreateUser {
  createUser(input: {
    username: "john"
    email: "john@example.com"
    role: "USER"
  }) {
    id
    username
    email
  }
}
`
