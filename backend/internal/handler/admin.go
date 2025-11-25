// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"net/http"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/middleware"
	"go-pro-backend/internal/service"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// AdminHandler handles admin-related HTTP requests.
type AdminHandler struct {
	services  *service.Services
	logger    logger.Logger
	validator validator.Validator
}

// NewAdminHandler creates a new admin handler.
func NewAdminHandler(services *service.Services, logger logger.Logger, validator validator.Validator) *AdminHandler {
	return &AdminHandler{
		services:  services,
		logger:    logger,
		validator: validator,
	}
}

// RegisterRoutes registers all admin routes.
// All routes require authentication AND admin role.
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	// All admin routes require both authentication and admin role.
	mux.Handle("GET /api/v1/admin/users",
		authMiddleware.AuthRequired(authMiddleware.AdminRequired(http.HandlerFunc(h.handleListUsers))))
	mux.Handle("GET /api/v1/admin/users/{id}",
		authMiddleware.AuthRequired(authMiddleware.AdminRequired(http.HandlerFunc(h.handleGetUser))))
	mux.Handle("PUT /api/v1/admin/users/{id}/role",
		authMiddleware.AuthRequired(authMiddleware.AdminRequired(http.HandlerFunc(h.handleUpdateUserRole))))
	mux.Handle("DELETE /api/v1/admin/users/{id}",
		authMiddleware.AuthRequired(authMiddleware.AdminRequired(http.HandlerFunc(h.handleDeleteUser))))
}

// handleListUsers lists all users with pagination (admin only).
// GET /api/v1/admin/users?page=1&page_size=20
func (h *AdminHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	// Get admin user from context (middleware ensures this is an admin).
	adminUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("admin user not authenticated"))
		return
	}

	h.logger.Info(r.Context(), "admin listing users", "admin_id", adminUser.ID)

	// Get pagination from context.
	pagination := getPaginationFromContext(r.Context())

	// Get all users from repository.
	users, total, err := h.services.User.GetAllUsers(r.Context(), pagination)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Build paginated response.
	totalPages := (int64(total) + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)
	response := &domain.ListResponse{
		Items: users,
		Pagination: &domain.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: int64(total),
			TotalPages: int(totalPages),
		},
	}

	h.writeSuccessResponse(w, r, response, "users retrieved successfully")
}

// handleGetUser retrieves a specific user by ID (admin only).
// GET /api/v1/admin/users/{id}
func (h *AdminHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Get admin user from context.
	adminUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("admin user not authenticated"))
		return
	}

	userID := r.PathValue("id")
	if userID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	h.logger.Info(r.Context(), "admin fetching user details", "admin_id", adminUser.ID, "target_user_id", userID)

	user, err := h.services.User.GetUserByID(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, user, "user retrieved successfully")
}

// handleUpdateUserRole updates a user's role (admin only).
// PUT /api/v1/admin/users/{id}/role
func (h *AdminHandler) handleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	// Get admin user from context.
	adminUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("admin user not authenticated"))
		return
	}

	targetUserID := r.PathValue("id")
	if targetUserID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	var req domain.UpdateUserRoleRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	// Validate role.
	if !req.Role.IsValid() {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("invalid role"))
		return
	}

	h.logger.Info(r.Context(), "admin updating user role",
		"admin_id", adminUser.ID,
		"target_user_id", targetUserID,
		"new_role", req.Role,
	)

	// Update user role.
	if err := h.services.Auth.UpdateUserRole(r.Context(), adminUser.ID, targetUserID, req.Role); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, map[string]interface{}{
		"user_id": targetUserID,
		"role":    req.Role,
	}, "user role updated successfully")
}

// handleDeleteUser deletes a user (admin only).
// DELETE /api/v1/admin/users/{id}
func (h *AdminHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get admin user from context.
	adminUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("admin user not authenticated"))
		return
	}

	targetUserID := r.PathValue("id")
	if targetUserID == "" {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("user ID is required"))
		return
	}

	// Prevent admin from deleting themselves.
	if adminUser.ID == targetUserID {
		h.writeErrorResponse(w, r, apierrors.NewBadRequestError("admins cannot delete themselves"))
		return
	}

	h.logger.Info(r.Context(), "admin deleting user",
		"admin_id", adminUser.ID,
		"target_user_id", targetUserID,
	)

	// Delete user.
	if err := h.services.User.DeleteUser(r.Context(), targetUserID); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, map[string]interface{}{
		"user_id": targetUserID,
		"deleted": true,
	}, "user deleted successfully")
}

// Helper methods.

// writeSuccessResponse writes a successful API response.
func (h *AdminHandler) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	// Delegate to the shared handler function from handler.go.
	handler := &Handler{
		services:  h.services,
		logger:    h.logger,
		validator: h.validator,
	}
	handler.writeSuccessResponse(w, r, data, message)
}

// writeErrorResponse writes an error API response.
func (h *AdminHandler) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// Delegate to the shared handler function from handler.go.
	handler := &Handler{
		services:  h.services,
		logger:    h.logger,
		validator: h.validator,
	}
	handler.writeErrorResponse(w, r, err)
}
