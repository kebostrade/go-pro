package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Repository handles data access operations for orders
type Repository struct {
	db              *sql.DB
	userServiceURL  string
	log             *logrus.Logger
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB, userServiceURL string, log *logrus.Logger) *Repository {
	return &Repository{
		db:             db,
		userServiceURL: userServiceURL,
		log:            log,
	}
}

// Initialize creates the orders table if it doesn't exist
func (r *Repository) Initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		items JSONB NOT NULL,
		total DECIMAL(10, 2) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	r.log.Info("Database initialized successfully")
	return nil
}

// Create inserts a new order into the database
func (r *Repository) Create(ctx context.Context, order *Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	query := `
	INSERT INTO orders (id, user_id, items, total, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.ExecContext(ctx, query,
		order.ID, order.UserID, itemsJSON, order.Total, order.Status,
		order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	r.log.WithFields(logrus.Fields{
		"order_id": order.ID,
		"user_id":  order.UserID,
		"total":    order.Total,
	}).Info("Order created successfully")

	return nil
}

// FindByID retrieves an order by ID
func (r *Repository) FindByID(ctx context.Context, id string) (*Order, error) {
	query := `
	SELECT id, user_id, items, total, status, created_at, updated_at
	FROM orders
	WHERE id = $1
	`

	var order Order
	var itemsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&itemsJSON,
		&order.Total,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	r.log.WithField("order_id", id).Debug("Order retrieved from database")
	return &order, nil
}

// FindAll retrieves all orders
func (r *Repository) FindAll(ctx context.Context) ([]Order, error) {
	query := `
	SELECT id, user_id, items, total, status, created_at, updated_at
	FROM orders
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		var itemsJSON []byte

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&itemsJSON,
			&order.Total,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		orders = append(orders, order)
	}

	r.log.WithField("count", len(orders)).Debug("Orders retrieved from database")
	return orders, nil
}

// FindByUserID retrieves all orders for a specific user
func (r *Repository) FindByUserID(ctx context.Context, userID string) ([]Order, error) {
	query := `
	SELECT id, user_id, items, total, status, created_at, updated_at
	FROM orders
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders by user ID: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		var itemsJSON []byte

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&itemsJSON,
			&order.Total,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		orders = append(orders, order)
	}

	r.log.WithFields(logrus.Fields{
		"user_id": userID,
		"count":   len(orders),
	}).Debug("Orders retrieved by user ID")

	return orders, nil
}

// UpdateStatus updates the status of an order
func (r *Repository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
	UPDATE orders
	SET status = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	r.log.WithFields(logrus.Fields{
		"order_id": id,
		"status":   status,
	}).Info("Order status updated")

	return nil
}

// Delete removes an order from the database
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	r.log.WithField("order_id", id).Info("Order deleted successfully")
	return nil
}

// ValidateUser checks if a user exists via User Service
func (r *Repository) ValidateUser(ctx context.Context, userID string) error {
	// In a real implementation, this would make an HTTP request to the User Service
	// For simplicity, we'll just log the validation attempt
	r.log.WithField("user_id", userID).Info("User validation attempt")

	// TODO: Implement actual HTTP call to User Service
	// Example:
	// resp, err := http.Get(r.userServiceURL + "/api/users/" + userID)
	// if err != nil || resp.StatusCode != http.StatusOK {
	//     return fmt.Errorf("user not found")
	// }

	return nil
}

// CheckHealth verifies database connectivity
func (r *Repository) CheckHealth(ctx context.Context) string {
	if err := r.db.PingContext(ctx); err != nil {
		r.log.WithError(err).Error("Database health check failed")
		return "unhealthy"
	}
	return "healthy"
}
