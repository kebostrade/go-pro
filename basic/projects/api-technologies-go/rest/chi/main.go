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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
)

// Models
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Role      string    `json:"role" validate:"required,oneof=admin user guest"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof=admin user guest"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=admin user guest"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// UserStore - In-memory storage
type UserStore struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserStore) Create(username, email, role string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	user := &User{
		ID:        s.nextID,
		Username:  username,
		Email:     email,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
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

func (s *UserStore) Update(id int, username, email, role string) (*User, bool) {
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
	if role != "" {
		user.Role = role
	}
	user.UpdatedAt = time.Now()
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

// API Server
type Server struct {
	router    *chi.Mux
	store     *UserStore
	validator *validator.Validate
}

func NewServer() *Server {
	s := &Server{
		router:    chi.NewRouter(),
		store:     NewUserStore(),
		validator: validator.New(),
	}
	s.setupMiddleware()
	s.setupRoutes()
	return s
}

func (s *Server) setupMiddleware() {
	// Built-in Chi middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Custom middleware: JSON Content-Type
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/health", s.handleHealth)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Users resource
		r.Route("/users", func(r chi.Router) {
			r.Get("/", s.handleListUsers)
			r.Post("/", s.handleCreateUser)

			// User by ID
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(s.userCtx) // Load user into context
				r.Get("/", s.handleGetUser)
				r.Put("/", s.handleUpdateUser)
				r.Delete("/", s.handleDeleteUser)
			})
		})
	})
}

// Middleware: Load user into context
func (s *Server) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		id, err := strconv.Atoi(userID)
		if err != nil {
			s.sendError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		user, exists := s.store.GetByID(id)
		if !exists {
			s.sendError(w, http.StatusNotFound, "User not found")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handlers
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.sendSuccess(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}, "")
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users := s.store.GetAll()
	s.sendSuccess(w, http.StatusOK, users, "")
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	s.sendSuccess(w, http.StatusOK, user, "")
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := s.validator.Struct(req); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := s.store.Create(req.Username, req.Email, req.Role)
	s.sendSuccess(w, http.StatusCreated, user, "User created successfully")
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := s.validator.Struct(req); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedUser, _ := s.store.Update(user.ID, req.Username, req.Email, req.Role)
	s.sendSuccess(w, http.StatusOK, updatedUser, "User updated successfully")
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	s.store.Delete(user.ID)
	s.sendSuccess(w, http.StatusOK, nil, "User deleted successfully")
}

// Helper methods
func (s *Server) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) sendError(w http.ResponseWriter, status int, message string) {
	s.sendJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

func (s *Server) sendSuccess(w http.ResponseWriter, status int, data interface{}, message string) {
	s.sendJSON(w, status, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func main() {
	server := NewServer()

	// Seed data
	server.store.Create("alice", "alice@example.com", "admin")
	server.store.Create("bob", "bob@example.com", "user")
	server.store.Create("charlie", "charlie@example.com", "guest")

	port := ":8080"
	fmt.Printf("🚀 Chi REST API server starting on http://localhost%s\n", port)
	fmt.Println("📚 Available endpoints:")
	fmt.Println("  GET    /health              - Health check")
	fmt.Println("  GET    /api/v1/users        - List all users")
	fmt.Println("  POST   /api/v1/users        - Create new user")
	fmt.Println("  GET    /api/v1/users/{id}   - Get user by ID")
	fmt.Println("  PUT    /api/v1/users/{id}   - Update user")
	fmt.Println("  DELETE /api/v1/users/{id}   - Delete user")

	if err := http.ListenAndServe(port, server.router); err != nil {
		log.Fatal(err)
	}
}

