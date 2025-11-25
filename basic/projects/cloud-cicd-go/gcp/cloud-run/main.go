package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Response represents API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Server represents the HTTP server
type Server struct {
	router *mux.Router
	users  map[string]*User
}

// NewServer creates a new server instance
func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
		users:  make(map[string]*User),
	}
	s.routes()
	s.seedData()
	return s
}

// routes sets up all HTTP routes
func (s *Server) routes() {
	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
	
	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", s.handleGetUsers).Methods("GET")
	api.HandleFunc("/users", s.handleCreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", s.handleGetUser).Methods("GET")
	api.HandleFunc("/users/{id}", s.handleUpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", s.handleDeleteUser).Methods("DELETE")
	
	// Info endpoint
	s.router.HandleFunc("/", s.handleInfo).Methods("GET")
	
	// Middleware
	s.router.Use(loggingMiddleware)
	s.router.Use(corsMiddleware)
}

// seedData adds sample data
func (s *Server) seedData() {
	s.users["1"] = &User{
		ID:        "1",
		Name:      "Alice Johnson",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
	}
	s.users["2"] = &User{
		ID:        "2",
		Name:      "Bob Smith",
		Email:     "bob@example.com",
		CreatedAt: time.Now(),
	}
}

// Handlers

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"status":    "UP",
			"timestamp": time.Now(),
			"service":   "cloud-run-demo",
		},
	})
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Cloud Run Demo API",
		Data: map[string]interface{}{
			"version":     "1.0.0",
			"environment": os.Getenv("ENV"),
			"endpoints": []string{
				"GET /health",
				"GET /api/v1/users",
				"POST /api/v1/users",
				"GET /api/v1/users/{id}",
				"PUT /api/v1/users/{id}",
				"DELETE /api/v1/users/{id}",
			},
		},
	})
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	user, exists := s.users[id]
	if !exists {
		respondJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}
	
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}
	
	user.ID = fmt.Sprintf("%d", len(s.users)+1)
	user.CreatedAt = time.Now()
	s.users[user.ID] = &user
	
	respondJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	user, exists := s.users[id]
	if !exists {
		respondJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}
	
	var updates User
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}
	
	if updates.Name != "" {
		user.Name = updates.Name
	}
	if updates.Email != "" {
		user.Email = updates.Email
	}
	
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	if _, exists := s.users[id]; !exists {
		respondJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}
	
	delete(s.users, id)
	
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// Middleware

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
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

// Helpers

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server := NewServer()
	
	log.Printf("🚀 Cloud Run server starting on port %s", port)
	log.Printf("📊 Environment: %s", os.Getenv("ENV"))
	log.Printf("🌐 Access the API at http://localhost:%s", port)
	
	if err := http.ListenAndServe(":"+port, server.router); err != nil {
		log.Fatal(err)
	}
}

