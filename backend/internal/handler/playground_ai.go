// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package handler provides HTTP request handlers for the API.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-pro-backend/internal/agents"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/service"
	"go-pro-backend/pkg/logger"
)

// PlaygroundAIHandler handles AI-powered playground requests.
type PlaygroundAIHandler struct {
	agentPool   *agents.AgentPool
	executor    service.ExecutorService
	logger      logger.Logger
	rateLimiter *playgroundRateLimiter
}

// NewPlaygroundAIHandler creates a new AI handler for playground.
func NewPlaygroundAIHandler(agentPool *agents.AgentPool, executor service.ExecutorService, logger logger.Logger) *PlaygroundAIHandler {
	return &PlaygroundAIHandler{
		agentPool:   agentPool,
		executor:    executor,
		logger:      logger,
		rateLimiter: newPlaygroundRateLimiter(),
	}
}

// playgroundRateLimiter manages rate limiting for AI requests.
type playgroundRateLimiter struct {
	requests map[string]*rateLimitEntry
}

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

func newPlaygroundRateLimiter() *playgroundRateLimiter {
	return &playgroundRateLimiter{
		requests: make(map[string]*rateLimitEntry),
	}
}

func (r *playgroundRateLimiter) allow(clientID string, limit int) bool {
	now := time.Now()
	entry, exists := r.requests[clientID]
	if !exists || now.After(entry.resetTime) {
		r.requests[clientID] = &rateLimitEntry{
			count:     1,
			resetTime: now.Add(time.Minute),
		}
		return true
	}
	if entry.count >= limit {
		return false
	}
	entry.count++
	return true
}

// AIAnalysisRequest represents a request for AI code analysis.
type AIAnalysisRequest struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language"`
}

// AIAnalysisResponse represents the AI analysis result.
type AIAnalysisResponse struct {
	Complexity   int                    `json:"complexity"`
	Issues       []agents.CodeIssue     `json:"issues"`
	Patterns     []string               `json:"patterns"`
	Strengths    []string               `json:"strengths"`
	Suggestions  []map[string]interface{} `json:"suggestions"`
	FunctionInfo []agents.FunctionInfo  `json:"functions"`
	Variables    []agents.VariableDecl  `json:"variables"`
}

// CodeCompletionRequest represents a code completion request.
type CodeCompletionRequest struct {
	Code     string `json:"code" validate:"required"`
	Position int    `json:"position"` // Cursor position
}

// CodeCompletionResponse represents code completion suggestions.
type CodeCompletionResponse struct {
	Completions []map[string]interface{} `json:"completions"`
}

// ErrorExplanationRequest represents an error explanation request.
type ErrorExplanationRequest struct {
	Code  string `json:"code" validate:"required"`
	Error string `json:"error" validate:"required"`
}

// ErrorExplanationResponse represents error explanations.
type ErrorExplanationResponse struct {
	Explanations []map[string]interface{} `json:"explanations"`
}

// TestGenerationRequest represents a test generation request.
type TestGenerationRequest struct {
	Code string `json:"code" validate:"required"`
}

// TestGenerationResponse represents generated test cases.
type TestGenerationResponse struct {
	TestCases []map[string]interface{} `json:"test_cases"`
}

// CodeExecutionRequest represents a code execution request with AI features.
type CodeExecutionRequest struct {
	Code      string `json:"code" validate:"required"`
	Language  string `json:"language"`
	SessionID string `json:"session_id"`
}

// CodeExecutionResponse represents code execution result with AI analysis.
type CodeExecutionResponse struct {
	Output          string                   `json:"output"`
	Error           string                   `json:"error,omitempty"`
	ExecutionTimeMs int64                    `json:"execution_time_ms"`
	Success         bool                     `json:"success"`
	AIAnalysis      *AIAnalysisResponse      `json:"ai_analysis,omitempty"`
	TestResults     []service.TestResult     `json:"test_results,omitempty"`
	DebugInfo       *DebugInfo               `json:"debug_info,omitempty"`
}

// DebugInfo contains debugging information.
type DebugInfo struct {
	Variables  []agents.VariableInfo `json:"variables,omitempty"`
	CallStack  []agents.StackFrame   `json:"call_stack,omitempty"`
	Breakpoints []agents.BreakpointInfo `json:"breakpoints,omitempty"`
}

// SessionResponse represents session data.
type SessionResponse struct {
	SessionID string                   `json:"session_id"`
	History   []agents.CodeHistoryEntry `json:"history,omitempty"`
}

