package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	appInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application information",
		},
		[]string{"version", "environment"},
	)
)

func init() {
	// Register metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(appInfo)

	// Set app info
	version := getEnv("APP_VERSION", "1.0.0")
	environment := getEnv("ENVIRONMENT", "development")
	appInfo.WithLabelValues(version, environment).Set(1)
}

type Response struct {
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func main() {
	printBanner()

	port := getEnv("PORT", "8080")
	
	r := mux.NewRouter()

	// Middleware
	r.Use(loggingMiddleware)
	r.Use(metricsMiddleware)

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/api/hello", helloHandler).Methods("GET")
	r.HandleFunc("/api/info", infoHandler).Methods("GET")
	
	// Health checks
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/health/live", livenessHandler).Methods("GET")
	r.HandleFunc("/health/ready", readinessHandler).Methods("GET")
	
	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// HTTP Server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%s", port)
		log.Printf("📊 Metrics available at http://localhost:%s/metrics", port)
		log.Printf("💚 Health check at http://localhost:%s/health", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
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

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Message:   "Welcome to DevOps with Go!",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"version":     getEnv("APP_VERSION", "1.0.0"),
			"environment": getEnv("ENVIRONMENT", "development"),
		},
	}
	jsonResponse(w, http.StatusOK, resp)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	resp := Response{
		Message:   fmt.Sprintf("Hello, %s!", name),
		Timestamp: time.Now().Format(time.RFC3339),
	}
	jsonResponse(w, http.StatusOK, resp)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	
	resp := Response{
		Message:   "Application Information",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"hostname":    hostname,
			"version":     getEnv("APP_VERSION", "1.0.0"),
			"environment": getEnv("ENVIRONMENT", "development"),
			"go_version":  "1.21",
		},
	}
	jsonResponse(w, http.StatusOK, resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	checks := map[string]string{
		"application": "healthy",
		"database":    "healthy", // Simulated
		"cache":       "healthy", // Simulated
	}

	resp := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    checks,
	}
	jsonResponse(w, http.StatusOK, resp)
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	// Liveness probe - is the app running?
	resp := map[string]string{
		"status": "alive",
	}
	jsonResponse(w, http.StatusOK, resp)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Readiness probe - is the app ready to serve traffic?
	resp := map[string]string{
		"status": "ready",
	}
	jsonResponse(w, http.StatusOK, resp)
}

// Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.RequestURI, rw.statusCode, duration)
	})
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start).Seconds()
		
		// Record metrics
		httpRequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			fmt.Sprintf("%d", rw.statusCode),
		).Inc()
		
		httpRequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(duration)
	})
}

// Response writer wrapper
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Helpers
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
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
║        🐳 DevOps with Go: Docker, K8s, Terraform            ║
║                                                              ║
║        A production-ready Go application demonstrating:     ║
║        • Docker containerization                            ║
║        • Kubernetes orchestration                           ║
║        • Terraform infrastructure                           ║
║        • Prometheus metrics                                 ║
║        • Health checks                                      ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

