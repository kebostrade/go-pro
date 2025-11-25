package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/agent"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/common"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// CodingAgentServer provides HTTP API for coding agents
type CodingAgentServer struct {
	config ServerConfig
	server *http.Server
	mux    *http.ServeMux
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port             string
	Agent            *agent.CodingExpertAgent
	LanguageRegistry *common.LanguageRegistry
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	MaxRequestSize   int64
}

// NewCodingAgentServer creates a new API server
func NewCodingAgentServer(config ServerConfig) *CodingAgentServer {
	mux := http.NewServeMux()
	
	server := &CodingAgentServer{
		config: config,
		mux:    mux,
		server: &http.Server{
			Addr:         ":" + config.Port,
			Handler:      mux,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}

	server.registerRoutes()
	return server
}

// registerRoutes sets up HTTP routes
func (s *CodingAgentServer) registerRoutes() {
	// API routes
	s.mux.HandleFunc("/api/v1/coding/ask", s.handleAsk)
	s.mux.HandleFunc("/api/v1/coding/analyze", s.handleAnalyze)
	s.mux.HandleFunc("/api/v1/coding/execute", s.handleExecute)
	s.mux.HandleFunc("/api/v1/coding/debug", s.handleDebug)
	s.mux.HandleFunc("/api/v1/health", s.handleHealth)
	s.mux.HandleFunc("/api/v1/languages", s.handleLanguages)

	// Wrap with middleware
	s.mux.Handle("/", s.loggingMiddleware(s.corsMiddleware(s.mux)))
}

// Start starts the HTTP server
func (s *CodingAgentServer) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *CodingAgentServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// handleAsk handles programming Q&A requests
func (s *CodingAgentServer) handleAsk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req AskRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Validate request
	if req.Query == "" {
		s.sendError(w, http.StatusBadRequest, "Query is required")
		return
	}

	// Run agent
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	result, err := s.config.Agent.Run(ctx, types.AgentInput{
		Query:   req.Query,
		Context: req.Context,
	})

	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Agent error: "+err.Error())
		return
	}

	// Send response
	s.sendJSON(w, http.StatusOK, AskResponse{
		Answer:   result.Output,
		Steps:    len(result.Steps),
		Metadata: result.Metadata,
	})
}

// handleAnalyze handles code analysis requests
func (s *CodingAgentServer) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req AnalyzeRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Validate request
	if req.Code == "" {
		s.sendError(w, http.StatusBadRequest, "Code is required")
		return
	}
	if req.Language == "" {
		s.sendError(w, http.StatusBadRequest, "Language is required")
		return
	}

	// Get language provider
	provider, err := s.config.LanguageRegistry.Get(req.Language)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, "Unsupported language: "+req.Language)
		return
	}

	// Analyze code
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	analysis, err := provider.Analyze(ctx, req.Code)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Analysis error: "+err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, analysis)
}

// handleExecute handles code execution requests
func (s *CodingAgentServer) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ExecuteRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Validate request
	if req.Code == "" {
		s.sendError(w, http.StatusBadRequest, "Code is required")
		return
	}
	if req.Language == "" {
		s.sendError(w, http.StatusBadRequest, "Language is required")
		return
	}

	// Get language provider
	provider, err := s.config.LanguageRegistry.Get(req.Language)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, "Unsupported language: "+req.Language)
		return
	}

	// Execute code
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	result, err := provider.Execute(ctx, types.ExecutionRequest{
		Code:     req.Code,
		Language: req.Language,
		Input:    req.Input,
		Timeout:  req.Timeout,
	})

	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Execution error: "+err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, result)
}

// handleDebug handles debugging requests
func (s *CodingAgentServer) handleDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req DebugRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Build debug query
	query := fmt.Sprintf("Debug this %s code:\n\n%s\n\nError: %s", req.Language, req.Code, req.Error)

	// Run agent
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	result, err := s.config.Agent.Run(ctx, types.AgentInput{
		Query: query,
	})

	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Debug error: "+err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, DebugResponse{
		Explanation: result.Output,
		Suggestions: []string{}, // TODO: Parse suggestions from output
	})
}

// handleHealth handles health check requests
func (s *CodingAgentServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

// handleLanguages lists supported languages
func (s *CodingAgentServer) handleLanguages(w http.ResponseWriter, r *http.Request) {
	languages := s.config.LanguageRegistry.List()
	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"languages": languages,
		"count":     len(languages),
	})
}

// decodeJSON decodes JSON request body
func (s *CodingAgentServer) decodeJSON(r *http.Request, v interface{}) error {
	// Limit request size
	r.Body = http.MaxBytesReader(nil, r.Body, s.config.MaxRequestSize)
	
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	
	return decoder.Decode(v)
}

// sendJSON sends JSON response
func (s *CodingAgentServer) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendError sends error response
func (s *CodingAgentServer) sendError(w http.ResponseWriter, status int, message string) {
	s.sendJSON(w, status, map[string]string{
		"error": message,
	})
}

// loggingMiddleware logs HTTP requests
func (s *CodingAgentServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] %s %s - %v\n", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// corsMiddleware adds CORS headers
func (s *CodingAgentServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Request/Response types

type AskRequest struct {
	Query   string                 `json:"query"`
	Context map[string]interface{} `json:"context,omitempty"`
}

type AskResponse struct {
	Answer   string               `json:"answer"`
	Steps    int                  `json:"steps"`
	Metadata types.AgentMetadata  `json:"metadata"`
}

type AnalyzeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type ExecuteRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Input    string `json:"input,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
}

type DebugRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Error    string `json:"error"`
}

type DebugResponse struct {
	Explanation string   `json:"explanation"`
	Suggestions []string `json:"suggestions"`
}

