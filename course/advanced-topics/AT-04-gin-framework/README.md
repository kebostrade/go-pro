# Building Web Applications with Go and Gin

Create production web applications using the Gin framework.

## Learning Objectives

- Set up Gin applications
- Handle routing and parameters
- Implement middleware
- Validate request data
- Serve static files and templates
- Handle errors gracefully

## Theory

### Basic Gin Setup

```go
func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080")
}
```

### Routing Patterns

```go
func setupRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    {
        api.GET("/users", listUsers)
        api.GET("/users/:id", getUser)
        api.POST("/users", createUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }

    admin := r.Group("/admin")
    admin.Use(authMiddleware)
    {
        admin.GET("/stats", getStats)
    }
}
```

### Request Binding

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
}

func createUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
}
```

### Custom Middleware

```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        latency := time.Since(start)
        statusCode := c.Writer.Status()

        log.Printf("[%s] %s %d %v",
            c.Request.Method, path, statusCode, latency)
    }
}

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

### Error Handling

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            switch e := err.(type) {
            case *ValidationError:
                c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
            case *NotFoundError:
                c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            }
        }
    }
}
```

## Security Considerations

```go
func RateLimitMiddleware(r *rate.Limiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !r.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Next()
    }
}
```

## Performance Tips

```go
var templateCache map[string]*template.Template

func LoadTemplates() {
    templateCache = make(map[string]*template.Template)
    for _, name := range []string{"index", "about", "contact"} {
        templateCache[name] = template.Must(
            template.ParseFiles("templates/"+name+".html"),
        )
    }
}

func renderTemplate(c *gin.Context, name string, data interface{}) {
    tmpl, ok := templateCache[name]
    if !ok {
        c.JSON(500, gin.H{"error": "template not found"})
        return
    }
    c.HTML(200, name, data)
}
```

## Exercises

1. Build a blog API with CRUD operations
2. Implement JWT authentication middleware
3. Add request validation
4. Create a file upload endpoint

## Validation

```bash
cd exercises
go test -v ./...
curl http://localhost:8080/api/v1/users
```

## Key Takeaways

- Use route groups for organization
- Validate all inputs with binding tags
- Chain middleware for cross-cutting concerns
- Handle errors consistently
- Cache templates in production

## Next Steps

**[AT-05: Microservices with Docker](../AT-05-microservices-docker/README.md)**

---

Gin makes web development fast and clean. 🚀
