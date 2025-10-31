package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/discovery"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/logger"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	serviceName = "api-gateway"
	servicePort = "8080"
)

func main() {
	if err := logger.Init(serviceName, true); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting API Gateway", zap.String("port", servicePort))

	// Wait for services to register
	time.Sleep(2 * time.Second)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RateLimitMiddleware(200, 400))

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Service discovery endpoint
	r.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		services := discovery.List()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"services": %v}`, services)
	}).Methods("GET")

	// Route to User Service
	r.PathPrefix("/api/users").HandlerFunc(proxyHandler("user-service"))

	// Route to Product Service
	r.PathPrefix("/api/products").HandlerFunc(proxyHandler("product-service"))

	// Route to Order Service
	r.PathPrefix("/api/orders").HandlerFunc(proxyHandler("order-service"))

	srv := &http.Server{
		Addr:         ":" + servicePort,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("API Gateway listening", zap.String("address", srv.Addr))
		logger.Info("Available routes:")
		logger.Info("  GET  /health")
		logger.Info("  GET  /services")
		logger.Info("  *    /api/users/*    -> user-service")
		logger.Info("  *    /api/products/* -> product-service")
		logger.Info("  *    /api/orders/*   -> order-service")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down API Gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("API Gateway stopped")
}

// proxyHandler creates a reverse proxy handler for a service
func proxyHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Discover service address
		serviceAddr, err := discovery.Discover(serviceName)
		if err != nil {
			logger.Error("Service not found", zap.String("service", serviceName), zap.Error(err))
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		// Build target URL
		targetURL := fmt.Sprintf("http://%s%s", serviceAddr, r.URL.Path)
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}

		// Create new request
		proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			logger.Error("Failed to create proxy request", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Send request
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(proxyReq)
		if err != nil {
			logger.Error("Failed to proxy request", zap.Error(err))
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Copy status code
		w.WriteHeader(resp.StatusCode)

		// Copy response body
		io.Copy(w, resp.Body)

		logger.Debug("Proxied request",
			zap.String("service", serviceName),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", resp.StatusCode),
		)
	}
}

