package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// DOMAIN MODELS
// ============================================================================

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// ============================================================================
// IN-MEMORY DATABASE
// ============================================================================

type Database struct {
	users map[string]*User
	mu    sync.RWMutex
	autoID int
}

func NewDatabase() *Database {
	return &Database{
		users: make(map[string]*User),
	}
}

func (db *Database) CreateUser(user *User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.autoID++
	user.ID = fmt.Sprintf("%d", db.autoID)
	user.CreatedAt = time.Now()
	db.users[user.ID] = user
	return nil
}

func (db *Database) GetUser(id string) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (db *Database) ListUsers() ([]User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]User, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, *user)
	}
	return users, nil
}

func (db *Database) UpdateUser(id string, updated *User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.users[id]; !ok {
		return fmt.Errorf("user not found")
	}

	updated.ID = id
	updated.CreatedAt = db.users[id].CreatedAt
	db.users[id] = updated
	return nil
}

func (db *Database) DeleteUser(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.users[id]; !ok {
		return fmt.Errorf("user not found")
	}

	delete(db.users, id)
	return nil
}

// ============================================================================
// HANDLERS
// ============================================================================

type Handler struct {
	db *Database
}

func NewHandler(db *Database) *Handler {
	return &Handler{db: db}
}

// CreateUser handles POST /users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Validation
	if user.Name == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		return
	}
	if user.Email == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email is required")
		return
	}

	// Create user
	if err := h.db.CreateUser(&user); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

// GetUser handles GET /users/{id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	if id == "" {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "User ID is required")
		return
	}

	user, err := h.db.GetUser(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// ListUsers handles GET /users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.ListUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list users")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": users,
		"meta": map[string]interface{}{
			"total": len(users),
		},
	})
}

// UpdateUser handles PUT /users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	if id == "" {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "User ID is required")
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Update user
	if err := h.db.UpdateUser(id, &user); err != nil {
		if err.Error() == "user not found" {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		} else {
			respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update user")
		}
		return
	}

	user.ID = id
	respondJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /users/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	if id == "" {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "User ID is required")
		return
	}

	if err := h.db.DeleteUser(id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// HELPERS
// ============================================================================

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, code, message string) {
	w.WriteHeader(status)
	var errResp ErrorResponse
	errResp.Error.Code = code
	errResp.Error.Message = message
	json.NewEncoder(w).Encode(errResp)
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	db := NewDatabase()
	handler := NewHandler(db)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/users", handler.ListUsers)
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetUser(w, r)
		case http.MethodPut:
			handler.UpdateUser(w, r)
		case http.MethodDelete:
			handler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/users/create", handler.CreateUser)
	mux.HandleFunc("/health", handler.HealthCheck)

	// Apply middleware
	handlerChain := loggingMiddleware(corsMiddleware(jsonMiddleware(mux)))

	// Start server
	addr := ":8080"
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("\nEndpoints:")
	log.Printf("  GET    /health")
	log.Printf("  GET    /users")
	log.Printf("  POST   /users/create")
	log.Printf("  GET    /users/{id}")
	log.Printf("  PUT    /users/{id}")
	log.Printf("  DELETE /users/{id}")

	if err := http.ListenAndServe(addr, handlerChain); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// ============================================================================
// EXAMPLE USAGE
// ============================================================================

/*
# Health check
curl http://localhost:8080/health

# List users
curl http://localhost:8080/users

# Create user
curl -X POST http://localhost:8080/users/create \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","age":30}'

# Get user
curl http://localhost:8080/users/1

# Update user
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com","age":25}'

# Delete user
curl -X DELETE http://localhost:8080/users/1
*/
