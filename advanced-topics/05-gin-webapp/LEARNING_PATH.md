# Gin Web Framework - Learning Path

This guide provides a structured learning path for mastering the Gin web framework.

## Prerequisites

Before starting this tutorial, you should have:

- ✅ Basic Go knowledge (variables, functions, structs, interfaces)
- ✅ Understanding of HTTP basics (methods, status codes, headers)
- ✅ Familiarity with JSON format
- ✅ Go 1.23 or higher installed

## Learning Path Overview

```
Level 1: Basics        → gin_basics.go
Level 2: REST API      → gin_rest_api.go
Level 3: Templates     → gin_templates.go
Level 4: Advanced      → gin_advanced.go
Level 5: Production    → main.go
```

---

## Level 1: Gin Basics

**File:** `examples/gin_basics.go`

### Learning Objectives

By the end of this level, you will understand:
- How to create a Gin router
- Basic routing (GET, POST, PUT, DELETE)
- Path parameters and query strings
- Middleware basics
- Request handling

### Concepts Covered

#### 1. Router Creation
```go
router := gin.Default()  // With Logger and Recovery middleware
router := gin.New()      // Without middleware
```

#### 2. Simple Routes
```go
router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Hello, Gin!"})
})
```

#### 3. Path Parameters
```go
router.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"user_id": id})
})
```

#### 4. Query Parameters
```go
router.GET("/search", func(c *gin.Context) {
    query := c.Query("q")
    page := c.DefaultQuery("page", "1")
    c.JSON(200, gin.H{"query": query, "page": page})
})
```

#### 5. Request Body
```go
router.POST("/users", func(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, gin.H{"user": user})
})
```

### Exercises

1. **Create a simple endpoint** that returns your name
2. **Add a parameterized route** to greet users by name
3. **Create a route group** for API endpoints
4. **Implement custom middleware** that logs request duration
5. **Add a wildcard route** for file paths

### Challenge
Build a simple todo API with:
- GET /todos - List all todos
- POST /todos - Create a todo
- GET /todos/:id - Get specific todo
- DELETE /todos/:id - Delete a todo

---

## Level 2: REST API

**File:** `examples/gin_rest_api.go`

### Learning Objectives

By the end of this level, you will understand:
- RESTful API design principles
- CRUD operations
- Request validation
- Error handling best practices
- Pagination
- Data repository pattern

### Concepts Covered

#### 1. Repository Pattern
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

#### 2. Request Validation
```go
type User struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}
```

#### 3. Consistent Response Format
```go
type Response struct {
    Success bool          `json:"success"`
    Data    interface{}   `json:"data,omitempty"`
    Message string        `json:"message,omitempty"`
    Error   *ErrorResponse `json:"error,omitempty"`
}
```

