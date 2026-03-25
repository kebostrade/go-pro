// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package handler provides HTTP request handlers for the API.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/middleware"
	"go-pro-backend/internal/service"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"

	apierrors "go-pro-backend/internal/errors"
)

// Handler represents the HTTP handler for the API.
type Handler struct {
	services  *service.Services
	logger    logger.Logger
	validator validator.Validator

	// Rate limiting for exercise submissions (per-user).
	submissionLimits map[string]*rateLimitState

	// AI-powered playground handler (optional, set after initialization).
	aiHandler *PlaygroundAIHandler

	// Interview handler for mock interviews.
	interviewHandler *InterviewHandler
}

// rateLimitState tracks submission rate limits per user.
type rateLimitState struct {
	count     int
	resetTime time.Time
}

// New creates a new HTTP handler.
func New(services *service.Services, logger logger.Logger, validator validator.Validator) *Handler {
	return &Handler{
		services:         services,
		logger:           logger,
		validator:        validator,
		submissionLimits: make(map[string]*rateLimitState),
	}
}

// SetAIHandler sets the AI-powered playground handler.
func (h *Handler) SetAIHandler(aiHandler *PlaygroundAIHandler) {
	h.aiHandler = aiHandler
}

// SetInterviewHandler sets the interview handler.
func (h *Handler) SetInterviewHandler(interviewHandler *InterviewHandler) {
	h.interviewHandler = interviewHandler
}

// RegisterRoutes registers all API routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	// Health check.
	mux.HandleFunc("GET /api/v1/health", h.handleHealth)

	// Course management.
	mux.HandleFunc("POST /api/v1/courses", h.handleCreateCourse)
	mux.HandleFunc("GET /api/v1/courses", h.handleGetCourses)
	mux.HandleFunc("GET /api/v1/courses/{id}", h.handleGetCourse)
	mux.HandleFunc("PUT /api/v1/courses/{id}", h.handleUpdateCourse)
	mux.HandleFunc("DELETE /api/v1/courses/{id}", h.handleDeleteCourse)

	// Course lessons.
	mux.HandleFunc("GET /api/v1/courses/{courseId}/lessons", h.handleGetCourseLessons)

	// Lesson management.
	mux.HandleFunc("POST /api/v1/lessons", h.handleCreateLesson)
	mux.HandleFunc("GET /api/v1/lessons/{id}", h.handleGetLesson)
	mux.HandleFunc("PUT /api/v1/lessons/{id}", h.handleUpdateLesson)
	mux.HandleFunc("DELETE /api/v1/lessons/{id}", h.handleDeleteLesson)

	// Exercise management.
	mux.HandleFunc("GET /api/v1/exercises/{id}", h.handleGetExercise)
	mux.HandleFunc("POST /api/v1/exercises/{id}/submit", h.handleSubmitExercise)

	// Progress tracking (REST-compliant).
	mux.HandleFunc("GET /api/v1/users/{userId}/progress", h.handleGetUserProgress)
	mux.HandleFunc("GET /api/v1/users/{userId}/progress/stats", h.handleGetProgressStats)
	mux.HandleFunc("POST /api/v1/users/{userId}/lessons/{lessonId}/progress", h.handleUpdateUserLessonProgress)

	// Legacy progress endpoints (kept for backward compatibility).
	// Note: Deprecated - use /api/v1/users/{userId}/progress instead
	mux.HandleFunc("GET /api/v1/progress/{userId}", h.handleGetProgress)
	mux.HandleFunc("POST /api/v1/progress/{userId}/lesson/{lessonId}", h.handleUpdateProgress)

	// Curriculum.
	mux.HandleFunc("GET /api/v1/curriculum", h.handleGetCurriculum)
	mux.HandleFunc("GET /api/v1/curriculum/lesson/{id}", h.handleGetLessonDetail)

	// Playground - code execution.
	mux.HandleFunc("POST /api/v1/playground/execute", h.handlePlaygroundExecute)

	// AI-powered playground endpoints (if AI handler is available).
	if h.aiHandler != nil {
		mux.HandleFunc("POST /api/v1/playground/analyze", h.aiHandler.HandleAnalyze)
		mux.HandleFunc("POST /api/v1/playground/complete", h.aiHandler.HandleComplete)
		mux.HandleFunc("POST /api/v1/playground/explain", h.aiHandler.HandleExplain)
		mux.HandleFunc("POST /api/v1/playground/generate-tests", h.aiHandler.HandleGenerateTests)
		mux.HandleFunc("POST /api/v1/playground/execute-ai", h.aiHandler.HandleExecuteWithAI)
		mux.HandleFunc("POST /api/v1/playground/sessions", h.aiHandler.HandleCreateSession)
		mux.HandleFunc("GET /api/v1/playground/sessions/{id}", h.aiHandler.HandleGetSession)
		mux.HandleFunc("GET /api/v1/playground/sessions/{id}/history", h.aiHandler.HandleGetHistory)
	}

	// Interview endpoints (if interview handler is available).
	if h.interviewHandler != nil {
		h.interviewHandler.RegisterRoutes(mux)
	}

	// API documentation.
	mux.HandleFunc("GET /", h.handleAPIDocumentation)
}

