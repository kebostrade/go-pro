package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/discovery"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/logger"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/middleware"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/services/user-service/internal"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	serviceName = "user-service"
	servicePort = "8081"
)

func main() {
	// Initialize logger
	if err := logger.Init(serviceName, true); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting User Service", zap.String("port", servicePort))

	// Create repository
	repo := internal.NewInMemoryRepository()

	// Create handler
	handler := internal.NewHandler(repo)

	// Setup router
	r := mux.NewRouter()

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RateLimitMiddleware(100, 200)) // 100 requests per second, burst of 200

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Register routes
	handler.RegisterRoutes(r)

	// Register service with discovery
	serviceAddr := fmt.Sprintf("localhost:%s", servicePort)
	discovery.Register(serviceName, serviceAddr)
	defer discovery.Deregister(serviceName)

	// Create server
	srv := &http.Server{
		Addr:         ":" + servicePort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("User Service listening", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down User Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("User Service stopped")
}