// HandleAnalyze handles POST /api/v1/playground/analyze
func (h *PlaygroundAIHandler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	var req AIAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid request body"))
		return
	}

	if req.Code == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code is required"))
		return
	}

	if req.Language == "" {
		req.Language = "go"
	}

	// Rate limit check
	clientID := getClientIP(r)
	if !h.rateLimiter.allow(clientID, 30) { // 30 requests per minute
		h.writeErrorResponse(w, r, apierrors.NewRateLimitError("rate limit exceeded"))
		return
	}

	// Create agent request
	agentReq := agents.AgentRequest{
		Type:      agents.AgentTypeAIAnalysis,
		SessionID: generateSessionID(),
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"code":     req.Code,
			"language": req.Language,
		},
		Context:  make(map[string]interface{}),
		Priority: 1,
	}

	// Process through AI agent
	ctx := r.Context()
	response, err := h.agentPool.ProcessRequest(ctx, agentReq)
	if err != nil {
		h.logger.Error(ctx, "AI analysis failed", "error", err)
		h.writeErrorResponse(w, r, apierrors.NewInternalError("AI analysis failed", err))
		return
	}

	// Extract analysis result
	analysisData, ok := response.Data["analysis"].(*agents.CodeAnalysisResult)
	var result AIAnalysisResponse
	if ok {
		result = AIAnalysisResponse{
			Complexity:   analysisData.Complexity,
			Issues:       analysisData.Issues,
			Patterns:     analysisData.Patterns,
			Strengths:    analysisData.Strengths,
			FunctionInfo: analysisData.Functions,
			Variables:    analysisData.Variables,
		}
	}

	// Get suggestions
	if suggestions, ok := response.Data["suggestions"].([]map[string]interface{}); ok {
		result.Suggestions = suggestions
	}

	h.writeSuccessResponse(w, r, result, "code analysis completed")
}

// HandleComplete handles POST /api/v1/playground/complete
func (h *PlaygroundAIHandler) HandleComplete(w http.ResponseWriter, r *http.Request) {
	var req CodeCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid request body"))
		return
	}

	if req.Code == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code is required"))
		return
	}

	// Rate limit check
	clientID := getClientIP(r)
	if !h.rateLimiter.allow(clientID, 60) { // 60 requests per minute for completions
		h.writeErrorResponse(w, r, apierrors.NewRateLimitError("rate limit exceeded"))
		return
	}

	// Create agent request
	agentReq := agents.AgentRequest{
		Type:      agents.AgentTypeAIAnalysis,
		SessionID: generateSessionID(),
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"code":     req.Code,
			"position": req.Position,
		},
		Context:  make(map[string]interface{}),
		Priority: 2,
	}

	// Process through AI agent
	ctx := r.Context()
	response, err := h.agentPool.ProcessRequest(ctx, agentReq)
	if err != nil {
		h.logger.Error(ctx, "Code completion failed", "error", err)
		h.writeErrorResponse(w, r, apierrors.NewInternalError("code completion failed", err))
		return
	}

	result := CodeCompletionResponse{}
	if completions, ok := response.Data["completion"].([]map[string]interface{}); ok {
		result.Completions = completions
	}

	h.writeSuccessResponse(w, r, result, "completions generated")
}

// HandleExplain handles POST /api/v1/playground/explain
func (h *PlaygroundAIHandler) HandleExplain(w http.ResponseWriter, r *http.Request) {
	var req ErrorExplanationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid request body"))
		return
	}

	if req.Code == "" || req.Error == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code and error are required"))
		return
	}

	// Rate limit check
	clientID := getClientIP(r)
	if !h.rateLimiter.allow(clientID, 20) {
		h.writeErrorResponse(w, r, apierrors.NewRateLimitError("rate limit exceeded"))
		return
	}

	// Create agent request
	agentReq := agents.AgentRequest{
		Type:      agents.AgentTypeAIAnalysis,
		SessionID: generateSessionID(),
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"code":  req.Code,
			"error": req.Error,
		},
		Context:  make(map[string]interface{}),
		Priority: 1,
	}

	// Process through AI agent
	ctx := r.Context()
	response, err := h.agentPool.ProcessRequest(ctx, agentReq)
	if err != nil {
		h.logger.Error(ctx, "Error explanation failed", "error", err)
		h.writeErrorResponse(w, r, apierrors.NewInternalError("error explanation failed", err))
		return
	}

	result := ErrorExplanationResponse{}
	if explanations, ok := response.Data["explanations"].([]map[string]interface{}); ok {
		result.Explanations = explanations
	}

	h.writeSuccessResponse(w, r, result, "error explained")
}

