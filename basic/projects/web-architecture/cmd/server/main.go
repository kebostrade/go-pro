package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/handler"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/middleware"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/repository"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Print banner
	printBanner()

	// Initialize repositories
	userRepo := repository.NewMemoryUserRepository()
	productRepo := repository.NewMemoryProductRepository()

	// Initialize services
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	userService := service.NewUserService(userRepo, jwtSecret)
	productService := service.NewProductService(productRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(chimiddleware.Compress(5))

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(userService))

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Get("/{id}", userHandler.GetByID)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})

			// Product routes
			r.Route("/products", func(r chi.Router) {
				r.Get("/", productHandler.List)
				r.Get("/{id}", productHandler.GetByID)
				r.Post("/", productHandler.Create)
				r.Put("/{id}", productHandler.Update)
				r.Delete("/{id}", productHandler.Delete)
			})
		})
	})

	// Start server
	port := getEnv("PORT", "8080")
	addr := ":" + port

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on http://localhost%s", addr)
		log.Printf("📚 API documentation: http://localhost%s/api/v1", addr)
		log.Printf("❤️  Health check: http://localhost%s/health", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║           🏗️  Web Architecture with Go                      ║
║                                                              ║
║           Clean Architecture • REST API • Middleware        ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	log.Println(banner)
}

