# Gin Web Framework Examples - Detailed Guide

This document provides detailed explanations for each example in the Gin Web Framework tutorial.

## Table of Contents

1. [Main Application (main.go)](#main-application)
2. [Gin Basics (gin_basics.go)](#gin-basics)
3. [REST API (gin_rest_api.go)](#rest-api)
4. [Templates (gin_templates.go)](#templates)
5. [Advanced Features (gin_advanced.go)](#advanced-features)

---

## Main Application

**File:** `main.go`

### Overview
A complete, production-ready web application demonstrating all major Gin features in a cohesive application.

### Features
- ✅ Complete RESTful API with CRUD operations
- ✅ Authentication middleware
- ✅ Request validation and error handling
- ✅ In-memory data store with thread-safe operations
- ✅ Structured JSON responses
- ✅ Health check endpoint
- ✅ Comprehensive logging

### Key Components

#### Data Store
```go
type UserStore struct {
    users  map[int]*User
    nextID int
    mu     sync.RWMutex
}
```
- Thread-safe in-memory storage using `sync.RWMutex`
- Concurrent read access with `RLock()`
- Exclusive write access with `Lock()`

#### Authentication Middleware
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token != "Bearer valid-token" {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```
- Validates Bearer token from Authorization header
- Uses `c.Abort()` to stop request processing on auth failure
- Stores user info in context with `c.Set()`

#### Handler Functions
```go
func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": getValidationErrors(err)})
        return
    }
    // Create user and return response
}
```
- Bind JSON to struct with validation
- Custom validation error messages
- Consistent response format

### Running the Application

```bash
# Start the server
go run main.go

# Test the API
curl http://localhost:8080/api/users
curl -X POST http://localhost:8080/api/users \
  -H 'Authorization: Bearer valid-token' \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","age":28}'
```

### Response Format

All responses follow this structure:

**Success:**
```json
{
  "success": true,
  "data": {...},
  "message": "Operation successful"
}
```

**Error:**
```json
{
  "success": false,
  "error": "Error message"
}
```

---

## Gin Basics

**File:** `examples/gin_basics.go`

### Overview
Comprehensive introduction to Gin framework fundamentals.

### Topics Covered

#### 1. Router Setup
```go
// Default router with Logger and Recovery middleware
router := gin.Default()

// Custom router without middleware
router := gin.New()
```

#### 2. HTTP Methods
```go
router.GET("/path", handler)
router.POST("/path", handler)
router.PUT("/path", handler)
router.DELETE("/path", handler)
router.PATCH("/path", handler)
```

#### 3. Route Parameters
```go
// Single parameter
router.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    // Use id
})

// Multiple parameters
router.GET("/users/:id/posts/:post_id", func(c *gin.Context) {
    userID := c.Param("id")
    postID := c.Param("post_id")
})

// Wildcard
router.GET("/files/*filepath", func(c *gin.Context) {
    filepath := c.Param("filepath")
})
```

#### 4. Query Parameters
```go
router.GET("/search", func(c *gin.Context) {
    query := c.Query("q")
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "10")
})
```

#### 5. Middleware
```go
// Custom middleware
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before request
        c.Next()
        // After request
    }
}

// Apply middleware
router.Use(CustomMiddleware())
```

#### 6. Route Groups
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

#### 7. Request Body Binding
```go
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
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

#### 8. Response Types
```go
// JSON
c.JSON(200, gin.H{"message": "success"})

// String
c.String(200, "Hello, %s", name)

// HTML
c.HTML(200, "index.html", gin.H{"title": "Home"})

// XML
c.XML(200, gin.H{"message": "success"})

// Redirect
c.Redirect(302, "/new-location")
```

---

## REST API

**File:** `examples/gin_rest_api.go`

### Overview
Production-ready REST API implementation following best practices.

### Architecture

#### Repository Pattern
```go
type UserStore struct {
    users  map[int]*User
    nextID int
    mu     sync.RWMutex
}

func (s *UserStore) Create(user *User) (*User, error)
func (s *UserStore) GetByID(id int) (*User, error)
func (s *UserStore) GetAll(page, pageSize int) ([]*User, int, error)
func (s *UserStore) Update(id int, updates *UpdateUserRequest) (*User, error)
func (s *UserStore) Delete(id int) error
```

