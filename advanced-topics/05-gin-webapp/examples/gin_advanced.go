// Gin Advanced - Advanced Gin framework features
//
// This example covers:
// - File upload handling (single and multiple files)
// - Session management

//go:build ignore
// - Advanced request binding
// - Custom validation
// - Streaming responses
// - Graceful shutdown
// - Custom middleware
// - Cookie management
// - Request context
// - Error recovery
//
// Run it: go run examples/gin_advanced.go
//
// Test file upload:
//   curl -X POST http://localhost:8080/upload -F "file=@test.txt" -F "name=Test File"

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Session represents a simple in-memory session
type Session struct {
	ID        string
	Data      map[string]interface{}
	ExpiresAt time.Time
}

// SessionManager manages sessions
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}
	// Start cleanup goroutine
	go sm.cleanupExpiredSessions()
	return sm
}

// CreateSession creates a new session
func (sm *SessionManager) CreateSession() *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		ID:        fmt.Sprintf("session_%d", time.Now().UnixNano()),
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	sm.sessions[session.ID] = session
	return session
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[id]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}

// cleanupExpiredSessions removes expired sessions
func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for id, session := range sm.sessions {
			if now.After(session.ExpiresAt) {
				delete(sm.sessions, id)
			}
		}
		sm.mu.Unlock()
	}
}

// FileMetadata represents uploaded file metadata
type FileMetadata struct {
	Filename string
	Size     int64
	Type     string
	Path     string
}

// Product represents a product with custom validation
type Product struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Description string  `json:"description" binding:"max=500"`
	SKU         string  `json:"sku" binding:"required,sku-valid"`
	Category    string  `json:"category" binding:"required,oneof=electronics books clothing food"`
}

