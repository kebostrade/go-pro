package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Models
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" binding:"required,min=3,max=50"`
	Email     string    `json:"email" binding:"required,email"`
	Role      string    `json:"role" binding:"required,oneof=admin user guest"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin user guest"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Role     string `json:"role,omitempty" binding:"omitempty,oneof=admin user guest"`
	Active   *bool  `json:"active,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserStore) GetAll(role string, active *bool) []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0)
	for _, user := range s.users {
		// Filter by role if provided
		if role != "" && user.Role != role {
			continue
		}
		// Filter by active status if provided
		if active != nil && user.Active != *active {
			continue
		}
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

func (s *UserStore) Update(id int, username, email, role string, active *bool) (*User, bool) {
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
	if active != nil {
		user.Active = *active
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
	router    *gin.Engine
	store     *UserStore
	validator *validator.Validate
}

func NewServer() *Server {
	// Set Gin to release mode for production
	// gin.SetMode(gin.ReleaseMode)

	s := &Server{
		router:    gin.Default(), // Includes Logger and Recovery middleware
		store:     NewUserStore(),
		validator: validator.New(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.handleHealth)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/login", s.handleLogin)

		// Users resource
		users := v1.Group("/users")
		{
			users.GET("", s.handleListUsers)
			users.POST("", s.handleCreateUser)
			users.GET("/:id", s.handleGetUser)
			users.PUT("/:id", s.handleUpdateUser)
			users.DELETE("/:id", s.handleDeleteUser)
			users.PATCH("/:id/activate", s.handleActivateUser)
			users.PATCH("/:id/deactivate", s.handleDeactivateUser)
		}

		// Admin routes (would normally require auth middleware)
		admin := v1.Group("/admin")
		{
			admin.GET("/stats", s.handleStats)
		}
	}
}

// Handlers
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock authentication (in real app, verify password hash)
	c.JSON(http.StatusOK, gin.H{
		"token":      "mock-jwt-token",
		"expires_in": 3600,
		"user": gin.H{
			"email": req.Email,
		},
	})
}

func (s *Server) handleListUsers(c *gin.Context) {
	// Query parameters for filtering
	role := c.Query("role")
	activeStr := c.Query("active")

	var active *bool
	if activeStr != "" {
		val := activeStr == "true"
		active = &val
	}

	users := s.store.GetAll(role, active)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"count":   len(users),
	})
}

func (s *Server) handleGetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, exists := s.store.GetByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func (s *Server) handleCreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := s.store.Create(req.Username, req.Email, req.Role)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
		"message": "User created successfully",
	})
}

func (s *Server) handleUpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := s.store.Update(id, req.Username, req.Email, req.Role, req.Active)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "User updated successfully",
	})
}

func (s *Server) handleDeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if !s.store.Delete(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

func (s *Server) handleActivateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	active := true
	user, exists := s.store.Update(id, "", "", "", &active)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "User activated successfully",
	})
}

func (s *Server) handleDeactivateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	active := false
	user, exists := s.store.Update(id, "", "", "", &active)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "User deactivated successfully",
	})
}

func (s *Server) handleStats(c *gin.Context) {
	allUsers := s.store.GetAll("", nil)

	stats := gin.H{
		"total_users":   len(allUsers),
		"active_users":  0,
		"admin_users":   0,
		"regular_users": 0,
		"guest_users":   0,
	}

	for _, user := range allUsers {
		if user.Active {
			stats["active_users"] = stats["active_users"].(int) + 1
		}
		switch user.Role {
		case "admin":
			stats["admin_users"] = stats["admin_users"].(int) + 1
		case "user":
			stats["regular_users"] = stats["regular_users"].(int) + 1
		case "guest":
			stats["guest_users"] = stats["guest_users"].(int) + 1
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

func main() {
	server := NewServer()

	// Seed data
	server.store.Create("alice", "alice@example.com", "admin")
	server.store.Create("bob", "bob@example.com", "user")
	server.store.Create("charlie", "charlie@example.com", "guest")

	port := ":8080"
	fmt.Printf("🚀 Gin REST API server starting on http://localhost%s\n", port)
	fmt.Println("📚 Available endpoints:")
	fmt.Println("  GET    /health                      - Health check")
	fmt.Println("  POST   /api/v1/login                - Login")
	fmt.Println("  GET    /api/v1/users                - List users (filter: ?role=admin&active=true)")
	fmt.Println("  POST   /api/v1/users                - Create user")
	fmt.Println("  GET    /api/v1/users/:id            - Get user")
	fmt.Println("  PUT    /api/v1/users/:id            - Update user")
	fmt.Println("  DELETE /api/v1/users/:id            - Delete user")
	fmt.Println("  PATCH  /api/v1/users/:id/activate   - Activate user")
	fmt.Println("  PATCH  /api/v1/users/:id/deactivate - Deactivate user")
	fmt.Println("  GET    /api/v1/admin/stats          - Get statistics")

	if err := server.router.Run(port); err != nil {
		panic(err)
	}
}

