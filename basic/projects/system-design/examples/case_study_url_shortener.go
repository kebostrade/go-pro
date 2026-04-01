// Package main provides the URL shortener case study demonstrating clean architecture.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goproject/system-design/internal/cache"
	"github.com/goproject/system-design/internal/circuit"
	"github.com/goproject/system-design/internal/concurrency"
)

// URL represents a shortened URL.
type URL struct {
	Code           string
	OriginalURL    string
	ShortURL       string
	AccessCount    int
	CreatedAt      time.Time
	LastAccessedAt time.Time
}

// URLRepository defines the interface for URL storage.
type URLRepository interface {
	Create(ctx context.Context, url *URL) error
	GetByCode(ctx context.Context, code string) (*URL, error)
	IncrementAccess(ctx context.Context, code string) error
}

// InMemoryURLRepository implements URLRepository with in-memory storage.
type InMemoryURLRepository struct {
	urls map[string]*URL
}

// NewInMemoryURLRepository creates a new in-memory URL repository.
func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		urls: make(map[string]*URL),
	}
}

// Create creates a new shortened URL.
func (r *InMemoryURLRepository) Create(ctx context.Context, url *URL) error {
	if _, exists := r.urls[url.Code]; exists {
		return fmt.Errorf("code already exists")
	}
	url.CreatedAt = time.Now()
	r.urls[url.Code] = url
	return nil
}

// GetByCode retrieves a URL by its short code.
func (r *InMemoryURLRepository) GetByCode(ctx context.Context, code string) (*URL, error) {
	url, exists := r.urls[code]
	if !exists {
		return nil, fmt.Errorf("URL not found")
	}
	return url, nil
}

// IncrementAccess increments the access count for a URL.
func (r *InMemoryURLRepository) IncrementAccess(ctx context.Context, code string) error {
	url, exists := r.urls[code]
	if !exists {
		return fmt.Errorf("URL not found")
	}
	url.AccessCount++
	url.LastAccessedAt = time.Now()
	return nil
}

// URLService handles URL shortening business logic.
type URLService struct {
	repo       URLRepository
	cache      *cache.Cache
	cb         *circuit.CircuitBreaker
	workerPool *concurrency.WorkerPool
}

// processAnalytics processes analytics in the background.
func processAnalytics(item concurrency.WorkItem) concurrency.Result {
	// In real implementation, this would access the repo
	return concurrency.Result{ID: item.ID}
}

// NewURLService creates a new URL service.
func NewURLService(repo URLRepository) *URLService {
	s := &URLService{
		repo:       repo,
		cache:      cache.New(cache.DefaultExpiration, cache.MaxSize),
		cb:         circuit.NewCircuitBreaker("url-db"),
		workerPool: concurrency.NewWorkerPool(10, processAnalytics),
	}
	s.workerPool.Start()
	return s
}

// Shorten creates a short URL.
func (s *URLService) Shorten(ctx context.Context, originalURL string) (*URL, error) {
	// Check cache first
	cacheKey := originalURL
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*URL), nil
	}

	// Generate short code
	code := generateShortCode(6)
	url := &URL{
		Code:        code,
		OriginalURL: originalURL,
		ShortURL:    fmt.Sprintf("http://localhost:8080/%s", code),
	}

	// Create with circuit breaker protection
	err := s.cb.Execute(ctx, func() error {
		return s.repo.Create(ctx, url)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, url)

	return url, nil
}

// Redirect retrieves the original URL and increments access count.
func (s *URLService) Redirect(ctx context.Context, code string) (*URL, error) {
	// Check cache first
	cacheKey := "redirect:" + code
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*URL), nil
	}

	var url *URL
	var err error

	// Get URL with circuit breaker protection
	s.cb.Execute(ctx, func() error {
		url, err = s.repo.GetByCode(ctx, code)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Increment access asynchronously via worker pool
	s.workerPool.Submit(concurrency.WorkItem{
		ID:      code,
		Payload: url.Code,
	})

	// Cache the result briefly
	s.cache.SetWithTTL(cacheKey, url, 5*time.Minute)

	return url, nil
}

// GetStats returns statistics for a URL.
func (s *URLService) GetStats(ctx context.Context, code string) (*URL, error) {
	return s.repo.GetByCode(ctx, code)
}

// generateShortCode generates a random short code.
func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// HTTPHandlers returns HTTP handlers for URL operations.
type HTTPHandlers struct {
	service *URLService
}

func (h *HTTPHandlers) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := parseJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidURL(req.URL) {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	url, err := h.service.Shorten(r.Context(), req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, url)
}

func (h *HTTPHandlers) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.service.Redirect(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

func (h *HTTPHandlers) Stats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.service.GetStats(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, url)
}

func (h *HTTPHandlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

// Helper functions
func parseJSON(r *http.Request, v interface{}) error {
	return nil // Simplified
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	fmt.Fprintf(w, "%+v", v)
}

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func main() {
	log.Println("URL Shortener Case Study Starting...")

	// Initialize components
	repo := NewInMemoryURLRepository()
	service := NewURLService(repo)
	handlers := &HTTPHandlers{service: service}

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/health", handlers.Health)
	r.Post("/api/shorten", handlers.Shorten)
	r.Get("/{code}", handlers.Redirect)
	r.Get("/api/stats/{code}", handlers.Stats)

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
