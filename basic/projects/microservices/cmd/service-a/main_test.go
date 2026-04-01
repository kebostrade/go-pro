package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	users := GetAllUsers()
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"existing user", "1", false},
		{"existing user 2", "2", false},
		{"existing user 3", "3", false},
		{"nonexistent user", "999", true},
		{"empty id", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetUserByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserCount(t *testing.T) {
	count := GetUserCount()
	if count != 3 {
		t.Errorf("expected 3 users, got %d", count)
	}
}

func TestHandleGetUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	HandleGetUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	users, ok := resp["users"].([]interface{})
	if !ok {
		t.Fatal("users field not found or not array")
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}

	count, ok := resp["count"].(float64)
	if !ok {
		t.Fatal("count field not found or not number")
	}
	if int(count) != 3 {
		t.Errorf("expected count 3, got %d", int(count))
	}
}

func TestHandleGetUserByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{"existing user", "1", http.StatusOK},
		{"nonexistent user", "999", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/users/"+tt.id, nil)
			w := httptest.NewRecorder()

			HandleGetUserByID(w, req, tt.id)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		name     string
		envPort  string
		wantPort int
	}{
		{"default when empty", "", 8001},
		{"valid port", "9000", 9000},
		{"invalid port string", "invalid", 8001},
		{"port out of range low", "0", 8001},
		{"port out of range high", "70000", 8001},
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

func TestUserJSON(t *testing.T) {
	user := User{ID: "1", Name: "Test", Email: "test@example.com", Age: 25}
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal user: %v", err)
	}

	if decoded.ID != user.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, user.ID)
	}
	if decoded.Name != user.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, user.Name)
	}
	if decoded.Email != user.Email {
		t.Errorf("Email mismatch: got %s, want %s", decoded.Email, user.Email)
	}
	if decoded.Age != user.Age {
		t.Errorf("Age mismatch: got %d, want %d", decoded.Age, user.Age)
	}
}
