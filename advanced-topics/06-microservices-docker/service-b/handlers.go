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

// CreateOrder handles POST /api/orders
func (h *Handlers) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Validate user exists
	if err := h.repo.ValidateUser(c.Request.Context(), req.UserID); err != nil {
		h.log.WithError(err).WithField("user_id", req.UserID).Error("User validation failed")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user",
			Message: "User not found",
		})
		return
	}

	order := &Order{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Items:     req.Items,
		Total:     req.Total,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.Create(c.Request.Context(), order); err != nil {
		h.log.WithError(err).Error("Failed to create order")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	h.log.WithFields(logrus.Fields{
		"order_id": order.ID,
		"user_id":  order.UserID,
	}).Info("Order created")

	c.JSON(http.StatusCreated, order)
}

// GetOrder handles GET /api/orders/:id
func (h *Handlers) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.log.WithError(err).WithField("order_id", id).Error("Order not found")
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrders handles GET /api/orders
func (h *Handlers) GetOrders(c *gin.Context) {
	orders, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve orders")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "query_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":  len(orders),
		"orders": orders,
	})
}

// GetUserOrders handles GET /api/orders/user/:user_id
func (h *Handlers) GetUserOrders(c *gin.Context) {
	userID := c.Param("user_id")

	orders, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		h.log.WithError(err).WithField("user_id", userID).Error("Failed to retrieve user orders")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "query_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"count":   len(orders),
		"orders":  orders,
	})
}

// UpdateOrderStatus handles PUT /api/orders/:id/status
func (h *Handlers) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	if err := h.repo.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		h.log.WithError(err).Error("Failed to update order status")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	h.log.WithFields(logrus.Fields{
		"order_id": id,
		"status":   req.Status,
	}).Info("Order status updated")

	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"status": req.Status,
	})
}

// DeleteOrder handles DELETE /api/orders/:id
func (h *Handlers) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		h.log.WithError(err).WithField("order_id", id).Error("Failed to delete order")
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Order not found",
		})
		return
	}

	h.log.WithField("order_id", id).Info("Order deleted")
	c.Status(http.StatusNoContent)
}

// Health handles GET /health
func (h *Handlers) Health(c *gin.Context) {
	dbStatus := h.repo.CheckHealth(c.Request.Context())

	status := "healthy"
	if dbStatus != "healthy" {
		status = "unhealthy"
	}

	response := HealthResponse{
		Status:    status,
		Service:   "order-service",
		Timestamp: time.Now().Format(time.RFC3339),
		DB:        dbStatus,
	}

	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