// Health endpoint.
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	healthStatus, err := h.services.Health.GetHealthStatus(r.Context())
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, healthStatus, "health check successful")
}

// Course handlers.
func (h *Handler) handleCreateCourse(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCourseRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	course, err := h.services.Course.CreateCourse(r.Context(), &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, course, "course created successfully")
}

func (h *Handler) handleGetCourses(w http.ResponseWriter, r *http.Request) {
	pagination := getPaginationFromContext(r.Context())

	response, err := h.services.Course.GetAllCourses(r.Context(), pagination)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, response, "courses retrieved successfully")
}

func (h *Handler) handleGetCourse(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("course ID is required"))
		return
	}

	course, err := h.services.Course.GetCourseByID(r.Context(), id)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, course, "course retrieved successfully")
}

func (h *Handler) handleUpdateCourse(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("course ID is required"))
		return
	}

	var req domain.UpdateCourseRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	course, err := h.services.Course.UpdateCourse(r.Context(), id, &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, course, "course updated successfully")
}

func (h *Handler) handleDeleteCourse(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("course ID is required"))
		return
	}

	if err := h.services.Course.DeleteCourse(r.Context(), id); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, nil, "course deleted successfully")
}

// Course lessons handler.
func (h *Handler) handleGetCourseLessons(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("courseId")
	if courseID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("course ID is required"))
		return
	}

	pagination := getPaginationFromContext(r.Context())

	response, err := h.services.Lesson.GetLessonsByCourseID(r.Context(), courseID, pagination)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, response, "course lessons retrieved successfully")
}

// Lesson handlers.
func (h *Handler) handleCreateLesson(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateLessonRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	lesson, err := h.services.Lesson.CreateLesson(r.Context(), &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, lesson, "lesson created successfully")
}

func (h *Handler) handleGetLesson(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("lesson ID is required"))
		return
	}

	lesson, err := h.services.Lesson.GetLessonByID(r.Context(), id)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, lesson, "lesson retrieved successfully")
}

func (h *Handler) handleUpdateLesson(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("lesson ID is required"))
		return
	}

	var req domain.UpdateLessonRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	lesson, err := h.services.Lesson.UpdateLesson(r.Context(), id, &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, lesson, "lesson updated successfully")
}

func (h *Handler) handleDeleteLesson(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("lesson ID is required"))
		return
	}

	if err := h.services.Lesson.DeleteLesson(r.Context(), id); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, nil, "lesson deleted successfully")
}

// Exercise handlers.
func (h *Handler) handleGetExercise(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("exercise ID is required"))
		return
	}

	exercise, err := h.services.Exercise.GetExerciseByID(r.Context(), id)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, exercise, "exercise retrieved successfully")
}

