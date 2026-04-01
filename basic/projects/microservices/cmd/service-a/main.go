// Package main provides the User Service (service-a) entry point.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// User represents a user entity.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// In-memory user store with seed data.
var users = map[string]User{
	"1": {ID: "1", Name: "Alice Johnson", Email: "alice@example.com", Age: 28},
	"2": {ID: "2", Name: "Bob Smith", Email: "bob@example.com", Age: 34},
	"3": {ID: "3", Name: "Charlie Brown", Email: "charlie@example.com", Age: 22},
}

func main() {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8001"
	}

	// Health check endpoint.
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "service-a",
			"port":    port,
		})
	})

	// GET /api/users - List all users.
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		userList := make([]User, 0, len(users))
		for _, u := range users {
			userList = append(userList, u)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": userList,
			"count": len(userList),
		})
	})

	// GET /api/users/{id} - Get user by ID.
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Path[len("/api/users/"):]
		user, ok := users[id]
		if !ok {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	addr := ":" + port
	log.Printf("🚀 User Service (service-a) starting on %s", addr)
	log.Printf("📋 Health: http://localhost%s/health", addr)
	log.Printf("📋 Users: http://localhost%s/api/users", addr)

	// Graceful shutdown with timeout.
	srv := &http.Server{
		Addr:         addr,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

// Handlers for testing.
func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	userList := make([]User, 0, len(users))
	for _, u := range users {
		userList = append(userList, u)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": userList,
		"count": len(userList),
	})
}

func HandleGetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserByID is exported for testing.
func GetUserByID(id string) (User, error) {
	user, ok := users[id]
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetAllUsers returns all users.
func GetAllUsers() []User {
	userList := make([]User, 0, len(users))
	for _, u := range users {
		userList = append(userList, u)
	}
	return userList
}

// GetUserCount returns the number of users.
func GetUserCount() int {
	return len(users)
}

// ParsePort parses and validates the port from environment.
func ParsePort() int {
	portStr := os.Getenv("SERVICE_PORT")
	if portStr == "" {
		portStr = "8001"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return 8001
	}
	return port
}
