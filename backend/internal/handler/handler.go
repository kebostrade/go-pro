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
	"time"

	"go-pro-backend/internal/domain"
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
}

// New creates a new HTTP handler.
func New(services *service.Services, logger logger.Logger, validator validator.Validator) *Handler {
	return &Handler{
		services:  services,
		logger:    logger,
		validator: validator,
	}
}

// RegisterRoutes registers all API routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
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

	// Progress tracking.
	mux.HandleFunc("GET /api/v1/progress/{userId}", h.handleGetProgress)
	mux.HandleFunc("POST /api/v1/progress/{userId}/lesson/{lessonId}", h.handleUpdateProgress)

	// Curriculum.
	mux.HandleFunc("GET /api/v1/curriculum", h.handleGetCurriculum)
	mux.HandleFunc("GET /api/v1/curriculum/lesson/{id}", h.handleGetLessonDetail)

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

	var req domain.SubmitExerciseRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
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

// Curriculum handlers.
func (h *Handler) handleGetCurriculum(w http.ResponseWriter, r *http.Request) {
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

	lesson, err := h.services.Curriculum.GetLessonDetail(r.Context(), lessonID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

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