func (h *Handler) handleSubmitExercise(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("exercise ID is required"))
		return
	}

	// Check rate limit (10 submissions per minute per user).
	// Phase 1: Use IP address. Phase 2: Use authenticated user ID.
	userKey := getClientIP(r)
	if !h.checkSubmissionRateLimit(userKey) {
		h.writeErrorResponse(w, r, apierrors.NewRateLimitError("submission rate limit exceeded (10 per minute)"))
		return
	}

	var req domain.SubmitExerciseRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Validate language
	validLanguages := map[string]bool{
		"go":         true,
		"python":     true,
		"javascript": true,
	}
	if !validLanguages[req.Language] {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid language: must be one of: go, python, javascript"))
		return
	}

	result, err := h.services.Exercise.SubmitExercise(r.Context(), id, &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, result, "exercise submitted and evaluated")
}

// Progress handlers.
func (h *Handler) handleGetProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	pagination := getPaginationFromContext(r.Context())

	response, err := h.services.Progress.GetProgressByUserID(r.Context(), userID, pagination)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, response, "progress retrieved successfully")
}

func (h *Handler) handleUpdateProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	lessonID := r.PathValue("lessonId")

	if userID == "" || lessonID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID and lesson ID are required"))
		return
	}

	var req domain.UpdateProgressRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	progress, err := h.services.Progress.UpdateProgress(r.Context(), userID, lessonID, &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, progress, "progress updated successfully")
}

// New progress handlers (REST-compliant).

// handleGetUserProgress retrieves user's lesson progress with pagination.
// GET /api/v1/users/:userId/progress?page=1&pageSize=20
func (h *Handler) handleGetUserProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	pagination := getPaginationFromContext(r.Context())

	// Get progress records from service.
	response, err := h.services.Progress.GetProgressByUserID(r.Context(), userID, pagination)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, response, "user progress retrieved successfully")
}

// handleGetProgressStats retrieves user's progress statistics.
// GET /api/v1/users/:userId/progress/stats
func (h *Handler) handleGetProgressStats(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	// For now, calculate stats from progress records.
	// In future, this could be cached or pre-calculated.
	pagination := &domain.PaginationRequest{
		Page:     1,
		PageSize: 1000, // Get all progress for stats calculation
	}

	response, err := h.services.Progress.GetProgressByUserID(r.Context(), userID, pagination)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Calculate statistics.
	stats := calculateProgressStats(response.Items)

	h.writeSuccessResponse(w, r, stats, "progress statistics retrieved successfully")
}

// handleUpdateUserLessonProgress updates progress for a specific lesson.
// POST /api/v1/users/:userId/lessons/:lessonId/progress
func (h *Handler) handleUpdateUserLessonProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	lessonID := r.PathValue("lessonId")

	if userID == "" || lessonID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID and lesson ID are required"))
		return
	}

	var req domain.UpdateProgressRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	progress, err := h.services.Progress.UpdateProgress(r.Context(), userID, lessonID, &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, progress, "progress updated successfully")
}

// handleGetProgressByID retrieves a specific progress record by ID.
// GET /api/v1/progress/:id
func (h *Handler) handleGetProgressByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("progress ID is required"))
		return
	}

	progress, err := h.services.Progress.GetProgressByID(r.Context(), id)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, progress, "progress retrieved successfully")
}

// Curriculum handlers.
func (h *Handler) handleGetCurriculum(w http.ResponseWriter, r *http.Request) {
	// Performance optimization: Add HTTP caching headers
	// Cache-Control: public allows CDN and browser caching
	// max-age=3600 caches for 1 hour (curriculum changes infrequently)
	w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=7200")

	// Add Vary header to ensure proper caching behavior
	w.Header().Set("Vary", "Accept-Encoding")

	curriculum, err := h.services.Curriculum.GetCurriculum(r.Context())
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, curriculum, "curriculum retrieved successfully")
}

