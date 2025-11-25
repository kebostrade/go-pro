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
		targetURL := buildTargetURL(serviceAddr, r.URL.Path, r.URL.RawQuery)

		// Create and send proxy request
		resp, err := sendProxyRequest(r, targetURL)
		if err != nil {
			logger.Error("Failed to proxy request", zap.Error(err))
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Copy response headers and body
		copyResponseHeaders(w, resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

// buildTargetURL constructs the target URL for the proxy request
func buildTargetURL(serviceAddr, path, rawQuery string) string {
	targetURL := fmt.Sprintf("http://%s%s", serviceAddr, path)
	if rawQuery != "" {
		targetURL += "?" + rawQuery
	}
	return targetURL
}

// sendProxyRequest creates and sends a proxy request to the target service
func sendProxyRequest(r *http.Request, targetURL string) (*http.Response, error) {
	// #nosec G602: URL is constructed from trusted service discovery source
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy request: %w", err)
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

	return client.Do(proxyReq)
}

// copyResponseHeaders copies response headers from the upstream service
func copyResponseHeaders(w http.ResponseWriter, headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}
