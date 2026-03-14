package main

import "time"

// Order represents an order in the system
type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Items     []Item    `json:"items"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Item represents an item in an order
type Item struct {
	Product  string  `json:"product"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Items  []Item  `json:"items" binding:"required,min=1"`
	Total  float64 `json:"total" binding:"required,gt=0"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing completed cancelled"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
	DB        string `json:"db"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
