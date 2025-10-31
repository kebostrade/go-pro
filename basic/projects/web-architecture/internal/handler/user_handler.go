package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/model"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/service"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/pkg/response"
	"github.com/go-chi/chi/v5"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.Register(r.Context(), &req)
	if err != nil {
		if err == service.ErrUserExists {
			response.Conflict(w, "User already exists")
			return
		}
		response.InternalServerError(w, "Failed to create user")
		return
	}

	response.Created(w, user)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	loginResp, err := h.service.Login(r.Context(), &req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Unauthorized(w, "Invalid credentials")
			return
		}
		response.InternalServerError(w, "Failed to login")
		return
	}

	response.Success(w, loginResp)
}

// GetByID retrieves a user by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, user)
}

// List retrieves a list of users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	users, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(w, "Failed to retrieve users")
		return
	}

	response.Success(w, users)
}

// Update updates a user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		response.InternalServerError(w, "Failed to update user")
		return
	}

	response.Success(w, user)
}

// Delete deletes a user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.InternalServerError(w, "Failed to delete user")
		return
	}

	response.NoContent(w)
}

