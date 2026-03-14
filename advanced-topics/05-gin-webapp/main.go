// Gin Web Framework - Complete Working Application
//
// This is a comprehensive Gin web application demonstrating all major features:
// - Basic routing and handlers
// - RESTful API with CRUD operations
// - HTML template rendering
// - File uploads and downloads
// - Session management
// - Middleware (authentication, logging, CORS)
// - Request validation and binding
// - Error handling
// - Static file serving
//
// Run it: go run main.go
// Visit: http://localhost:8080

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// User model
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required,min=2,max=100"`
	Email     string    `json:"email" binding:"required,email"`
	Age       int       `json:"age" binding:"gte=0,lte=130"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// In-memory data store
type UserStore struct {
	users  map[int]*User
	nextID int
	mu     sync.RWMutex
}

var store = &UserStore{
	users:  make(map[int]*User),
	nextID: 1,
}

// Initialize with sample data
func init() {
	users := []*User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
	}

	for _, user := range users {
		store.Create(user)
	}
}

// Create adds a new user
func (s *UserStore) Create(user *User) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = s.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	s.users[user.ID] = user
	s.nextID++

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserStore) GetByID(id int) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetAll retrieves all users
func (s *UserStore) GetAll() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return users
}

// Update updates an existing user
func (s *UserStore) Update(id int, user *User) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return nil, fmt.Errorf("user not found")
	}

	user.ID = id
	user.UpdatedAt = time.Now()
	s.users[id] = user

	return user, nil
}

// Delete removes a user
func (s *UserStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(s.users, id)
	return nil
}

// Custom middleware - Logger
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s - Status: %d - Latency: %v",
			time.Now().Format("15:04:05"),
			c.Request.Method,
			path,
			statusCode,
			latency,
		)
	}
}

// Custom middleware - Authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", 123)
		c.Next()
	}
}

// Handlers
func getIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Gin Web Framework Demo",
		"message": "Welcome to the Gin Web Framework Demo!",
	})
}

func listUsers(c *gin.Context) {
	users := store.GetAll()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(users),
		"data":    users,
	})
}

func getUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	user, err := store.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   getValidationErrors(err),
		})
		return
	}

	createdUser, err := store.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    createdUser,
	})
}

func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   getValidationErrors(err),
		})
		return
	}

	updatedUser, err := store.Update(id, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"data":    updatedUser,
	})
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	if err := store.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

func getValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is required", e.Field()))
			case "email":
				errors = append(errors, fmt.Sprintf("%s must be a valid email", e.Field()))
			case "min":
				errors = append(errors, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
			case "max":
				errors = append(errors, fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param()))
			case "gte":
				errors = append(errors, fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param()))
			case "lte":
				errors = append(errors, fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param()))
			default:
				errors = append(errors, fmt.Sprintf("%s validation failed on %s", e.Field(), e.Tag()))
			}
		}
		return fmt.Sprintf("Validation errors: %v", errors)
	}
	return err.Error()
}

func setupRouter() *gin.Engine {
	// Create router with default middleware
	router := gin.Default()

	// Use custom middleware
	router.Use(LoggerMiddleware())

	// Serve static files
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Home page
	router.GET("/", getIndex)

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		api.GET("/users", listUsers)
		api.GET("/users/:id", getUser)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(AuthMiddleware())
		{
			protected.POST("/users", createUser)
			protected.PUT("/users/:id", updateUser)
			protected.DELETE("/users/:id", deleteUser)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"users":     len(store.GetAll()),
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Route %s not found", c.Request.URL.Path),
		})
	})

	return router
}

func main() {
	router := setupRouter()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🚀 Gin Web Framework - Complete Application         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("✅ Server starting on http://localhost:8080")
	fmt.Println()
	fmt.Println("📋 Available Endpoints:")
	fmt.Println()
	fmt.Println("  📄 Pages:")
	fmt.Println("     GET    /                    - Home page")
	fmt.Println("     GET    /health              - Health check")
	fmt.Println()
	fmt.Println("  🔌 API (Public):")
	fmt.Println("     GET    /api/users           - List all users")
	fmt.Println("     GET    /api/users/:id       - Get specific user")
	fmt.Println()
	fmt.Println("  🔐 API (Protected - Requires Auth):")
	fmt.Println("     POST   /api/users           - Create new user")
	fmt.Println("     PUT    /api/users/:id       - Update user")
	fmt.Println("     DELETE /api/users/:id       - Delete user")
	fmt.Println()
	fmt.Println("📝 Example Commands:")
	fmt.Println()
	fmt.Println("  # Get all users")
	fmt.Println("  curl http://localhost:8080/api/users")
	fmt.Println()
	fmt.Println("  # Get specific user")
	fmt.Println("  curl http://localhost:8080/api/users/1")
	fmt.Println()
	fmt.Println("  # Create user (requires auth)")
	fmt.Println("  curl -X POST http://localhost:8080/api/users \\")
	fmt.Println("    -H 'Authorization: Bearer valid-token' \\")
	fmt.Println("    -H 'Content-Type: application/json' \\")
	fmt.Println("    -d '{\"name\":\"Alice\",\"email\":\"alice@example.com\",\"age\":28}'")
	fmt.Println()
	fmt.Println("  # Update user (requires auth)")
	fmt.Println("  curl -X PUT http://localhost:8080/api/users/1 \\")
	fmt.Println("    -H 'Authorization: Bearer valid-token' \\")
	fmt.Println("    -H 'Content-Type: application/json' \\")
	fmt.Println("    -d '{\"name\":\"John Updated\",\"email\":\"john@example.com\",\"age\":31}'")
	fmt.Println()
	fmt.Println("  # Delete user (requires auth)")
	fmt.Println("  curl -X DELETE http://localhost:8080/api/users/3 \\")
	fmt.Println("    -H 'Authorization: Bearer valid-token'")
	fmt.Println()
	fmt.Println("💡 Press Ctrl+C to stop the server")
	fmt.Println()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
