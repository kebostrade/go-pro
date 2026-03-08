// Gin Basics - Fundamental concepts of the Gin web framework
//
// This example covers:
// - Router setup and configuration
// - Basic routing (GET, POST, PUT, DELETE, PATCH)

//go:build ignore
// - Path parameters and query strings
// - Route groups
// - Custom middleware
// - Error handling
// - JSON responses
//
// Run it: go run examples/gin_basics.go
// Test: curl http://localhost:8080

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// User represents a simple user model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

// Custom middleware - Logger
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), c.Request.Method, path)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		fmt.Printf("[%s] %s %s - Status: %d - Latency: %v\n",
			time.Now().Format("15:04:05"), c.Request.Method, path, statusCode, latency)
	}
}

// Custom middleware - Authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// In production, validate the token properly
		if token != "Bearer valid-token" {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", 123)
		c.Next()
	}
}

// Custom middleware - Rate limiter (simple version)
func RateLimitMiddleware() gin.HandlerFunc {
	// In production, use a proper rate limiting library
	requestCount := 0
	maxRequests := 5

	return func(c *gin.Context) {
		requestCount++
		if requestCount > maxRequests {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	// Create router with default middleware (Logger and Recovery)
	router := gin.Default()

	// Set Gin mode (Debug, Release, Test)
	// gin.SetMode(gin.ReleaseMode)

	// Use custom middleware
	router.Use(LoggerMiddleware())

	// Global middleware applied to all routes
	router.Use(RateLimitMiddleware())

	// Basic routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Gin Basics!",
			"status":  "running",
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Query parameters
	// Example: GET /search?q=golang&page=1
	router.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		c.JSON(200, gin.H{
			"query": query,
			"page":  page,
			"limit": limit,
		})
	})

	// Path parameters
	// Example: GET /users/123
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Find user
		for _, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				c.JSON(200, gin.H{
					"user": user,
				})
				return
			}
		}

		c.JSON(404, gin.H{"error": "User not found"})
	})

	// Multiple path parameters
	// Example: GET /users/123/posts/456
	router.GET("/users/:id/posts/:post_id", func(c *gin.Context) {
		userID := c.Param("id")
		postID := c.Param("post_id")

		c.JSON(200, gin.H{
			"user_id":  userID,
			"post_id":  postID,
			"endpoint": "user post",
		})
	})

	// Wildcard route
	// Example: GET /files/images/photo.jpg
	router.GET("/files/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.JSON(200, gin.H{
			"filepath": filepath,
			"type":     "file",
		})
	})

	// POST request with JSON body
	// Example: POST /users -d '{"name":"John","email":"john@example.com"}'
	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user.ID = len(users) + 1
		users = append(users, user)

		c.JSON(201, gin.H{
			"message": "User created successfully",
			"user":    user,
		})
	})

	// PUT request for updates
	router.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updateUser User

		if err := c.ShouldBindJSON(&updateUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Find and update user
		for i, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				updateUser.ID = user.ID
				users[i] = updateUser
				c.JSON(200, gin.H{
					"message": "User updated successfully",
					"user":    updateUser,
				})
				return
			}
		}

		c.JSON(404, gin.H{"error": "User not found"})
	})

	// DELETE request
	router.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		for i, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				users = append(users[:i], users[i+1:]...)
				c.JSON(200, gin.H{
					"message": "User deleted successfully",
				})
				return
			}
		}

		c.JSON(404, gin.H{"error": "User not found"})
	})

	// PATCH request for partial updates
	router.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updates map[string]interface{}

		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Find and partially update user
		for i, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				if name, ok := updates["name"].(string); ok {
					users[i].Name = name
				}
				if email, ok := updates["email"].(string); ok {
					users[i].Email = email
				}
				c.JSON(200, gin.H{
					"message": "User patched successfully",
					"user":    users[i],
				})
				return
			}
		}

		c.JSON(404, gin.H{"error": "User not found"})
	})

	// Route groups - API v1
	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"version": "v1",
				"users":   users,
			})
		})
		v1.GET("/users/:id", func(c *gin.Context) {
			id := c.Param("id")
			for _, user := range users {
				if fmt.Sprintf("%d", user.ID) == id {
					c.JSON(200, gin.H{
						"version": "v1",
						"user":    user,
					})
					return
				}
			}
			c.JSON(404, gin.H{"error": "User not found"})
		})
	}

	// Route groups - API v2 with authentication
	v2 := router.Group("/api/v2")
	v2.Use(AuthMiddleware())
	{
		v2.GET("/users", func(c *gin.Context) {
			// Get user info from context
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{
				"version":  "v2",
				"user_id":  userID,
				"users":    users,
			})
		})

		v2.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{
				"message": "Access granted to protected resource",
				"user_id": userID,
			})
		})
	}

	// Admin routes with multiple middleware
	admin := router.Group("/admin")
	admin.Use(AuthMiddleware())
	admin.Use(func(c *gin.Context) {
		// Additional admin check
		userID, _ := c.Get("user_id")
		if userID != 123 {
			c.JSON(403, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	})
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to admin dashboard",
				"stats": gin.H{
					"total_users": len(users),
					"uptime":      "24 hours",
				},
			})
		})
	}

	// Error handling routes
	router.GET("/error", func(c *gin.Context) {
		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Something went wrong",
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not found",
			"path":    c.Request.URL.Path,
			"message": "The requested resource does not exist",
		})
	})

	// 405 method not allowed
	router.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error":   "Method not allowed",
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
		})
	})

	return router
}

func main() {
	router := setupRouter()

	// Run the server
	fmt.Println("🚀 Gin Basics server starting on :8080")
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  GET    /")
	fmt.Println("  GET    /health")
	fmt.Println("  GET    /search?q=golang&page=1")
	fmt.Println("  GET    /users/:id")
	fmt.Println("  GET    /users/:id/posts/:post_id")
	fmt.Println("  GET    /files/*filepath")
	fmt.Println("  POST   /users")
	fmt.Println("  PUT    /users/:id")
	fmt.Println("  PATCH  /users/:id")
	fmt.Println("  DELETE /users/:id")
	fmt.Println("  GET    /api/v1/users")
	fmt.Println("  GET    /api/v2/users (requires auth)")
	fmt.Println("  GET    /admin/dashboard (requires admin)")
	fmt.Println("\nExample commands:")
	fmt.Println("  curl http://localhost:8080")
	fmt.Println("  curl http://localhost:8080/users/1")
	fmt.Println("  curl http://localhost:8080/search?q=golang")
	fmt.Println("  curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{\"name\":\"Bob\",\"email\":\"bob@example.com\"}'")
	fmt.Println("  curl http://localhost:8080/api/v2/users -H 'Authorization: Bearer valid-token'")
	fmt.Println()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