// Custom validator for SKU
func skuValid(fl validator.FieldLevel) bool {
	sku := fl.Field().String()
	if len(sku) < 3 || len(sku) > 50 {
		return false
	}
	// Simple validation: should contain only alphanumeric and hyphens
	for _, c := range sku {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

// Upload directory
const uploadDir = "./uploads"

func setupRouter() *gin.Engine {
	// Create router
	router := gin.Default()

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("sku-valid", skuValid)
	}

	// Ensure upload directory exists
	os.MkdirAll(uploadDir, 0755)

	// Initialize session manager
	sessionManager := NewSessionManager()

	// Custom middleware - Request ID
	router.Use(func(c *gin.Context) {
		c.Set("request_id", fmt.Sprintf("req_%d", time.Now().UnixNano()))
		c.Next()
	})

	// Custom middleware - Security headers
	router.Use(func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})

	// File upload routes
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No file uploaded",
			})
			return
		}

		// Get additional form fields
		name := c.PostForm("name")
		description := c.PostForm("description")

		// Create unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filepath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file",
			})
			return
		}

		// Get file info
		fileInfo, _ := os.Stat(filepath)

		metadata := FileMetadata{
			Filename: filename,
			Size:     fileInfo.Size(),
			Type:     file.Header.Get("Content-Type"),
			Path:     filepath,
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "File uploaded successfully",
			"metadata":    metadata,
			"name":        name,
			"description": description,
		})
	})

	// Multiple file upload
	router.POST("/upload/multiple", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse multipart form",
			})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No files uploaded",
			})
			return
		}

		var uploadedFiles []FileMetadata

		for _, file := range files {
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			filepath := filepath.Join(uploadDir, filename)

			if err := c.SaveUploadedFile(file, filepath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to save file: %s", file.Filename),
				})
				return
			}

			fileInfo, _ := os.Stat(filepath)

			uploadedFiles = append(uploadedFiles, FileMetadata{
				Filename: filename,
				Size:     fileInfo.Size(),
				Type:     file.Header.Get("Content-Type"),
				Path:     filepath,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully uploaded %d files", len(uploadedFiles)),
			"files":   uploadedFiles,
		})
	})

	// Session routes
	router.POST("/session/create", func(c *gin.Context) {
		session := sessionManager.CreateSession()

		// Set session cookie
		c.SetCookie(
			"session_id",
			session.ID,
			86400, // 1 day
			"/",
			"",
			false, // secure
			true,  // http only
		)

		c.JSON(http.StatusOK, gin.H{
			"message":  "Session created",
			"session_id": session.ID,
		})
	})

	router.POST("/session/set", func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No session found",
			})
			return
		}

		session, exists := sessionManager.GetSession(sessionID)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session",
			})
			return
		}

		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})
			return
		}

		// Store data in session
		for key, value := range data {
			session.Data[key] = value
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Session data updated",
			"data":    session.Data,
		})
	})

	router.GET("/session/get", func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No session found",
			})
			return
		}

		session, exists := sessionManager.GetSession(sessionID)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"session_id": session.ID,
			"data":       session.Data,
		})
	})

	router.POST("/session/delete", func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No session found",
			})
			return
		}

		// Clear session cookie
		c.SetCookie(
			"session_id",
			"",
			-1,
			"/",
			"",
			false,
			true,
		)

		// Remove session from manager
		delete(sessionManager.sessions, sessionID)

		c.JSON(http.StatusOK, gin.H{
			"message": "Session deleted",
		})
	})

	// Advanced binding with custom validation
	router.POST("/products", func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Product created successfully",
			"product": product,
		})
	})

	// Streaming response example
	router.GET("/stream", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		c.Stream(func(w io.Writer) bool {
			for i := 1; i <= 10; i++ {
				fmt.Fprintf(w, "Message %d\n", i)
				w.(http.Flusher).Flush()
				time.Sleep(500 * time.Millisecond)
			}
			return false
		})
	})

	// Server-Sent Events (SSE)
	router.GET("/events", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		c.Stream(func(w io.Writer) bool {
			for i := 1; i <= 5; i++ {
				event := fmt.Sprintf("event: message\ndata: Event %d at %s\n\n", i, time.Now().Format(time.RFC3339))
				fmt.Fprint(w, event)
				w.(http.Flusher).Flush()
				time.Sleep(1 * time.Second)
			}
			return false
		})
	})

	// Request context example
	router.GET("/context", func(c *gin.Context) {
		requestID, _ := c.Get("request_id")

		c.JSON(http.StatusOK, gin.H{
			"request_id":   requestID,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"user_agent":   c.Request.UserAgent(),
			"client_ip":    c.ClientIP(),
			"content_type": c.ContentType(),
		})
	})

	// Cookie management
	router.POST("/cookies/set", func(c *gin.Context) {
		var data struct {
			Name  string `json:"name" binding:"required"`
			Value string `json:"value" binding:"required"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})
			return
		}

		c.SetCookie(
			data.Name,
			data.Value,
			3600,   // 1 hour
			"/",    // path
			"",     // domain
			false,  // secure
			false,  // http only
		)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Cookie '%s' set successfully", data.Name),
		})
	})

	router.GET("/cookies/get", func(c *gin.Context) {
		cookieName := c.Query("name")
		if cookieName == "" {
			// Return all cookies
			cookies := c.Request.Cookies()
			cookieMap := make(map[string]string)
			for _, cookie := range cookies {
				cookieMap[cookie.Name] = cookie.Value
			}

			c.JSON(http.StatusOK, gin.H{
				"cookies": cookieMap,
			})
			return
		}

		cookie, err := c.Cookie(cookieName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Cookie '%s' not found", cookieName),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name":  cookieName,
			"value": cookie,
		})
	})

	router.DELETE("/cookies/delete", func(c *gin.Context) {
		var data struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})
			return
		}

		c.SetCookie(
			data.Name,
			"",
			-1,
			"/",
			"",
			false,
			false,
		)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Cookie '%s' deleted successfully", data.Name),
		})
	})

	// Custom error handling
	router.GET("/error/:code", func(c *gin.Context) {
		codeStr := c.Param("code")
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			code = 500
		}

		c.JSON(code, gin.H{
			"error": fmt.Sprintf("Error %d occurred", code),
			"path":  c.Request.URL.Path,
		})
	})

	// Panic recovery (custom)
	router.GET("/panic", func(c *gin.Context) {
		panic("This is a test panic!")
	})

	// Health check with detailed info
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(time.Now()).String(),
		})
	})

	return router
}

func main() {
	router := setupRouter()

	// Create server with timeouts
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("🚀 Gin Advanced server starting on :8080")
	fmt.Println("\n📝 Available endpoints:")
	fmt.Println("  POST   /upload")
	fmt.Println("  POST   /upload/multiple")
	fmt.Println("  POST   /session/create")
	fmt.Println("  POST   /session/set")
	fmt.Println("  GET    /session/get")
	fmt.Println("  POST   /session/delete")
	fmt.Println("  POST   /products")
	fmt.Println("  GET    /stream")
	fmt.Println("  GET    /events")
	fmt.Println("  GET    /context")
	fmt.Println("  POST   /cookies/set")
	fmt.Println("  GET    /cookies/get")
	fmt.Println("  DELETE /cookies/delete")
	fmt.Println("  GET    /error/:code")
	fmt.Println("  GET    /panic")
	fmt.Println("  GET    /health")
	fmt.Println("\n💡 Example commands:")
	fmt.Println("  # Upload single file")
	fmt.Println("  curl -X POST http://localhost:8080/upload \\")
	fmt.Println("    -F 'file=@test.txt' \\")
	fmt.Println("    -F 'name=Test File'")
	fmt.Println()
	fmt.Println("  # Upload multiple files")
	fmt.Println("  curl -X POST http://localhost:8080/upload/multiple \\")
	fmt.Println("    -F 'files=@test1.txt' \\")
	fmt.Println("    -F 'files=@test2.txt'")
	fmt.Println()
	fmt.Println("  # Create session")
	fmt.Println("  curl -X POST http://localhost:8080/session/create \\")
	fmt.Println("    -c cookies.txt")
	fmt.Println()
	fmt.Println("  # Set session data")
	fmt.Println("  curl -X POST http://localhost:8080/session/set \\")
	fmt.Println("    -b cookies.txt \\")
	fmt.Println("    -H 'Content-Type: application/json' \\")
	fmt.Println("    -d '{\"user\":\"john\",\"role\":\"admin\"}'")
	fmt.Println()
	fmt.Println("  # Get session data")
	fmt.Println("  curl http://localhost:8080/session/get -b cookies.txt")
	fmt.Println()
	fmt.Println("  # Create product with validation")
	fmt.Println("  curl -X POST http://localhost:8080/products \\")
	fmt.Println("    -H 'Content-Type: application/json' \\")
	fmt.Println("    -d '{\"name\":\"Laptop\",\"price\":999.99,\"sku\":\"LAPTOP-001\",\"category\":\"electronics\"}'")
	fmt.Println()
	fmt.Println("  # Stream response")
	fmt.Println("  curl http://localhost:8080/stream")
	fmt.Println()
	fmt.Println("  # Server-Sent Events")
	fmt.Println("  curl http://localhost:8080/events")
	fmt.Println()

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	fmt.Println("✅ Server started successfully")
	fmt.Println("\n⚡ Press Ctrl+C to gracefully shutdown the server")
	fmt.Println()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\n⏳ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("✅ Server shutdown complete")
}