// HandleGenerateTests handles POST /api/v1/playground/generate-tests
func (h *PlaygroundAIHandler) HandleGenerateTests(w http.ResponseWriter, r *http.Request) {
	var req TestGenerationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid request body"))
		return
	}

	if req.Code == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code is required"))
		return
	}

	// Rate limit check
	clientID := getClientIP(r)
	if !h.rateLimiter.allow(clientID, 15) {
		h.writeErrorResponse(w, r, apierrors.NewRateLimitError("rate limit exceeded"))
		return
	}

	// Create agent request for AI analysis (test generation)
	agentReq := agents.AgentRequest{
		Type:      agents.AgentTypeAIAnalysis,
		SessionID: generateSessionID(),
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"code": req.Code,
		},
		Context:  make(map[string]interface{}),
		Priority: 1,
	}

	// Process through AI agent
	ctx := r.Context()
	response, err := h.agentPool.ProcessRequest(ctx, agentReq)
	if err != nil {
		h.logger.Error(ctx, "Test generation failed", "error", err)
		h.writeErrorResponse(w, r, apierrors.NewInternalError("test generation failed", err))
		return
	}

	result := TestGenerationResponse{}
	if testCases, ok := response.Data["test_generation"].([]map[string]interface{}); ok {
		result.TestCases = testCases
	}

	h.writeSuccessResponse(w, r, result, "test cases generated")
}

// HandleExecuteWithAI handles POST /api/v1/playground/execute-ai
// Combines code execution with AI analysis in a single request
func (h *PlaygroundAIHandler) HandleExecuteWithAI(w http.ResponseWriter, r *http.Request) {
	var req CodeExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid request body"))
		return
	}

	if req.Code == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("code is required"))
		return
	}

	if req.Language == "" {
		req.Language = "go"
	}

	// Rate limit check
	clientID := getClientIP(r)
	if !h.rateLimiter.allow(clientID, 20) {
		h.writeErrorResponse(w, r, apierrors.NewRateLimitError("rate limit exceeded"))
		return
	}

	ctx := r.Context()
	startTime := time.Now()

	// Execute code
	execReq := &service.ExecuteRequest{
		Code:     req.Code,
		Language: req.Language,
		Timeout:  5 * time.Second,
		TestCases: []service.TestCase{
			{
				Name:     "Playground Output",
				Input:    "",
				Expected: "",
			},
		},
	}

	execResult, err := h.executor.ExecuteCode(ctx, execReq)
	executionTime := time.Since(startTime).Milliseconds()

	result := CodeExecutionResponse{
		ExecutionTimeMs: executionTime,
	}

	if err != nil {
		result.Error = err.Error()
		result.Success = false
	} else {
		// Build output from results (same pattern as playground.go)
		output := ""
		var execErr string
		if execResult.Error != nil {
			execErr = fmt.Sprintf("Error: %v", execResult.Error)
		} else if len(execResult.Results) > 0 {
			output = execResult.Results[0].Actual
			if execResult.Results[0].Error != "" {
				execErr = execResult.Results[0].Error
			}
		}
		result.Output = output
		result.Error = execErr
		result.Success = execResult.Error == nil && execErr == ""
		result.TestResults = execResult.Results
	}

	// Run AI analysis in parallel with execution context
	agentReq := agents.AgentRequest{
		Type:      agents.AgentTypeAIAnalysis,
		SessionID: req.SessionID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"code":  req.Code,
			"error": result.Error,
		},
		Context: map[string]interface{}{
			"execution": map[string]interface{}{
				"output":   result.Output,
				"success":  result.Success,
				"duration": executionTime,
			},
		},
		Priority: 1,
	}

	aiResponse, aiErr := h.agentPool.ProcessRequest(ctx, agentReq)
	if aiErr == nil && aiResponse.Success {
		if analysisData, ok := aiResponse.Data["analysis"].(*agents.CodeAnalysisResult); ok {
			result.AIAnalysis = &AIAnalysisResponse{
				Complexity:   analysisData.Complexity,
				Issues:       analysisData.Issues,
				Patterns:     analysisData.Patterns,
				Strengths:    analysisData.Strengths,
				FunctionInfo: analysisData.Functions,
				Variables:    analysisData.Variables,
			}
			if suggestions, ok := aiResponse.Data["suggestions"].([]map[string]interface{}); ok {
				result.AIAnalysis.Suggestions = suggestions
			}
		}
	}

	// Save to session history if session ID provided
	if req.SessionID != "" {
		stateReq := agents.AgentRequest{
			Type:      agents.AgentTypeStateManager,
			SessionID: req.SessionID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"action":    "save_code",
				"code":      req.Code,
				"language":  req.Language,
				"output":    result.Output,
				"session_id": req.SessionID,
			},
		}
		h.agentPool.ProcessRequest(ctx, stateReq)
	}

	h.writeSuccessResponse(w, r, result, "code executed with AI analysis")
}