#### Request/Response DTOs
```go
// Separate DTOs for create vs update
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}

type UpdateUserRequest struct {
    Name  *string `json:"name" binding:"omitempty,min=2,max=100"`
    Email *string `json:"email" binding:"omitempty,email"`
    Age   *int    `json:"age" binding:"omitempty,gte=0,lte=130"`
}
```

#### Validation
```go
type User struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}
```

**Validation Tags:**
- `required` - Field must be present
- `min` - Minimum length/ value
- `max` - Maximum length/ value
- `email` - Must be valid email format
- `gte` - Greater than or equal to
- `lte` - Less than or equal to
- `omitempty` - Skip validation if field is empty

#### Pagination
```go
type PaginationParams struct {
    Page     int `form:"page" binding:"omitempty,min=1"`
    PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type PaginationResponse struct {
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

#### Error Handling
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Details string `json:"details,omitempty"`
}

type Response struct {
    Success bool          `json:"success"`
    Data    interface{}   `json:"data,omitempty"`
    Message string        `json:"message,omitempty"`
    Error   *ErrorResponse `json:"error,omitempty"`
}
```

### API Endpoints

#### GET /api/users
Get all users with pagination.

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 10, max: 100)

**Response:**
```json
{
  "success": true,
  "data": {
    "users": [...],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

#### GET /api/users/:id
Get specific user by ID.

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
  }
}
```

#### POST /api/users
Create new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}
```

**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {...}
}
```

#### PUT /api/users/:id
Update entire user.

**Request Body:**
```json
{
  "name": "Updated Name",
  "email": "updated@example.com",
  "age": 35
}
```

#### DELETE /api/users/:id
Delete user.

**Response:**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

---

## Templates

**File:** `examples/gin_templates.go`

### Overview
HTML template rendering with dynamic data binding.

### Template Features

#### 1. Loading Templates
```go
// Load all templates from directory
router.LoadHTMLGlob("templates/*")

// Load specific templates
router.LoadHTMLFiles("templates/layout.html", "templates/home.html")
```

#### 2. Rendering Templates
```go
router.GET("/", func(c *gin.Context) {
    c.HTML(200, "index.html", gin.H{
        "title": "Home Page",
        "users": users,
    })
})
```

#### 3. Template Inheritance
```html
<!-- base.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    {{block "content" .}}{{end}}
</body>
</html>

<!-- index.html -->
{{define "content"}}
<h1>Welcome!</h1>
{{end}}
```

#### 4. Template Functions
```go
func customTemplateFuncs() template.FuncMap {
    return template.FuncMap{
        "formatDate": func(date string) string {
            // Format date
        },
        "upper": func(s string) string {
            return strings.ToUpper(s)
        },
    }
}
```

#### 5. Loops and Conditionals
```html
<!-- Loop -->
{{range .Users}}
    <div>{{ .Name }}</div>
{{end}}

<!-- Conditional -->
{{if eq .Role "Admin"}}
    <span>Admin</span>
{{else}}
    <span>User</span>
{{end}}
```

#### 6. Static Files
```go
// Serve static files from directory
router.Static("/static", "./static")

// Serve single file
router.StaticFile("/favicon.ico", "./static/favicon.ico")
```

---

## Advanced Features

**File:** `examples/gin_advanced.go`

### Overview
Advanced Gin features for production applications.

### Topics Covered

#### 1. File Uploads

**Single File Upload:**
```go
router.POST("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }

    filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
    filepath := filepath.Join(uploadDir, filename)

    if err := c.SaveUploadedFile(file, filepath); err != nil {
        c.JSON(500, gin.H{"error": "Failed to save file"})
        return
    }

    c.JSON(200, gin.H{"message": "File uploaded successfully"})
})
```

**Multiple File Upload:**
```go
router.POST("/upload/multiple", func(c *gin.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": "Failed to parse form"})
        return
    }

    files := form.File["files"]
    for _, file := range files {
        // Save each file
    }
})
```

**Test File Upload:**
```bash
# Single file
curl -X POST http://localhost:8080/upload \
  -F "file=@test.txt" \
  -F "name=Test File"

# Multiple files
curl -X POST http://localhost:8080/upload/multiple \
  -F "files=@test1.txt" \
  -F "files=@test2.txt"
```

#### 2. Session Management

**Session Manager:**
```go
type SessionManager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
}

func (sm *SessionManager) CreateSession() *Session
func (sm *SessionManager) GetSession(id string) (*Session, bool)
```

