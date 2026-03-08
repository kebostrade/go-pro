# Gin Web Framework - Quick Reference Guide

## Quick Start Commands

```bash
# Run the complete application
go run main.go

# Run individual examples
go run examples/gin_basics.go
go run examples/gin_rest_api.go
go run examples/gin_templates.go
go run examples/gin_advanced.go

# Use the quick start script
./quick-start.sh

# Test the API (after starting server)
./test-api.sh
```

## Common Patterns

### Router Setup

```go
// Default router (with Logger and Recovery)
router := gin.Default()

// Custom router
router := gin.New()
router.Use(gin.Logger(), gin.Recovery())
```

### Basic Routing

```go
router.GET("/path", handler)
router.POST("/path", handler)
router.PUT("/path", handler)
router.DELETE("/path", handler)
router.PATCH("/path", handler)
```

### Path Parameters

```go
// Single parameter: /users/:id
router.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
})

// Multiple parameters: /users/:id/posts/:post_id
router.GET("/users/:id/posts/:post_id", func(c *gin.Context) {
    userID := c.Param("id")
    postID := c.Param("post_id")
})

// Wildcard: /files/*filepath
router.GET("/files/*filepath", func(c *gin.Context) {
    filepath := c.Param("filepath")
})
```

### Query Parameters

```go
router.GET("/search", func(c *gin.Context) {
    query := c.Query("q")
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "10")
})
```

### Request Binding

```go
type User struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}

router.POST("/users", func(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // Process user
})
```

### Response Types

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

### Middleware

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
router.Use(Logger())

// Apply to specific routes
router.GET("/protected", AuthMiddleware(), handler)
```

### Route Groups

```go
// API v1 group
v1 := router.Group("/api/v1")
{
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)
}

// API v2 with authentication
v2 := router.Group("/api/v2")
v2.Use(AuthMiddleware())
{
    v2.GET("/users", getUsersV2)
}
```

## Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field must be present | `binding:"required"` |
| `min` | Minimum value/length | `binding:"min=3"` |
| `max` | Maximum value/length | `binding:"max=100"` |
| `email` | Valid email format | `binding:"email"` |
| `gte` | Greater than or equal | `binding:"gte=0"` |
| `lte` | Less than or equal | `binding:"lte=130"` |
| `oneof` | One of the listed values | `binding:"oneof=red blue green"` |
| `omitempty` | Skip validation if empty | `binding:"omitempty,min=3"` |

## Common HTTP Status Codes

| Code | Name | Usage |
|------|------|-------|
| 200 | OK | Successful request |
| 201 | Created | Resource created successfully |
| 204 | No Content | Successful request with no response body |
| 400 | Bad Request | Invalid request data |
| 401 | Unauthorized | Authentication required |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 405 | Method Not Allowed | HTTP method not supported |
| 409 | Conflict | Resource conflict |
| 500 | Internal Server Error | Server error |

## File Upload

```go
// Single file upload
router.POST("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }

    filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
    if err := c.SaveUploadedFile(file, filename); err != nil {
        c.JSON(500, gin.H{"error": "Failed to save file"})
        return
    }

    c.JSON(200, gin.H{"message": "File uploaded successfully"})
})

// Multiple file upload
router.POST("/upload/multiple", func(c *gin.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": "Failed to parse form"})
        return
    }

    files := form.File["files"]
    for _, file := range files {
        c.SaveUploadedFile(file, file.Filename)
    }

    c.JSON(200, gin.H{"message": "Files uploaded successfully"})
})
```

## Cookies

```go
// Set cookie
c.SetCookie(
    "name",      // cookie name
    "value",     // cookie value
    3600,        // max age in seconds
    "/",         // path
    "",          // domain
    false,       // secure
    true,        // http only
)

// Get cookie
value, err := c.Cookie("name")

