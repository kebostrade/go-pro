package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllOrders(t *testing.T) {
	orders := GetAllOrders()
	if len(orders) != 3 {
		t.Errorf("expected 3 orders, got %d", len(orders))
	}
}

func TestGetOrderByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"existing order", "1", false},
		{"existing order 2", "2", false},
		{"existing order 3", "3", false},
		{"nonexistent order", "999", true},
		{"empty id", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetOrderByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetOrdersByUserID(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		wantCount int
	}{
		{"user with orders", "1", 2},
		{"user with one order", "2", 1},
		{"user with no orders", "999", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := GetOrdersByUserID(tt.userID)
			if len(orders) != tt.wantCount {
				t.Errorf("GetOrdersByUserID() got %d orders, want %d", len(orders), tt.wantCount)
			}
		})
	}
}

func TestGetOrderCount(t *testing.T) {
	count := GetOrderCount()
	if count != 3 {
		t.Errorf("expected 3 orders, got %d", count)
	}
}

func TestHandleGetOrders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	w := httptest.NewRecorder()

	HandleGetOrders(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	orders, ok := resp["orders"].([]interface{})
	if !ok {
		t.Fatal("orders field not found or not array")
	}
	if len(orders) != 3 {
		t.Errorf("expected 3 orders, got %d", len(orders))
	}

	count, ok := resp["count"].(float64)
	if !ok {
		t.Fatal("count field not found or not number")
	}
	if int(count) != 3 {
		t.Errorf("expected count 3, got %d", int(count))
	}
}

func TestHandleGetOrderByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{"existing order", "1", http.StatusOK},
		{"nonexistent order", "999", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/orders/"+tt.id, nil)
			w := httptest.NewRecorder()

			HandleGetOrderByID(w, req, tt.id)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestParsePortServiceB(t *testing.T) {
	tests := []struct {
		name     string
		envPort  string
		wantPort int
	}{
		{"default when empty", "", 8002},
		{"valid port", "9000", 9000},
		{"invalid port string", "invalid", 8002},
		{"port out of range low", "0", 8002},
		{"port out of range high", "70000", 8002},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVICE_PORT", tt.envPort)
			got := ParsePort()
			if got != tt.wantPort {
				t.Errorf("ParsePort() = %d, want %d", got, tt.wantPort)
			}
		})
	}
}

func TestOrderJSON(t *testing.T) {
	order := Order{ID: "1", UserID: "1", Product: "Widget A", Amount: 100, Status: "shipped"}
	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal order: %v", err)
	}

	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal order: %v", err)
	}

	if decoded.ID != order.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, order.ID)
	}
	if decoded.Product != order.Product {
		t.Errorf("Product mismatch: got %s, want %s", decoded.Product, order.Product)
	}
	if decoded.Amount != order.Amount {
		t.Errorf("Amount mismatch: got %d, want %d", decoded.Amount, order.Amount)
	}
	if decoded.Status != order.Status {
		t.Errorf("Status mismatch: got %s, want %s", decoded.Status, order.Status)
	}
}
