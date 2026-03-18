// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package main is the entry point for the GO-PRO Learning Platform API server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go-pro-backend/internal/config"
	"go-pro-backend/internal/container"
	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/handler"
	"go-pro-backend/internal/middleware"
	"go-pro-backend/internal/service"
	applogger "go-pro-backend/pkg/logger"
)

const version = "1.0.0"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	applog := applogger.New(cfg.Logger.Level, cfg.Logger.Format)

	ctx := context.Background()
	applog.Info(ctx, "Starting GO-PRO API Server",
		"version", version,
		"port", cfg.Server.Port,
		"log_level", cfg.Logger.Level,
	)

	appContainer, err := container.NewContainer(&container.ContainerConfig{
		Config: cfg,
		Logger: applog,
	})
	if err != nil {
		applog.Error(ctx, "Failed to initialize container", "error", err)
		os.Exit(1)
	}
	defer func() {
		if shutdownErr := appContainer.Shutdown(ctx); shutdownErr != nil {
			applog.Error(ctx, "Failed to shutdown container", "error", shutdownErr)
		}
	}()

	if err := initializeSampleData(ctx, appContainer.Services); err != nil {
		applog.Error(ctx, "Failed to initialize sample data", "error", err)
		os.Exit(1)
	}

	if err := appContainer.Services.Auth.Initialize(ctx); err != nil {
		devMode := os.Getenv("DEV_MODE") == "true"
		projectID := os.Getenv("FIREBASE_PROJECT_ID")

		if devMode && projectID == "" {
			applog.Warn(ctx, "Firebase not configured, continuing in mock DEV_MODE")
		} else {
			applog.Error(ctx, "Failed to initialize Firebase Auth service", "error", err)
			if !devMode {
				os.Exit(1)
			}
		}
	} else {
		applog.Info(ctx, "Firebase Auth service initialized successfully")
	}

	authServiceAdapter := &firebaseAuthAdapter{authService: appContainer.Services.Auth}

	authMiddleware := middleware.NewAuthMiddleware(
		authServiceAdapter,
		appContainer.Repositories.User,
		applog,
	)

	httpHandler := handler.New(appContainer.Services, applog, appContainer.Validator)
	authHandler := handler.NewAuthHandler(appContainer.Services, applog, appContainer.Validator)
	adminHandler := handler.NewAdminHandler(appContainer.Services, applog, appContainer.Validator)

	// Initialize AI-powered playground handler if AgentPool is available
	if appContainer.AgentPool != nil {
		aiHandler := handler.NewPlaygroundAIHandler(appContainer.AgentPool, appContainer.Services.Executor, applog)
		httpHandler.SetAIHandler(aiHandler)
		applog.Info(ctx, "AI-powered playground handler initialized")
	}

	// Setup routes.
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux, authMiddleware)
	authHandler.RegisterRoutes(mux, authMiddleware)
	adminHandler.RegisterRoutes(mux, authMiddleware)

	// Setup middleware chain.
	middlewares := []middleware.Middleware{
		middleware.RequestID,
		middleware.Logging(applog),
		middleware.Recovery(applog),
		middleware.CORS(cfg.CORS.AllowedOrigins),
		middleware.Security(),
		middleware.ContentType("application/json"),
		middleware.Timeout(30 * time.Second),
		middleware.RateLimit(100, time.Minute),
		middleware.Pagination(10, 100),
	}

	handler := middleware.Chain(mux, middlewares...)

	// Create HTTP server.
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		applog.Info(ctx, "GO-PRO API Server starting",
			"address", "http://"+cfg.Server.Host+":"+cfg.Server.Port,
			"documentation", "http://"+cfg.Server.Host+":"+cfg.Server.Port,
			"health_check", "http://"+cfg.Server.Host+":"+cfg.Server.Port+"/api/v1/health",
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			applog.Error(ctx, "Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	applog.Info(ctx, "Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		applog.Error(ctx, "Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	applog.Info(ctx, "Server exited gracefully")
}

// initializeSampleData populates the repositories with sample data.
func initializeSampleData(ctx context.Context, services *service.Services) error {
	// Create sample course.
	courseReq := &domain.CreateCourseRequest{
		Title: "GO-PRO: Complete Go Programming Mastery",
		Description: "Master Go programming from basics to advanced microservices. " +
			"Learn Go's syntax, concurrency patterns, web development, testing, " +
			"and best practices through hands-on exercises and real-world projects.",
	}

	course, err := services.Course.CreateCourse(ctx, courseReq)
	if err != nil {
		return err
	}

	// Add more sample data here as services are implemented.
	_ = course // Prevent unused variable warning

	return nil
}

// firebaseAuthAdapter adapts service.AuthService to middleware.AuthService interface.
type firebaseAuthAdapter struct {
	authService service.AuthService
}

// VerifyToken verifies a Firebase ID token and returns the token information.
func (a *firebaseAuthAdapter) VerifyToken(ctx context.Context, token string) (*middleware.FirebaseToken, error) {
	claims, err := a.authService.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &middleware.FirebaseToken{
		UID:         claims.UserID,
		Email:       claims.Email,
		DisplayName: claims.DisplayName,
		PhotoURL:    claims.Picture,
	}, nil
}
