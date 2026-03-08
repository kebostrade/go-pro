package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	version = "1.0.0"
	startTime = time.Now()
)

// Response represents the API response structure
type Response struct {
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	Hostname  string    `json:"hostname"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	Config    Config    `json:"config"`
}

// Config represents application configuration
type Config struct {
	Env      string `json:"env"`
	LogLevel string `json:"log_level"`
	ApiKey   string `json:"api_key_present"`
}

func getConfig() Config {
	// Get environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		apiKey = "yes"
	}

	return Config{
		Env:      getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		ApiKey:   apiKey,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	hostname := getHostname()
	uptime := time.Since(startTime)

	response := Response{
		Message:   "Hello from Kubernetes!",
		Version:   version,
		Hostname:  hostname,
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Config:    getConfig(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func healthLiveHandler(w http.ResponseWriter, r *http.Request) {
	// Liveness probe - check if app is alive
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func healthReadyHandler(w http.ResponseWriter, r *http.Request) {
	// Readiness probe - check if app is ready to serve traffic
	// Add any checks here (database connections, etc.)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	hostname := getHostname()
	uptime := time.Since(startTime)

	metrics := fmt.Sprintf(`
# HELP app_version Application version
# TYPE app_version gauge
app_version{version="%s"} 1

# HELP app_uptime_seconds Application uptime in seconds
# TYPE app_uptime_seconds gauge
app_uptime_seconds %f

# HELP app_hostname Application hostname
# TYPE app_hostname gauge
app_hostname{hostname="%s"} 1
`, version, uptime.Seconds(), hostname)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
}

func main() {
	port := getEnv("PORT", "8080")

	// HTTP handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/health/live", healthLiveHandler)
	http.HandleFunc("/health/ready", healthReadyHandler)
	http.HandleFunc("/metrics", metricsHandler)

	// Log configuration
	log.Printf("Starting application v%s", version)
	log.Printf("Environment: %s", getEnv("APP_ENV", "development"))
	log.Printf("Log Level: %s", getEnv("LOG_LEVEL", "info"))
	log.Printf("Listening on port %s", port)

	// Start server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
