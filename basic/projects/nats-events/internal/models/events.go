package models

import "time"

// OrderEvent represents an order-related event
type OrderEvent struct {
	EventType string    `json:"event_type"` // "created", "updated", "deleted"
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// UserEvent represents a user-related event
type UserEvent struct {
	EventType string    `json:"event_type"` // "registered", "updated", "deleted"
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

// TaskEvent represents a task processing event
type TaskEvent struct {
	TaskID     string    `json:"task_id"`
	Type       string    `json:"type"` // "process", "email", "webhook"
	Payload    string    `json:"payload"`
	Priority   int       `json:"priority"`
	RetryCount int       `json:"retry_count"`
	Timestamp  time.Time `json:"timestamp"`
}
