package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/internal/service"
	apperrors "github.com/DimaJoyti/go-pro/basic/projects/rest-api/pkg/errors"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns the health status of the service
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UserHandler handles HTTP requests for users
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler creates a new UserHandler with the given service
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	user, err := h.svc.Create(r.Context(), req.Name, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// ListUsers handles GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser handles GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	user, err := h.svc.Update(r.Context(), id, req.Name, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, apperrors.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to write error responses
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
