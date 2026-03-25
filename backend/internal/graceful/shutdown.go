// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package graceful provides graceful shutdown tools for server
package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownManager handles graceful shutdown of HTTP servers.
type ShutdownManager struct {
	server          *http.Server
	shutdownTimeout time.Duration
	handlers        []func()
	wg              sync.WaitGroup
	mu              sync.Mutex
}

// NewShutdownManager creates a new shutdown manager.
func NewShutdownManager(server *http.Server, shutdownTimeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		server:          server,
		shutdownTimeout: shutdownTimeout,
		handlers:        make([]func(), 0),
		wg:              sync.WaitGroup{},
		mu:              sync.Mutex{},
	}
}

// Register adds a shutdown handler function in the server.
func (sm *ShutdownManager) Register(handler func()) {
	sm.mu.Lock()
	sm.handlers = append(sm.handlers, handler)
	sm.mu.Unlock()
}

// Start begins listening for shutdown signals.
func (sm *ShutdownManager) Start() {
	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), sm.shutdownTimeout)
	defer cancel()

	// Channel to listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-quit

	log.Printf("Shutdown signal received, starting graceful shutdown...")

	// Increment wait group for each handler
	sm.wg.Add(len(sm.handlers))

	// Call all registered shutdown handlers
	for _, handler := range sm.handlers {
		go func(h func()) {
			h()
			sm.wg.Done()
		}(handler)
	}

	// Wait for all handlers to complete or timeout
	go func() {
		select {
		case <-ctx.Done():
			log.Printf("Shutdown timeout exceeded")
		default:
			log.Printf("All shutdown handlers completed")
		}
	}()

	// Wait for all handlers to complete
	sm.wg.Wait()

	// Shutdown the server
	if err := sm.server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Printf("Server shutdown completed")
}
