package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goproject/ml-gorgonia/internal/api"
	"github.com/goproject/ml-gorgonia/internal/model"
)

func main() {
	log.Println("Starting ML Inference Server...")

	// Load model if path provided, otherwise use mock
	modelPath := os.Getenv("MODEL_PATH")
	var m model.Model
	if modelPath != "" {
		var err error
		m, err = model.LoadONNXModel(modelPath)
		if err != nil {
			log.Printf("Warning: Failed to load model from %s: %v", modelPath, err)
			log.Println("Using mock model for demonstration")
			m = model.MockModel()
		}
	} else {
		log.Println("No MODEL_PATH provided, using mock model")
		m = model.MockModel()
	}

	// Create handler
	handler := api.NewHandler(m)
	router := handler.Routes()

	// Create server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
