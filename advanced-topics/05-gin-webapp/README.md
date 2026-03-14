# Gin Web Framework - Complete Guide

This guide provides comprehensive examples for building web applications using the [Gin framework](https://gin-gonic.com/), one of the most popular Go web frameworks.

## What is Gin?

Gin is a high-performance HTTP web framework written in Go. It features a Martini-like API with much better performance -- up to 40 times faster. If you need performance and good productivity, you will love Gin.

### Key Features

- **Fast**: Radix tree based routing, small memory use. No reflection. Predictable API performance.
- **Middleware Support**: An HTTP request can pass through a list of middleware.
- **CRUD Validation**: Validate incoming JSON requests easily.
- **Route Groups**: Organize your routes hierarchically.
- **JSON Validation**: Parse and validate request JSON.
- **Error Management**: Provides a convenient way to collect all errors occurred during HTTP request.
- **Built-in Rendering**: Native JSON, XML, and HTML rendering.
- **Extensible**: Easy to add new middleware and handlers.

## Installation

```bash
go get -u github.com/gin-gonic/gin
```

## Project Structure

```
05-gin-webapp/
├── README.md                  # This guide
├── main.go                    # Complete working application
└── examples/
    ├── gin_basics.go          # Basic routing, middleware, handlers
    ├── gin_rest_api.go        # REST API example
    ├── gin_templates.go       # HTML template rendering
    └── gin_advanced.go        # File uploads, sessions, binding
```

## Quick Start

```bash
# Run the complete application
go run main.go

# Run individual examples
go run examples/gin_basics.go
go run examples/gin_rest_api.go
go run examples/gin_templates.go
go run examples/gin_advanced.go
```

## Core Concepts

### 1. Router & Engine

The Gin engine is the core of your application:

```go
import "github.com/gin-gonic/gin"

func main() {
    // Create a Gin router with default middleware:
    // Logger and Recovery (crash-free)
    r := gin.Default()

    // Or create a router without middleware
    r := gin.New()

    // Define routes
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Gin!",
        })
    })

    // Run the server
    r.Run(":8080")
}
```

### 2. Routing

Gin supports various HTTP methods and parameterized routes:

```go
// Simple route
r.GET("/ping", func(c *gin.Context) {
    c.String(200, "pong")
})

// Route with parameters
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"user_id": id})
})

// Route with wildcard
r.GET("/files/*filepath", func(c *gin.Context) {
    filepath := c.Param("filepath")
    c.JSON(200, gin.H{"filepath": filepath})
})

// Query parameters
r.GET("/search", func(c *gin.Context) {
    query := c.Query("q")
    page := c.DefaultQuery("page", "1")
    c.JSON(200, gin.H{"query": query, "page": page})
})
```

### 3. Middleware

Middleware functions can modify the request or response:

```go
// Custom middleware
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("Request received")
        c.Next()
        fmt.Println("Response sent")
    }
}

// Apply middleware globally
r.Use(Logger())

// Apply middleware to specific routes
r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "authenticated"})
})
```

### 4. Route Groups

Organize related routes:

```go
// API v1 group
v1 := r.Group("/api/v1")
{
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)
}

// API v2 group
v2 := r.Group("/api/v2")
v2.Use(AuthMiddleware())
{
    v2.GET("/users", getUsersV2)
}
```

### 5. Request Binding

Bind request data to structs:

```go
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}

r.POST("/users", func(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"user": user})
})
```

### 6. Response Rendering

Various response types:

```go
// JSON response
c.JSON(200, gin.H{"message": "success"})

// String response
c.String(200, "Hello, %s", name)

// HTML response
c.HTML(200, "index.html", gin.H{"title": "Home"})

// XML response
c.XML(200, gin.H{"message": "success"})

// File download
c.File("files/file.pdf")

// Redirect
c.Redirect(302, "/new-location")
```

## Example Files

### gin_basics.go

Covers:
- Router setup and configuration
- Basic routing (GET, POST, PUT, DELETE)
- Path parameters and query strings
- Route groups
- Custom middleware
- Error handling

**Run it:**
```bash
go run examples/gin_basics.go
# Visit: http://localhost:8080
```

### gin_rest_api.go

Covers:
- Complete REST API implementation
- CRUD operations
- Request validation
- Error handling
- JSON responses
- Route organization

**Run it:**
```bash
go run examples/gin_rest_api.go
# Test with curl or Postman
```

**Example API calls:**
```bash
# Get all users
curl http://localhost:8080/api/users

# Get specific user
curl http://localhost:8080/api/users/1

# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","age":30}'

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john@example.com","age":31}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/1
```

### gin_templates.go

Covers:
- HTML template rendering
- Template inheritance
- Dynamic data binding
- Static file serving
- Multi-template support

**Run it:**
```bash
go run examples/gin_templates.go
# Visit: http://localhost:8080
```

### gin_advanced.go

Covers:
- File upload handling
- Session management
- Advanced binding and validation
- Custom validation tags
- Streaming responses
- Graceful shutdown

**Run it:**
```bash
go run examples/gin_advanced.go
# Test file upload:
curl -X POST http://localhost:8080/upload \
  -F "file=@test.txt" \
  -F "name=Test File"
```

## Common Patterns

### Authentication Middleware

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        // Validate token
        c.Next()
    }
}
```

### Error Handling

```go
func handleError(c *gin.Context, err error) {
    // Log error
    log.Error(err)

    // Return appropriate response
    c.JSON(500, gin.H{
        "error": "Internal server error",
        "message": err.Error(),
    })
}
```

### Logging Middleware

```go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()

        log.Printf("[%s] %s %s %s %d %v",
            time.Now().Format(time.RFC3339),
            clientIP,
            method,
            path,
            statusCode,
            latency,
        )
    }
}
```

### CORS Middleware

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

## Best Practices

1. **Use Route Groups**: Organize routes logically (e.g., `/api/v1`, `/api/v2`)
2. **Implement Middleware**: Use middleware for cross-cutting concerns (auth, logging, CORS)
3. **Validate Input**: Always validate incoming request data
4. **Handle Errors**: Implement proper error handling and logging
5. **Use Structs**: Define structs for request/response bodies
6. **Keep Handlers Lean**: Move business logic to service layers
7. **Use Context**: Pass request-scoped data using `c.Set()` and `c.Get()`
8. **Graceful Shutdown**: Implement proper server shutdown handling
9. **Security**: Always use HTTPS in production, implement rate limiting
10. **Testing**: Write tests for handlers and middleware

## Testing Gin Handlers

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)

func TestGetUser(t *testing.T) {
    // Set Gin to test mode
    gin.SetMode(gin.TestMode)

    // Create test router
    router := setupTestRouter()

    // Create test request
    req, _ := http.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()

    // Perform request
    router.ServeHTTP(w, req)

    // Assert response
    if w.Code != 200 {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

## Performance Tips

1. **Use gin.Default() for development**, customize for production
2. **Enable JSON validation** with binding tags
3. **Use sync.Pool** for reusable objects
4. **Disable debug mode** in production: `gin.SetMode(gin.ReleaseMode)`
5. **Implement caching** for expensive operations
6. **Use connection pooling** for database connections

## Resources

- [Official Gin Documentation](https://gin-gonic.com/docs/)
- [Gin GitHub Repository](https://github.com/gin-gonic/gin)
- [Gin Examples](https://github.com/gin-gonic/examples)

## Running the Examples

Each example is self-contained and runnable:

```bash
# Run the complete application
go run main.go

# Or run individual examples
go run examples/gin_basics.go
go run examples/gin_rest_api.go
go run examples/gin_templates.go
go run examples/gin_advanced.go
```

## What You'll Learn

After completing these examples, you'll be able to:
- Create HTTP servers with Gin
- Implement RESTful APIs
- Handle various HTTP methods and request types
- Use middleware for cross-cutting concerns
- Validate and bind request data
- Render HTML templates
- Handle file uploads
- Manage sessions
- Implement authentication and authorization
- Build production-ready web applications

## Next Steps

1. Start with `gin_basics.go` to understand fundamentals
2. Move to `gin_rest_api.go` for practical API development
3. Learn `gin_templates.go` for server-side rendering
4. Master advanced concepts with `gin_advanced.go`
5. Build your own web application using Gin

Happy coding with Gin! 🚀
