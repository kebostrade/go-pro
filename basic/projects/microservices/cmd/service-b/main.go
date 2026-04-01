// Package main provides the Order Service (service-b) entry point.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Order represents an order entity.
type Order struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Product   string `json:"product"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// User represents a user entity (for validation).
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory order store with seed data.
var orders = map[string]Order{
	"1": {ID: "1", UserID: "1", Product: "Widget A", Amount: 100, Status: "shipped", CreatedAt: "2024-01-15T10:00:00Z"},
	"2": {ID: "2", UserID: "2", Product: "Widget B", Amount: 250, Status: "processing", CreatedAt: "2024-01-16T14:30:00Z"},
	"3": {ID: "3", UserID: "1", Product: "Widget C", Amount: 75, Status: "delivered", CreatedAt: "2024-01-17T09:15:00Z"},
}

func main() {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://service-a:8001"
	}

	// Health check endpoint.
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":       "ok",
			"service":      "service-b",
			"port":         port,
			"user_service": userServiceURL,
		})
	})

	// GET /api/orders - List all orders.
	http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		orderList := make([]Order, 0, len(orders))
		for _, o := range orders {
			orderList = append(orderList, o)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"orders": orderList,
			"count":  len(orderList),
		})
	})

	// GET /api/orders/{id} - Get order by ID.
	http.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Path[len("/api/orders/"):]
		order, ok := orders[id]
		if !ok {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	})

	// GET /api/orders/user/{user_id} - Get orders by user ID (calls user service for validation).
	http.HandleFunc("/api/orders/user/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Path[len("/api/orders/user/"):]

		// Validate user exists by calling user service.
		resp, err := http.Get(userServiceURL + "/api/users/" + userID)
		if err != nil {
			http.Error(w, "Failed to validate user", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Get user's orders.
		var userOrders []Order
		for _, o := range orders {
			if o.UserID == userID {
				userOrders = append(userOrders, o)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": userID,
			"orders":  userOrders,
			"count":   len(userOrders),
		})
	})

	addr := ":" + port
	log.Printf("🚀 Order Service (service-b) starting on %s", addr)
	log.Printf("📋 Health: http://localhost%s/health", addr)
	log.Printf("📋 Orders: http://localhost%s/api/orders", addr)
	log.Printf("🔗 User Service URL: %s", userServiceURL)

	// Graceful shutdown with timeout.
	srv := &http.Server{
		Addr:         addr,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

// Handlers for testing.
func HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	orderList := make([]Order, 0, len(orders))
	for _, o := range orders {
		orderList = append(orderList, o)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orders": orderList,
		"count":  len(orderList),
	})
}

func HandleGetOrderByID(w http.ResponseWriter, r *http.Request, id string) {
	order, ok := orders[id]
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetOrderByID is exported for testing.
func GetOrderByID(id string) (Order, error) {
	order, ok := orders[id]
	if !ok {
		return Order{}, fmt.Errorf("order not found")
	}
	return order, nil
}

// GetOrdersByUserID returns orders for a specific user.
func GetOrdersByUserID(userID string) []Order {
	var result []Order
	for _, o := range orders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}
	return result
}

// GetAllOrders returns all orders.
func GetAllOrders() []Order {
	orderList := make([]Order, 0, len(orders))
	for _, o := range orders {
		orderList = append(orderList, o)
	}
	return orderList
}

// GetOrderCount returns the number of orders.
func GetOrderCount() int {
	return len(orders)
}

// ValidateUserExists checks if a user exists via the user service.
func ValidateUserExists(userID string, userServiceURL string) (bool, error) {
	resp, err := http.Get(userServiceURL + "/api/users/" + userID)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	return resp.StatusCode == http.StatusOK, nil
}

// ParsePort parses and validates the port from environment.
func ParsePort() int {
	portStr := os.Getenv("SERVICE_PORT")
	if portStr == "" {
		portStr = "8002"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return 8002
	}
	return port
}
