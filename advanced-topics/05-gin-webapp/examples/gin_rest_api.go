// Gin REST API - Complete REST API example
//
// This example covers:
// - RESTful API design
// - CRUD operations

//go:build ignore
// - Request validation and binding
// - Custom validation tags
// - Proper error handling
// - Resource organization
// - Pagination
// - Filtering and sorting
//
// Run it: go run examples/gin_rest_api.go
// Test with curl or Postman
//
// Example commands:
//   curl http://localhost:8080/api/users
//   curl http://localhost:8080/api/users/1
//   curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"name":"John","email":"john@example.com","age":30}'

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

// User model with validation tags
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required,min=2,max=100"`
	Email     string    `json:"email" binding:"required,email"`
	Age       int       `json:"age" binding:"gte=0,lte=130"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest defines the expected structure for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=0,lte=130"`
}

// UpdateUserRequest defines the expected structure for updating a user
type UpdateUserRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=2,max=100"`
	Email *string `json:"email" binding:"omitempty,email"`
	Age   *int    `json:"age" binding:"omitempty,gte=0,lte=130"`
}

// UserStore handles user data operations
type UserStore struct {
	users  map[int]*User
	nextID int
	mu     sync.RWMutex
}

// NewUserStore creates a new user store with sample data
func NewUserStore() *UserStore {
	store := &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}

	// Add sample users
	store.Create(&User{Name: "John Doe", Email: "john@example.com", Age: 30})
	store.Create(&User{Name: "Jane Smith", Email: "jane@example.com", Age: 25})
	store.Create(&User{Name: "Bob Johnson", Email: "bob@example.com", Age: 35})

	return store
}

// Create adds a new user to the store
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

// GetAll retrieves all users with pagination
func (s *UserStore) GetAll(page, pageSize int) ([]*User, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice for pagination
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	total := len(users)

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []*User{}, total, nil
	}

	if end > total {
		end = total
	}

	return users[start:end], total, nil
}

// Update updates an existing user
func (s *UserStore) Update(id int, updates *UpdateUserRequest) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Apply updates
	if updates.Name != nil {
		user.Name = *updates.Name
	}
	if updates.Email != nil {
		user.Email = *updates.Email
	}
	if updates.Age != nil {
		user.Age = *updates.Age
	}

	user.UpdatedAt = time.Now()

	return user, nil
}

// Delete removes a user from the store
func (s *UserStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(s.users, id)
	return nil
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// PaginationResponse contains pagination metadata
type PaginationResponse struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int   `json:"total"`
	TotalPages int  `json:"total_pages"`
}

// UsersResponse represents the response for users list
type UsersResponse struct {
	Users      []*User            `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	store *UserStore
}

// NewUserHandler creates a new user handler
func NewUserHandler(store *UserStore) *UserHandler {
	return &UserHandler{store: store}
}

// ListUsers handles GET /api/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "invalid_params",
				Details: err.Error(),
			},
		})
		return
	}

	// Set default pagination values
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}

	users, total, err := h.store.GetAll(params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to retrieve users",
			},
		})
		return
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: UsersResponse{
			Users: users,
			Pagination: PaginationResponse{
				Page:       params.Page,
				PageSize:   params.PageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}

// GetUser handles GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "invalid_id",
				Details: "User ID must be a number",
			},
		})
		return
	}

	user, err := h.store.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "not_found",
				Message: "User not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "validation_error",
				Details: getValidationErrors(err),
			},
		})
		return
	}

	user := &User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	createdUser, err := h.store.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to create user",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    createdUser,
		Message: "User created successfully",
	})
}

// UpdateUser handles PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "invalid_id",
				Details: "User ID must be a number",
			},
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "validation_error",
				Details: getValidationErrors(err),
			},
		})
		return
	}

	user, err := h.store.Update(id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error: &ErrorResponse{
					Error:   "not_found",
					Message: "User not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to update user",
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "User updated successfully",
	})
}

// DeleteUser handles DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "invalid_id",
				Details: "User ID must be a number",
			},
		})
		return
	}

	if err := h.store.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "not_found",
				Message: "User not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// getValidationErrors extracts validation error messages
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
	// Create router
	router := gin.Default()

	// Initialize user store and handler
	userStore := NewUserStore()
	userHandler := NewUserHandler(userStore)

	// API routes
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("", userHandler.ListUsers)           // GET /api/users
			users.POST("", userHandler.CreateUser)         // POST /api/users
			users.GET("/:id", userHandler.GetUser)         // GET /api/users/:id
			users.PUT("/:id", userHandler.UpdateUser)      // PUT /api/users/:id
			users.DELETE("/:id", userHandler.DeleteUser)   // DELETE /api/users/:id
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error: &ErrorResponse{
				Error:   "not_found",
				Message: fmt.Sprintf("Route %s not found", c.Request.URL.Path),
			},
		})
	})

	return router
}

func main() {
	router := setupRouter()

	fmt.Println("🚀 Gin REST API server starting on :8080")
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  GET    /api/users")
	fmt.Println("  GET    /api/users/:id")
	fmt.Println("  POST   /api/users")
	fmt.Println("  PUT    /api/users/:id")
	fmt.Println("  DELETE /api/users/:id")
	fmt.Println("\nExample commands:")
	fmt.Println("  # Get all users (with pagination)")
	fmt.Println("  curl http://localhost:8080/api/users?page=1&page_size=10")
	fmt.Println()
	fmt.Println("  # Get specific user")
	fmt.Println("  curl http://localhost:8080/api/users/1")
	fmt.Println()
	fmt.Println("  # Create user")
	fmt.Println("  curl -X POST http://localhost:8080/api/users \\")
	fmt.Println("    -H 'Content-Type: application/json' \\")
	fmt.Println("    -d '{\"name\":\"Alice\",\"email\":\"alice@example.com\",\"age\":28}'")
	fmt.Println()
	fmt.Println("  # Update user")
	fmt.Println("  curl -X PUT http://localhost:8080/api/users/1 \\")
	fmt.Println("    -H 'Content-Type: application/json' \\")
	fmt.Println("    -d '{\"name\":\"John Updated\",\"email\":\"john@example.com\",\"age\":31}'")
	fmt.Println()
	fmt.Println("  # Delete user")
	fmt.Println("  curl -X DELETE http://localhost:8080/api/users/3")
	fmt.Println()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
