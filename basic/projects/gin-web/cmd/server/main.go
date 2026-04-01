package main

import (
	"log"

	"github.com/DimaJoyti/go-pro/basic/projects/gin-web/internal/handler"
	"github.com/DimaJoyti/go-pro/basic/projects/gin-web/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create new Gin router
	router := gin.New()

	// Core middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Custom middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	// Static files
	router.Static("/static", "./static")

	// Templates
	router.LoadHTMLGlob("internal/views/*.html")

	// Routes
	router.GET("/", handler.Home)
	router.GET("/about", handler.About)
	router.GET("/api/v1/health", handler.HealthCheck)

	// Start server
	log.Println("Starting Gin Web App on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