**Session Usage:**
```go
// Create session
router.POST("/session/create", func(c *gin.Context) {
    session := sessionManager.CreateSession()
    c.SetCookie("session_id", session.ID, 86400, "/", "", false, true)
    c.JSON(200, gin.H{"session_id": session.ID})
})

// Set session data
router.POST("/session/set", func(c *gin.Context) {
    sessionID, _ := c.Cookie("session_id")
    session, _ := sessionManager.GetSession(sessionID)
    session.Data["key"] = "value"
})
```

#### 3. Custom Validation

**Register Custom Validator:**
```go
func skuValid(fl validator.FieldLevel) bool {
    sku := fl.Field().String()
    // Custom validation logic
    return true
}

// Register validator
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("sku-valid", skuValid)
}
```

**Use Custom Validator:**
```go
type Product struct {
    SKU string `json:"sku" binding:"required,sku-valid"`
}
```

#### 4. Streaming Responses

**Server-Sent Events (SSE):**
```go
router.GET("/events", func(c *gin.Context) {
    c.Header("Content-Type", "text/event-stream")
    c.Stream(func(w io.Writer) bool {
        event := fmt.Sprintf("data: Event at %s\n\n", time.Now())
        fmt.Fprint(w, event)
        w.(http.Flusher).Flush()
        return true
    })
})
```

#### 5. Cookie Management

```go
// Set cookie
c.SetCookie("name", "value", 3600, "/", "", false, false)

// Get cookie
value, err := c.Cookie("name")

// Delete cookie
c.SetCookie("name", "", -1, "/", "", false, false)
```

#### 6. Graceful Shutdown

```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      router,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
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

---

## Testing the Examples

### Manual Testing

#### Using curl
```bash
# Get all users
curl http://localhost:8080/api/users

# Create user
curl -X POST http://localhost:8080/api/users \
  -H 'Authorization: Bearer valid-token' \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","age":28}'

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H 'Authorization: Bearer valid-token' \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Updated","email":"alice@example.com","age":29}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/1 \
  -H 'Authorization: Bearer valid-token'
```

#### Using the Test Script
```bash
# Make sure server is running
go run main.go

# In another terminal, run tests
./test-api.sh
```

### Automated Testing

The `test-api.sh` script provides:
- ✅ Automated endpoint testing
- ✅ Response validation
- ✅ Authentication testing
- ✅ Error handling validation
- ✅ Detailed test results

---

## Best Practices

### 1. Code Organization
```
project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── repository/
│   ├── middleware/
│   └── models/
├── pkg/
│   └── utils/
└── go.mod
```

### 2. Error Handling
```go
// Always handle errors
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
    return
}

// Use consistent error format
type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}
```

### 3. Validation
```go
// Always validate input
type User struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
}

if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### 4. Middleware
```go
// Use middleware for cross-cutting concerns
router.Use(LoggerMiddleware())
router.Use(AuthMiddleware())
router.Use(CORSMiddleware())
```

### 5. Route Groups
```go
// Organize routes logically
v1 := router.Group("/api/v1")
v2 := router.Group("/api/v2")
admin := router.Group("/admin")
admin.Use(AuthMiddleware())
```

### 6. Context Usage
```go
// Store request-scoped data
c.Set("user_id", 123)

// Retrieve data
userID, exists := c.Get("user_id")
if !exists {
    // Handle missing data
}
```

---

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

#### 2. Module Dependencies
```bash
# Update dependencies
go mod tidy

# Verify dependencies
go mod verify
```

#### 3. Template Not Found
```bash
# Ensure templates directory exists
mkdir templates

# Check template path
ls templates/
```

#### 4. File Upload Issues
```bash
# Create upload directory
mkdir uploads
chmod 755 uploads
```

---

## Additional Resources

- [Official Gin Documentation](https://gin-gonic.com/docs/)
- [Gin GitHub Repository](https://github.com/gin-gonic/gin)
- [Gin Examples](https://github.com/gin-gonic/examples)
- [Go Validator Documentation](https://github.com/go-playground/validator)

---

## Summary

This Gin Web Framework tutorial provides:

1. **Complete Application** - Production-ready example
2. **Fundamentals** - Basic routing, middleware, handlers
3. **REST API** - Full CRUD with validation
4. **Templates** - HTML rendering with dynamic data
5. **Advanced Features** - File uploads, sessions, streaming

Each example is self-contained and runnable. Start with the basics, then progress to more advanced topics as needed.