// HandleCreateSession handles POST /api/v1/playground/sessions
func (h *PlaygroundAIHandler) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	// Create agent request for session creation
	agentReq := agents.AgentRequest{
		Type:      agents.AgentTypeStateManager,
		SessionID: "",
		UserID:    getClientIP(r), // Use IP as temporary user ID
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action": "create_session",
		},
	}

	ctx := r.Context()
	response, err := h.agentPool.ProcessRequest(ctx, agentReq)
	if err != nil {
		h.writeErrorResponse(w, r, apierrors.NewInternalError("failed to create session", err))
		return
	}

	sessionID, _ := response.Data["session_id"].(string)
	result := SessionResponse{
		SessionID: sessionID,
	}

	h.writeSuccessResponse(w, r, result, "session created")
}

// HandleGetSession handles GET /api/v1/playground/sessions/{id}
func (h *PlaygroundAIHandler) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("session ID is required"))
		return
	}

	// Get session info
	sessionReq := agents.AgentRequest{
		Type:      agents.AgentTypeStateManager,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":     "get_session",
			"session_id": sessionID,
		},
	}

	ctx := r.Context()
	sessionResp, err := h.agentPool.ProcessRequest(ctx, sessionReq)
	if err != nil {
		h.writeErrorResponse(w, r, apierrors.NewNotFoundError("session not found"))
		return
	}

	// Get history
	historyReq := agents.AgentRequest{
		Type:      agents.AgentTypeStateManager,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":     "get_history",
			"session_id": sessionID,
			"limit":      20,
		},
	}

	result := SessionResponse{
		SessionID: sessionID,
	}

	historyResp, _ := h.agentPool.ProcessRequest(ctx, historyReq)
	if historyResp.Success {
		if history, ok := historyResp.Data["history"].([]agents.CodeHistoryEntry); ok {
			result.History = history
		}
	}

	// Add session data
	if sessionData, ok := sessionResp.Data["session"]; ok {
		_ = sessionData // Session data available if needed
	}

	h.writeSuccessResponse(w, r, result, "session retrieved")
}

// HandleGetHistory handles GET /api/v1/playground/sessions/{id}/history
func (h *PlaygroundAIHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("session ID is required"))
		return
	}

	limit := 10 // default

	// Create agent request
	agentReq := agents.AgentRequest{
		Type:      agents.AgentTypeStateManager,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":     "get_history",
			"session_id": sessionID,
			"limit":      limit,
		},
	}

	ctx := r.Context()
	response, err := h.agentPool.ProcessRequest(ctx, agentReq)
	if err != nil {
		h.writeErrorResponse(w, r, apierrors.NewNotFoundError("session not found"))
		return
	}

	result := struct {
		History []agents.CodeHistoryEntry `json:"history"`
	}{}

	if history, ok := response.Data["history"].([]agents.CodeHistoryEntry); ok {
		result.History = history
	}

	h.writeSuccessResponse(w, r, result, "history retrieved")
}

// Helper methods

func (h *PlaygroundAIHandler) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success":   true,
		"data":      data,
		"message":   message,
		"timestamp": time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), "Failed to encode response", "error", err)
	}
}

func (h *PlaygroundAIHandler) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apierrors.APIError
	var statusCode int

	if e, ok := err.(*apierrors.APIError); ok {
		apiErr = e
		statusCode = e.StatusCode
	} else {
		statusCode = http.StatusInternalServerError
		apiErr = apierrors.NewInternalError("internal server error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": apiErr.Message,
			"type":    apiErr.Type,
		},
		"timestamp": time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), "Failed to encode error response", "error", err)
	}
}

func generateSessionID() string {
	return "session_" + time.Now().Format("20060102150405")
}
