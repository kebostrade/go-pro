// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"net/http"

	"go-pro-backend/internal/domain"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/middleware"
	"go-pro-backend/internal/service"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	services  *service.Services
	logger    logger.Logger
	validator validator.Validator
}

// NewAuthHandler creates a new authentication handler.
func NewAuthHandler(services *service.Services, logger logger.Logger, validator validator.Validator) *AuthHandler {
	return &AuthHandler{
		services:  services,
		logger:    logger,
		validator: validator,
	}
}

// RegisterRoutes registers all authentication routes.
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	// Public authentication endpoint (no auth required).
	mux.HandleFunc("POST /api/v1/auth/verify", h.handleVerifyToken)

	// Protected endpoints (require authentication).
	mux.Handle("GET /api/v1/auth/me", authMiddleware.AuthRequired(http.HandlerFunc(h.handleGetProfile)))
	mux.Handle("PUT /api/v1/auth/me", authMiddleware.AuthRequired(http.HandlerFunc(h.handleUpdateProfile)))
}

// handleVerifyToken verifies a Firebase ID token and syncs the user.
// POST /api/v1/auth/verify
func (h *AuthHandler) handleVerifyToken(w http.ResponseWriter, r *http.Request) {
	var req domain.VerifyTokenRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	response, err := h.services.Auth.VerifyAndSyncUser(r.Context(), req.IDToken)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	h.writeSuccessResponse(w, r, response, "token verified successfully")
}

// handleGetProfile retrieves the current user's profile.
// GET /api/v1/auth/me
func (h *AuthHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	profile := user.ToProfileResponse()
	h.writeSuccessResponse(w, r, profile, "profile retrieved successfully")
}

// handleUpdateProfile updates the current user's profile.
// PUT /api/v1/auth/me
func (h *AuthHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, r, apierrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	var req domain.UpdateUserRequest
	if err := h.validator.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	updatedUser, err := h.services.Auth.UpdateUserProfile(r.Context(), user.ID, &req)
	if err != nil {
		h.writeErrorResponse(w, r, err)
		return
	}

	profile := updatedUser.ToProfileResponse()
	h.writeSuccessResponse(w, r, profile, "profile updated successfully")
}

// Helper methods.

// writeSuccessResponse writes a successful API response.
func (h *AuthHandler) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	// Delegate to the shared handler function from handler.go.
	handler := &Handler{
		services:  h.services,
		logger:    h.logger,
		validator: h.validator,
	}
	handler.writeSuccessResponse(w, r, data, message)
}

// writeErrorResponse writes an error API response.
func (h *AuthHandler) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// Delegate to the shared handler function from handler.go.
	handler := &Handler{
		services:  h.services,
		logger:    h.logger,
		validator: h.validator,
	}
	handler.writeErrorResponse(w, r, err)
}
