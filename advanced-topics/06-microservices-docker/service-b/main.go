package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	// Load configuration
	port := getEnv("SERVICE_PORT", "8002")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "orders_db")
	userServiceURL := getEnv("USER_SERVICE_URL", "http://service-a:8001")

	// Connect to database
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("Failed to ping database")
	}

	log.Info("Database connected successfully")

	// Initialize repository
	repo := NewRepository(db, userServiceURL, log)
	if err := repo.Initialize(); err != nil {
		log.WithError(err).Fatal("Failed to initialize repository")
	}

	// Initialize handlers
	handlers := NewHandlers(repo, log)

	// Setup Gin router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware(log))
	router.Use(correlationMiddleware())

	// Routes
	router.GET("/health", handlers.Health)

	v1 := router.Group("/api")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("", handlers.GetOrders)
			orders.GET("/:id", handlers.GetOrder)
			orders.GET("/user/:user_id", handlers.GetUserOrders)
			orders.PUT("/:id/status", handlers.UpdateOrderStatus)
			orders.DELETE("/:id", handlers.DeleteOrder)
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.WithField("address", addr).Info("Order service starting")

	if err := router.Run(addr); err != nil {
		log.WithError(err).Fatal("Failed to start server")
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
