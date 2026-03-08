package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	// Load configuration
	port := getEnv("GATEWAY_PORT", "8080")
	userServiceURL := getEnv("USER_SERVICE_URL", "http://service-a:8001")
	orderServiceURL := getEnv("ORDER_SERVICE_URL", "http://service-b:8002")

	// Initialize service registry
	registry := NewServiceRegistry(userServiceURL, orderServiceURL, log)

	// Setup Gin router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware(log))
	router.Use(correlationMiddleware())
	router.Use(corsMiddleware())

	// Routes
	router.GET("/health", handleHealth(registry))

	// Proxy routes
	api := router.Group("/api")
	{
		// User service routes
		users := api.Group("/users")
		{
			users.GET("", proxyHandler("users", "/api/users"))
			users.POST("", proxyHandler("users", "/api/users"))
			users.GET("/:id", proxyHandler("users", "/api/users/"))
			users.PUT("/:id", proxyHandler("users", "/api/users/"))
			users.DELETE("/:id", proxyHandler("users", "/api/users/"))
		}

		// Order service routes
		orders := api.Group("/orders")
		{
			orders.GET("", proxyHandler("orders", "/api/orders"))
			orders.POST("", proxyHandler("orders", "/api/orders"))
			orders.GET("/:id", proxyHandler("orders", "/api/orders/"))
			orders.GET("/user/:user_id", proxyHandler("orders", "/api/orders/user/"))
			orders.PUT("/:id/status", proxyHandler("orders", "/api/orders/"))
			orders.DELETE("/:id", proxyHandler("orders", "/api/orders/"))
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.WithField("address", addr).Info("API Gateway starting")

	if err := router.Run(addr); err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}

// handleHealth returns health status of gateway and downstream services
func handleHealth(registry *ServiceRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{
			"status":    "healthy",
			"service":   "api-gateway",
			"timestamp": time.Now().Format(time.RFC3339),
			"services": gin.H{
				"user-service":  checkHealth(registry, "service-a"),
				"order-service": checkHealth(registry, "service-b"),
			},
		}

		c.JSON(http.StatusOK, status)
	}
}

// checkHealth checks the health of a service
func checkHealth(registry *ServiceRegistry, service string) string {
	if err := registry.CheckServiceHealth(service); err != nil {
		return "unhealthy"
	}
	return "healthy"
}

// proxyHandler creates a handler that proxies requests to the specified service
func proxyHandler(service, pathPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		registry := c.MustGet("registry").(*ServiceRegistry)

		// Construct the path to proxy
		path := pathPrefix
		if pathPrefix[len(pathPrefix)-1] != '/' && c.Param("id") != "" {
			path += c.Param("id") + c.Request.URL.Path
		} else {
			path += c.Request.URL.Path
		}

		// Read request body
		var body []byte
		if c.Request.Body != nil {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "failed to read request body",
				})
				return
			}
		}

		// Proxy the request
		resp, err := registry.ProxyRequest(service, path, c.Request.Method, body, c.Request.Header)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "failed to proxy request",
				"message": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := ReadResponseBody(resp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to read response",
			})
			return
		}

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Write response
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// loggingMiddleware logs all requests
func loggingMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get correlation ID from context (set by correlationMiddleware)
		correlationID := c.GetString("correlation_id")

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		log.WithFields(logrus.Fields{
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"status":         c.Writer.Status(),
			"duration_ms":    duration.Milliseconds(),
			"correlation_id": correlationID,
			"client_ip":      c.ClientIP(),
		}).Info("Request processed")
	}
}

// correlationMiddleware adds a correlation ID to each request
func correlationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get correlation ID from header or generate new one
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = generateCorrelationID()
		}

		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-ID", correlationID)

		// Store registry in context
		registry := &ServiceRegistry{
			userServiceURL:  getEnv("USER_SERVICE_URL", "http://service-a:8001"),
			orderServiceURL: getEnv("ORDER_SERVICE_URL", "http://service-b:8002"),
			log:            logrus.New(),
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
		c.Set("registry", registry)

		c.Next()
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Correlation-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// generateCorrelationID generates a unique correlation ID
func generateCorrelationID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
