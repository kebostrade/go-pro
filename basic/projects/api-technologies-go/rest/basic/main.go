package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// User represents a user in our system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// UserStore is an in-memory user storage
type UserStore struct {
	mu      sync.RWMutex
	users   map[int]*User
	nextID  int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserStore) Create(username, email string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &User{
		ID:        s.nextID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserStore) GetAll() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *UserStore) GetByID(id int) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

func (s *UserStore) Update(id int, username, email string) (*User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, false
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	return user, true
}

func (s *UserStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.users[id]
	if exists {
		delete(s.users, id)
	}
	return exists
}

// API handlers
type API struct {
	store *UserStore
}

func NewAPI() *API {
	return &API{
		store: NewUserStore(),
	}
}

// Middleware: Logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed in %v", time.Since(start))
	}
}

// Middleware: CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Middleware: Content-Type JSON
func jsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// Chain middlewares
func chain(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

// Helper functions
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	})
}

func sendSuccess(w http.ResponseWriter, status int, data interface{}, message string) {
	sendJSON(w, status, SuccessResponse{
		Data:    data,
		Message: message,
	})
}

// Handler: GET /users - List all users
func (api *API) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users := api.store.GetAll()
	sendSuccess(w, http.StatusOK, users, "")
}

// Handler: GET /users/{id} - Get user by ID
func (api *API) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, exists := api.store.GetByID(id)
	if !exists {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	sendSuccess(w, http.StatusOK, user, "")
}

// Handler: POST /users - Create new user
func (api *API) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validation
	if req.Username == "" || req.Email == "" {
		sendError(w, http.StatusBadRequest, "Username and email are required")
		return
	}

	user := api.store.Create(req.Username, req.Email)
	sendSuccess(w, http.StatusCreated, user, "User created successfully")
}

// Handler: PUT /users/{id} - Update user
func (api *API) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, exists := api.store.Update(id, req.Username, req.Email)
	if !exists {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	sendSuccess(w, http.StatusOK, user, "User updated successfully")
}

// Handler: DELETE /users/{id} - Delete user
func (api *API) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if !api.store.Delete(id) {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	sendSuccess(w, http.StatusOK, nil, "User deleted successfully")
}

// Router
func (api *API) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Apply middleware chain
	mux.HandleFunc("/users", chain(api.handleUsersRouter, loggingMiddleware, corsMiddleware, jsonMiddleware))
	mux.HandleFunc("/users/", chain(api.handleUserRouter, loggingMiddleware, corsMiddleware, jsonMiddleware))

	return mux
}

func (api *API) handleUsersRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		api.handleGetUsers(w, r)
	case http.MethodPost:
		api.handleCreateUser(w, r)
	default:
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (api *API) handleUserRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		api.handleGetUser(w, r)
	case http.MethodPut:
		api.handleUpdateUser(w, r)
	case http.MethodDelete:
		api.handleDeleteUser(w, r)
	default:
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func main() {
	api := NewAPI()

	// Seed some data
	api.store.Create("alice", "alice@example.com")
	api.store.Create("bob", "bob@example.com")

	mux := api.setupRoutes()

	port := ":8080"
	fmt.Printf("🚀 Basic REST API server starting on http://localhost%s\n", port)
	fmt.Println("📚 Available endpoints:")
	fmt.Println("  GET    /users      - List all users")
	fmt.Println("  POST   /users      - Create new user")
	fmt.Println("  GET    /users/{id} - Get user by ID")
	fmt.Println("  PUT    /users/{id} - Update user")
	fmt.Println("  DELETE /users/{id} - Delete user")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