// Delete cookie
c.SetCookie("name", "", -1, "/", "", false, false)
```

## HTML Templates

```go
// Load templates
router.LoadHTMLGlob("templates/*")

// Render template
c.HTML(200, "index.html", gin.H{
    "title": "Home Page",
    "users": users,
})

// Serve static files
router.Static("/static", "./static")
router.StaticFile("/favicon.ico", "./static/favicon.ico")
```

## Error Handling

```go
// Panic recovery (built-in)
router.Use(gin.Recovery())

// Custom error handling
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            c.JSON(500, gin.H{"error": err.Error()})
        }
    }
}

// 404 handler
router.NoRoute(func(c *gin.Context) {
    c.JSON(404, gin.H{"error": "Not found"})
})
```

## Context Methods

```go
// Set value in context
c.Set("user_id", 123)

// Get value from context
userID, exists := c.Get("user_id")

// Get with type assertion
if userID, ok := c.Get("user_id"); ok {
    id := userID.(int)
}

// Get header
token := c.GetHeader("Authorization")

// Set header
c.Header("X-Custom-Header", "value")
```

## Testing Endpoints

```bash
# GET request
curl http://localhost:8080/api/users

# POST request
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"John","email":"john@example.com"}'

# PUT request
curl -X PUT http://localhost:8080/api/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"John Updated"}'

# DELETE request
curl -X DELETE http://localhost:8080/api/users/1

# File upload
curl -X POST http://localhost:8080/upload \
  -F 'file=@test.txt' \
  -F 'name=Test File'

# With authentication
curl -H 'Authorization: Bearer token' \
  http://localhost:8080/api/protected
```

## Best Practices

### 1. Use Route Groups
```go
api := router.Group("/api")
{
    v1 := api.Group("/v1")
    v2 := api.Group("/v2")
}
```

### 2. Validate Input
```go
if err := c.ShouldBindJSON(&data); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### 3. Use Middleware for Cross-Cutting Concerns
```go
router.Use(Logger())
router.Use(AuthMiddleware())
router.Use(CORSMiddleware())
```

### 4. Consistent Response Format
```go
c.JSON(200, gin.H{
    "success": true,
    "data": data,
    "message": "Operation successful",
})
```

### 5. Handle Errors Properly
```go
if err != nil {
    log.Error(err)
    c.JSON(500, gin.H{"error": "Internal server error"})
    return
}
```

### 6. Use Structs for Request/Response
```go
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}
```

### 7. Graceful Shutdown
```go
srv := &http.Server{
    Addr:    ":8080",
    Handler: router,
}

go srv.ListenAndServe()

// Wait for interrupt signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

## Common Middleware Patterns

### CORS
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### Authentication
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Authorization required"})
            c.Abort()
            return
        }

        // Validate token
        // ...

        c.Set("user_id", userID)
        c.Next()
    }
}
```

### Logging
```go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        latency := time.Since(start)
        log.Printf("[%s] %s %s - %d - %v",
            time.Now().Format("15:04:05"),
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            latency,
        )
    }
}
```

## Project Structure

```
project/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handler/             # HTTP handlers
│   ├── repository/          # Data access layer
│   ├── middleware/          # Custom middleware
│   └── models/              # Data models
├── pkg/                     # Public packages
├── templates/               # HTML templates
├── static/                  # Static files
├── go.mod
└── go.sum
```

## Resources

- [Official Documentation](https://gin-gonic.com/docs/)
- [GitHub Repository](https://github.com/gin-gonic/gin)
- [Examples](https://github.com/gin-gonic/examples)
- [Go Validator](https://github.com/go-playground/validator)

## Tips

1. Use `gin.Default()` for development
2. Set `gin.SetMode(gin.ReleaseMode)` for production
3. Always validate input data
4. Use meaningful route names
5. Implement proper error handling
6. Add logging middleware
7. Use HTTPS in production
8. Implement rate limiting
9. Keep handlers lean
10. Write tests for your handlers