#### 4. Pagination
```go
type PaginationResponse struct {
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

### Exercises

1. **Create a product API** with full CRUD operations
2. **Add filtering** to your list endpoint (e.g., filter by category)
3. **Implement sorting** (e.g., sort by name, date, price)
4. **Add search functionality** with partial matching
5. **Implement PATCH** for partial updates

### Challenge
Build a blog API with:
- Posts (CRUD)
- Comments on posts
- Categories/Tags
- Search functionality
- Pagination and filtering

---

## Level 3: Templates

**File:** `examples/gin_templates.go`

### Learning Objectives

By the end of this level, you will understand:
- HTML template rendering
- Template inheritance
- Dynamic data binding
- Static file serving
- Template functions

### Concepts Covered

#### 1. Loading Templates
```go
router.LoadHTMLGlob("templates/*")
router.LoadHTMLFiles("templates/layout.html", "templates/home.html")
```

#### 2. Rendering Templates
```go
c.HTML(200, "index.html", gin.H{
    "title": "Home Page",
    "users": users,
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
```

#### 4. Template Functions
```go
func formatDate(date string) string {
    // Format date
}
```

### Exercises

1. **Create a personal website** with multiple pages
2. **Add a contact form** with form handling
3. **Implement pagination** in your template
4. **Add dynamic content** based on URL parameters
5. **Create a dashboard** with charts/tables

### Challenge
Build a blog frontend with:
- Home page with post list
- Individual post pages
- Category pages
- Search results page
- About/contact pages

---

## Level 4: Advanced Features

**File:** `examples/gin_advanced.go`

### Learning Objectives

By the end of this level, you will understand:
- File upload handling
- Session management
- Custom validators
- Streaming responses
- Cookie management
- Graceful shutdown

### Concepts Covered

#### 1. File Upload
```go
file, err := c.FormFile("file")
c.SaveUploadedFile(file, filepath)
```

#### 2. Session Management
```go
session := sessionManager.CreateSession()
c.SetCookie("session_id", session.ID, 86400, "/", "", false, true)
```

#### 3. Custom Validation
```go
v.RegisterValidation("sku-valid", skuValid)
```

#### 4. Streaming
```go
c.Stream(func(w io.Writer) bool {
    fmt.Fprint(w, "data\n")
    w.(http.Flusher).Flush()
    return true
})
```

### Exercises

1. **Build a file upload service** with progress tracking
2. **Implement user authentication** with sessions
3. **Create real-time updates** with Server-Sent Events
4. **Build a custom validator** for phone numbers
5. **Implement remember me** with cookies

### Challenge
Build a user dashboard with:
- File upload/management
- User profile with avatar upload
- Real-time notifications
- Session-based authentication
- Remember me functionality

---

## Level 5: Production Application

**File:** `main.go`

### Learning Objectives

By the end of this level, you will understand:
- Building production-ready applications
- Security best practices
- Performance optimization
- Error handling
- Logging
- Graceful shutdown

### Concepts Covered

#### 1. Application Structure
```
cmd/server/main.go
internal/handler/
internal/repository/
internal/middleware/
```

#### 2. Security
```go
c.Header("X-Frame-Options", "DENY")
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-XSS-Protection", "1; mode=block")
```

#### 3. Graceful Shutdown
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

### Exercises

1. **Add rate limiting** to your API
2. **Implement caching** for expensive operations
3. **Add metrics** and monitoring
4. **Implement proper logging** with levels
5. **Add database integration** (PostgreSQL, MongoDB)

### Challenge
Build a production-ready SaaS application with:
- User authentication (JWT)
- Rate limiting
- Database integration
- Caching layer
- Monitoring and logging
- API documentation
- Tests

---

## Project Ideas

### Beginner
1. **URL Shortener** - Create a service to shorten URLs
2. **Task Manager** - Simple todo application
3. **Guestbook** - Collect and display guest messages
4. **Contact Form** - Form with email sending
5. **Link Directory** - Organized list of links

### Intermediate
1. **Blog Platform** - Full blog with posts, comments, categories
2. **URL Shortener with Analytics** - Track clicks and statistics
3. **File Manager** - Upload, organize, and download files
4. **Wiki System** - Collaborative documentation
5. **API Gateway** - Proxy and route requests to multiple services

### Advanced
1. **E-commerce Platform** - Products, cart, checkout
2. **Chat Application** - Real-time messaging
3. **Social Network** - Users, posts, followers
4. **Project Management Tool** - Tasks, projects, teams
5. **Monitoring Dashboard** - System metrics and alerts

---

## Best Practices Checklist

### Code Organization
- [ ] Use route groups for organization
- [ ] Separate handlers from business logic
- [ ] Implement repository pattern for data access
- [ ] Use middleware for cross-cutting concerns
- [ ] Keep handlers focused and lean

### Security
- [ ] Validate all input data
- [ ] Use HTTPS in production
- [ ] Implement authentication/authorization
- [ ] Set security headers
- [ ] Sanitize user input
- [ ] Implement rate limiting

### Error Handling
- [ ] Handle all errors appropriately
- [ ] Return consistent error responses
- [ ] Log errors for debugging
- [ ] Don't expose sensitive information
- [ ] Use appropriate HTTP status codes

### Performance
- [ ] Use connection pooling
- [ ] Implement caching where appropriate
- [ ] Optimize database queries
- [ ] Use compression for large responses
- [ ] Monitor performance metrics

### Testing
- [ ] Write unit tests for handlers
- [ ] Write integration tests
- [ ] Test error cases
- [ ] Use test fixtures
- [ ] Mock external dependencies

---

## Resources

### Official Resources
- [Gin Documentation](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [Gin Examples](https://github.com/gin-gonic/examples)

### Learning Resources
- [Go Web Examples](https://gowebexamples.com/)
- [Building Web Apps with Gin](https://blog.logrocket.com/building-web-applications-with-golang-and-gin/)
- [Gin Tutorial](https://tutorialedge.net/golang/gin-framework-tutorial/)

### Related Technologies
- [Go Validator](https://github.com/go-playground/validator)
- [sqlx](https://github.com/jmoiron/sqlx) - Database library
- [testify](https://github.com/stretchr/testify) - Testing toolkit

---

## Next Steps

After completing this learning path:

1. **Build a Real Project**
   - Start with something small but useful
   - Apply what you've learned
   - Iterate and improve

2. **Learn Related Technologies**
   - Database integration (PostgreSQL, MongoDB)
   - Authentication (JWT, OAuth)
   - Testing ( testify, go-sqlmock)
   - Docker for deployment

3. **Explore Advanced Topics**
   - Microservices architecture
   - GraphQL APIs
   - WebSockets
   - gRPC

4. **Contribute to Open Source**
   - Find Gin-based projects to contribute to
   - Share your own examples
   - Help others learn

5. **Build a Portfolio**
   - Create real applications
   - Write blog posts about your learning
   - Share code on GitHub

---

## Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Find and kill process
lsof -i :8080
kill -9 <PID>
```

#### Module Dependencies
```bash
go mod tidy
go mod verify
```

#### Template Not Found
```bash
# Check template path
ls templates/
```

#### Import Errors
```bash
# Ensure all dependencies are downloaded
go mod download
```

---

## Summary

This learning path takes you from:

**Beginner** → Understanding basic routing and handlers
**Intermediate** → Building RESTful APIs with validation
**Advanced** → Implementing complex features like file uploads and sessions
**Production** → Building secure, scalable applications

Each level builds on the previous one, so complete them in order. Don't rush - understanding the concepts is more important than speed.

Good luck with your Gin journey! 🚀
