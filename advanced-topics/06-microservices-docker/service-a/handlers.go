package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	repo *Repository
	log  *logrus.Logger
}

// NewHandlers creates a new handlers instance
func NewHandlers(repo *Repository, log *logrus.Logger) *Handlers {
	return &Handlers{
		repo: repo,
		log:  log,
	}
}

// CreateUser handles POST /api/users
func (h *Handlers) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	user := &User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.Create(c.Request.Context(), user); err != nil {
		h.log.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	h.log.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User created")

	c.JSON(http.StatusCreated, user)
}

// GetUser handles GET /api/users/:id
func (h *Handlers) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.log.WithError(err).WithField("user_id", id).Error("User not found")
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers handles GET /api/users
func (h *Handlers) GetUsers(c *gin.Context) {
	users, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve users")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "query_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(users),
		"users": users,
	})
}

// UpdateUser handles PUT /api/users/:id
func (h *Handlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Check if user exists
	user, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.log.WithError(err).WithField("user_id", id).Error("User not found")
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "User not found",
		})
		return
	}

	// Update user fields
	user.Name = req.Name
	user.Email = req.Email
	user.UpdatedAt = time.Now()

	if err := h.repo.Update(c.Request.Context(), user); err != nil {
		h.log.WithError(err).Error("Failed to update user")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	h.log.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User updated")

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /api/users/:id
func (h *Handlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		h.log.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "User not found",
		})
		return
	}

	h.log.WithField("user_id", id).Info("User deleted")
	c.Status(http.StatusNoContent)
}

// Health handles GET /health
func (h *Handlers) Health(c *gin.Context) {
	dbStatus, redisStatus := h.repo.CheckHealth(c.Request.Context())

	status := "healthy"
	if dbStatus != "healthy" || redisStatus != "healthy" {
		status = "unhealthy"
	}

	response := HealthResponse{
		Status:    status,
		Service:   "user-service",
		Timestamp: time.Now().Format(time.RFC3339),
		DB:        dbStatus,
		Redis:     redisStatus,
	}

 statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
