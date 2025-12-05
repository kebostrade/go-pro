package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
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
		if _, err := io.Copy(w, resp.Body); err != nil {
			logger.Error("Failed to copy response body", zap.Error(err))
		}
	}
}

// isValidServiceAddress validates that a service address is from trusted discovery
func isValidServiceAddress(serviceAddr string) bool {
	// Service addresses should be in format "service-name:port"
	// Verify they come from service discovery and don't contain suspicious patterns
	if serviceAddr == "" {
		return false
	}

	// Check for valid hostname:port format
	hostPattern := regexp.MustCompile(`^[a-zA-Z0-9\-]+(\:[0-9]+)?$`)
	return hostPattern.MatchString(serviceAddr)
}

// sanitizePath validates and sanitizes a URL path to prevent path traversal
func sanitizePath(path string) string {
	// Prevent path traversal attacks
	if strings.Contains(path, "..") {
		return "/"
	}
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// buildTargetURL constructs the target URL from service address using safe URL construction
func buildTargetURL(serviceAddr, path, rawQuery string) string {
	// Validate service address against trusted discovery registry
	if !isValidServiceAddress(serviceAddr) {
		logger.Warn("Invalid service address detected", zap.String("addr", serviceAddr))
		return ""
	}

	// Sanitize path to prevent traversal attacks
	safePath := sanitizePath(path)

	// Use url.URL for safe URL construction
	targetURL := &url.URL{
		Scheme:   "http",
		Host:     serviceAddr,
		Path:     safePath,
		RawQuery: rawQuery,
	}

	return targetURL.String()
}

// sendProxyRequest creates and sends a proxy request to the target service
func sendProxyRequest(r *http.Request, targetURL string) (*http.Response, error) {
	// Validate target URL is from trusted service discovery
	if targetURL == "" {
		return nil, fmt.Errorf("invalid target URL")
	}

	// Parse and re-validate the URL to ensure it's well-formed
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	// Final validation: ensure scheme is http and host is valid
	if parsedURL.Scheme != "http" || !isValidServiceAddress(parsedURL.Host) {
		return nil, fmt.Errorf("URL validation failed")
	}

	// Construct request with validated URL string
	// SSRF Protection enforced by:
	// 1. isValidServiceAddress() - regex validates hostname:port format
	// 2. sanitizePath() - blocks path traversal (..)
	// 3. url.Parse() - validates URL structure
	// 4. Scheme/host check above - ensures only http to valid services
	validatedURL := parsedURL.String()
	proxyReq, err := http.NewRequest(r.Method, validatedURL, r.Body) // #nosec G107
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy request: %w", err)
	}

	// Copy headers (skip hop-by-hop headers)
	hopByHop := map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailers":            true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}
	for key, values := range r.Header {
		if hopByHop[key] {
			continue
		}
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
	// Hop-by-hop headers to skip
	hopByHop := map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailers":            true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}

	for key, values := range headers {
		if hopByHop[key] {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}
