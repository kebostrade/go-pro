// Package main provides the API Gateway entry point.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/microservices/internal/gateway"
	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://service-a:8001"
	}

	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		orderServiceURL = "http://service-b:8002"
	}

	// Create registry and proxy.
	registry := gateway.NewRegistry(userServiceURL, orderServiceURL)
	proxy := gateway.NewProxy(registry)

	// Create chi router.
	r := chi.NewRouter()

	// Health endpoint.
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "gateway",
			"port":    port,
		})
	})

	// Mount routes.
	gateway.MountRoutes(r, proxy)

	addr := ":" + port
	log.Printf("🚀 API Gateway starting on %s", addr)
	log.Printf("📋 Health: http://localhost%s/health", addr)
	log.Printf("🔗 User Service: %s", userServiceURL)
	log.Printf("🔗 Order Service: %s", orderServiceURL)

	// Graceful shutdown with timeout.
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