func (h *Handler) handleGetLessonDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("lesson ID is required"))
		return
	}

	// Convert string to int.
	var lessonID int
	if _, err := fmt.Sscanf(idStr, "%d", &lessonID); err != nil {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid lesson ID format"))
		return
	}

	// Performance optimization: ETag support for conditional requests
	// Generate ETag based on lesson ID and timestamp (simplified approach)
	etag := fmt.Sprintf(`"lesson-%d-%d"`, lessonID, time.Now().Unix()/(60*15)) // 15-minute granularity

	// Check If-None-Match header for cached version
	if match := r.Header.Get("If-None-Match"); match != "" {
		// Remove weak ETag prefix if present
		if len(match) > 2 && match[:2] == "W/" {
			match = match[2:]
		}
		if match == etag {
			// Content hasn't changed, return 304 Not Modified
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	lesson, err := h.services.Curriculum.GetLessonDetail(r.Context(), lessonID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Set ETag header for this response
	w.Header().Set("ETag", etag)
	// Cache for 5 minutes on client, must revalidate with server
	w.Header().Set("Cache-Control", "private, max-age=300, must-revalidate")

	h.writeSuccessResponse(w, r, lesson, "lesson detail retrieved successfully")
}

// API documentation handler.
func (h *Handler) handleAPIDocumentation(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>GO-PRO API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { margin: 20px 0; padding: 15px; border-left: 4px solid #007acc; background: #f5f5f5; }
        .method { font-weight: bold; color: #007acc; }
        code { background: #e8e8e8; padding: 2px 4px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🚀 GO-PRO Learning Platform API</h1>
    <p>Welcome to the GO-PRO API! This RESTful API powers the Go programming learning platform.</p>

    <h2>📋 Available Endpoints</h2>

    <div class="endpoint">
        <div class="method">GET</div>
        <strong>/api/v1/health</strong>
        <p>Health check endpoint to verify API status</p>
    </div>

    <div class="endpoint">
        <div class="method">GET</div>
        <strong>/api/v1/courses</strong>
        <p>Retrieve all available courses (supports pagination)</p>
    </div>

    <div class="endpoint">
        <div class="method">POST</div>
        <strong>/api/v1/courses</strong>
        <p>Create a new course</p>
    </div>

    <div class="endpoint">
        <div class="method">GET</div>
        <strong>/api/v1/courses/{id}</strong>
        <p>Get details of a specific course</p>
    </div>

    <div class="endpoint">
        <div class="method">PUT</div>
        <strong>/api/v1/courses/{id}</strong>
        <p>Update a specific course</p>
    </div>

    <div class="endpoint">
        <div class="method">DELETE</div>
        <strong>/api/v1/courses/{id}</strong>
        <p>Delete a specific course</p>
    </div>

    <div class="endpoint">
        <div class="method">GET</div>
        <strong>/api/v1/courses/{courseId}/lessons</strong>
        <p>Get all lessons for a specific course</p>
    </div>

    <div class="endpoint">
        <div class="method">GET</div>
        <strong>/api/v1/exercises/{id}</strong>
        <p>Get details of a specific exercise</p>
    </div>

    <div class="endpoint">
        <div class="method">POST</div>
        <strong>/api/v1/exercises/{id}/submit</strong>
        <p>Submit an exercise solution for evaluation</p>
    </div>

    <div class="endpoint">
        <div class="method">GET</div>
        <strong>/api/v1/progress/{userId}</strong>
        <p>Get learning progress for a specific user</p>
    </div>

    <div class="endpoint">
        <div class="method">POST</div>
        <strong>/api/v1/progress/{userId}/lesson/{lessonId}</strong>
        <p>Update progress for a specific lesson</p>
    </div>

    <h2>🧪 Try It Out</h2>
    <p>Test the API with these example requests:</p>
    <ul>
        <li><a href="/api/v1/health">Health Check</a></li>
        <li><a href="/api/v1/courses">All Courses</a></li>
        <li><a href="/api/v1/courses/go-pro">GO-PRO Course</a></li>
        <li><a href="/api/v1/progress/demo-user">Demo User Progress</a></li>
    </ul>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		h.logger.Error(r.Context(), "Failed to write HTML response", "error", err)
	}
}

// Helper methods.

// writeSuccessResponse writes a successful API response.
func (h *Handler) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	response := &domain.APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		RequestID: logger.GetRequestID(r.Context()),
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.LogError(h.logger, r.Context(), err, "failed to encode response")
	}
}

// writeErrorResponse writes an error API response.
func (h *Handler) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apierrors.APIError
	var statusCode int

	if errors.As(err, &apiErr) {
		statusCode = apiErr.StatusCode
	} else {
		statusCode = http.StatusInternalServerError
		apiErr = apierrors.NewInternalError("internal server error", err)
	}

	response := &domain.APIResponse{
		Success: false,
		Error: &domain.APIError{
			Type:    apiErr.Type,
			Message: apiErr.Message,
		},
		RequestID: logger.GetRequestID(r.Context()),
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.LogError(h.logger, r.Context(), err, "failed to encode error response")
	}

	// Log the error.
	logger.LogError(h.logger, r.Context(), apiErr, "HTTP request error",
		"status_code", statusCode,
		"error_type", apiErr.Type,
	)
}

// getPaginationFromContext retrieves pagination from request context.
func getPaginationFromContext(ctx context.Context) *domain.PaginationRequest {
	if pagination, ok := ctx.Value("pagination").(*domain.PaginationRequest); ok {
		return pagination
	}

	return &domain.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
}

// checkSubmissionRateLimit checks if submission rate limit is exceeded.
// Returns true if allowed, false if rate limit exceeded.
func (h *Handler) checkSubmissionRateLimit(userKey string) bool {
	const (
		maxSubmissions = 10
		windowDuration = time.Minute
	)

	now := time.Now()

	// Get or create rate limit state for user.
	state, exists := h.submissionLimits[userKey]
	if !exists || now.After(state.resetTime) {
		// Create new window.
		h.submissionLimits[userKey] = &rateLimitState{
			count:     1,
			resetTime: now.Add(windowDuration),
		}
		return true
	}

	// Check if limit exceeded.
	if state.count >= maxSubmissions {
		return false
	}

	// Increment count.
	state.count++
	return true
}

// getClientIP extracts client IP from request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies).
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one.
		if idx := len(xff); idx > 0 {
			if commaIdx := 0; commaIdx < idx {
				for i, c := range xff {
					if c == ',' {
						commaIdx = i
						break
					}
				}
				if commaIdx > 0 {
					return xff[:commaIdx]
				}
			}
			return xff
		}
	}

	// Check X-Real-IP header.
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr.
	// RemoteAddr format is "IP:port", extract just the IP.
	if idx := len(r.RemoteAddr); idx > 0 {
		for i := idx - 1; i >= 0; i-- {
			if r.RemoteAddr[i] == ':' {
				return r.RemoteAddr[:i]
			}
		}
	}

	return r.RemoteAddr
}

// calculateProgressStats calculates progress statistics from progress records.
func calculateProgressStats(items interface{}) map[string]interface{} {
	stats := map[string]interface{}{
		"total_lessons":       0,
		"completed_lessons":   0,
		"in_progress_lessons": 0,
		"average_score":       0.0,
		"total_time_spent":    0,
	}

	// Type assert items to progress array.
	progressList, ok := items.([]interface{})
	if !ok {
		return stats
	}

	totalLessons := len(progressList)
	completedCount := 0
	inProgressCount := 0
	totalScore := 0
	scoreCount := 0

	for _, item := range progressList {
		progress, ok := item.(*domain.Progress)
		if !ok {
			continue
		}

		switch progress.Status {
		case domain.StatusCompleted:
			completedCount++
			// For now, assume 100 score for completed lessons.
			// In future, this should come from actual exercise scores.
			totalScore += 100
			scoreCount++
		case domain.StatusInProgress:
			inProgressCount++
		}
	}

	stats["total_lessons"] = totalLessons
	stats["completed_lessons"] = completedCount
	stats["in_progress_lessons"] = inProgressCount

	// Calculate average score.
	if scoreCount > 0 {
		stats["average_score"] = float64(totalScore) / float64(scoreCount)
	}

	// Total time spent would need to be tracked separately.
	// For now, estimate based on completed lessons (30 minutes per lesson).
	stats["total_time_spent"] = completedCount * 30

	return stats
}

// getPaginationParams extracts pagination params from gin context.
func getPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}
