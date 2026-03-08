//go:build integration
// +build integration

package testingexamples

// This file contains integration test examples
// Run with: go test -tags=integration -v ./...

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// User type for testing
type User struct {
	ID    int
	Name  string
	Email string
}

// ========================================
// EXAMPLE 1: HTTP TESTING WITH HTTPTEST
// ========================================

// Simple HTTP handler
type UserHandler struct {
	users map[string]*User
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		users: make(map[string]*User),
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodDelete:
		h.deleteUser(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email parameter required", http.StatusBadRequest)
		return
	}

	user, exists := h.users[email]
	if !exists {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if user.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	h.users[user.Email] = &user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email parameter required", http.StatusBadRequest)
		return
	}

	delete(h.users, email)
	w.WriteHeader(http.StatusNoContent)
}

func TestUserHandler_GetUser(t *testing.T) {
	handler := NewUserHandler()

	// Add a test user
	handler.users["alice@example.com"] = &User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	tests := []struct {
		name           string
		email          string
		expectedStatus int
		expectedBody   string
		checkUser      func(*testing.T, *User)
	}{
		{
			name:           "existing user",
			email:          "alice@example.com",
			expectedStatus: http.StatusOK,
			checkUser: func(t *testing.T, user *User) {
				t.Helper()
				if user.Name != "Alice" {
					t.Errorf("expected name Alice, got %s", user.Name)
				}
			},
		},
		{
			name:           "user not found",
			email:          "nonexistent@example.com",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing email parameter",
			email:          "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			url := fmt.Sprintf("/user?email=%s", tt.email)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body if successful
			if tt.expectedStatus == http.StatusOK && tt.checkUser != nil {
				var user User
				if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				tt.checkUser(t, &user)
			}
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	handler := NewUserHandler()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "valid user creation",
			requestBody:    `{"id": 1, "name": "Bob", "email": "bob@example.com"}`,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				var user User
				if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if user.Name != "Bob" {
					t.Errorf("expected name Bob, got %s", user.Name)
				}
			},
		},
		{
			name:           "missing email",
			requestBody:    `{"id": 1, "name": "Charlie"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	handler := NewUserHandler()

	// Add a test user
	handler.users["alice@example.com"] = &User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	t.Run("delete existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/user?email=alice@example.com", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}

		// Verify user is deleted
		if _, exists := handler.users["alice@example.com"]; exists {
			t.Error("user still exists after deletion")
		}
	})
}

// ========================================
// EXAMPLE 2: HTTP CLIENT TESTING
// ========================================

// APIClient is a simple HTTP client wrapper
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *APIClient) GetUser(email string) (*User, error) {
	url := fmt.Sprintf("%s/user?email=%s", c.baseURL, email)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *APIClient) CreateUser(user *User) (*User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/user", c.baseURL)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdUser User
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func TestAPIClient(t *testing.T) {
	// Create test server
	handler := NewUserHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create client pointing to test server
	client := NewAPIClient(server.URL)

	t.Run("get user", func(t *testing.T) {
		// First create a user
		user := &User{ID: 1, Name: "David", Email: "david@example.com"}
		createdUser, err := client.CreateUser(user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		// Now get the user
		fetchedUser, err := client.GetUser(createdUser.Email)
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		if fetchedUser.Name != "David" {
			t.Errorf("expected name David, got %s", fetchedUser.Name)
		}
	})

	t.Run("get non-existent user", func(t *testing.T) {
		_, err := client.GetUser("nonexistent@example.com")
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}

// ========================================
// EXAMPLE 3: DATABASE INTEGRATION TESTING
// ========================================

// UserRepository for database operations
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`
	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email FROM users WHERE email = $1`
	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// For this example, we'll use SQLite
	// In production, you might use PostgreSQL or another database
	// Note: You need to import _ "github.com/mattn/go-sqlite3" to use this

	// For now, we'll skip actual database setup if SQLite is not available
	// and just show the test structure

	t.Skip("database integration test - requires SQLite driver")

	/*
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open database: %v", err)
		}

		// Create schema
		schema := `
			CREATE TABLE users (
				id INTEGER PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE
			)
		`
		if _, err := db.Exec(schema); err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}

		t.Cleanup(func() {
			db.Close()
		})

		return db
	*/

	return nil
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Verify user was created
	found, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to find user: %v", err)
	}

	if found.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", found.Name)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test users
	users := []*User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	for _, user := range users {
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
	}

	t.Run("existing user", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, "alice@example.com")
		if err != nil {
			t.Fatalf("failed to find user: %v", err)
		}
		if found.Name != "Alice" {
			t.Errorf("expected name Alice, got %s", found.Name)
		}
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, "nonexistent@example.com")
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}

// ========================================
// EXAMPLE 4: TRANSACTION TESTING
// ========================================

func TestUserRepository_Transaction(t *testing.T) {
	t.Skip("transaction tests require actual database setup")
}

// ========================================
// EXAMPLE 5: TESTING WITH CONTEXT
// ========================================

func (r *UserRepository) FindWithTimeout(ctx context.Context, id int) (*User, error) {
	// Simulate slow query
	time.Sleep(10 * time.Millisecond)

	query := `SELECT id, name, email FROM users WHERE id = $1`
	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func TestUserRepository_ContextTimeout(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := repo.FindWithTimeout(ctx, 1)
	if err == nil {
		t.Error("expected timeout error")
	}
}

// ========================================
// EXAMPLE 6: TESTING MIDDLEWARE
// ========================================

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call next handler
		next.ServeHTTP(w, r)

		// Log request
		duration := time.Since(start)
		fmt.Printf("%s %s took %v\n", r.Method, r.URL.Path, duration)
	})
}

// AuthenticationMiddleware checks for API key
func AuthenticationMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key != apiKey {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrapped := LoggingMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAuthenticationMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	apiKey := "test-api-key-123"
	wrapped := AuthenticationMiddleware(apiKey)(handler)

	t.Run("valid API key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-API-Key", apiKey)
		w := httptest.NewRecorder()

		wrapped.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("invalid API key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-API-Key", "wrong-key")
		w := httptest.NewRecorder()

		wrapped.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("missing API key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		wrapped.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}

// ========================================
// EXAMPLE 7: TESTING JSON ENDPOINTS
// ========================================

type Request struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func ProcessRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	resp := Response{
		Status:  "success",
		Message: "processed: " + req.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func TestProcessRequestHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, Response)
	}{
		{
			name:           "valid request",
			requestBody:    `{"action": "process", "data": "test data"}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp Response) {
				t.Helper()
				if resp.Status != "success" {
					t.Errorf("expected status success, got %s", resp.Status)
				}
				if !contains(resp.Message, "test data") {
					t.Errorf("expected message to contain 'test data', got %s", resp.Message)
				}
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty request",
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/process", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			ProcessRequestHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK && tt.checkResponse != nil {
				var resp Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				tt.checkResponse(t, resp)
			}
		})
	}
}

// ========================================
// EXAMPLE 8: TESTING FILE UPLOADS
// ========================================

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"status":     "success",
		"filename":   header.Filename,
		"size":       len(content),
		"content_type": header.Header.Get("Content-Type"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func TestUploadHandler(t *testing.T) {
	t.Run("successful upload", func(t *testing.T) {
		// Create temporary file
		content := []byte("test file content")
		tmpFile, err := os.CreateTemp("", "test-upload-*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(content); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		// Create multipart form
		body := &bytes.Buffer{}
		// Note: In real code, use multipart writer here
		// For simplicity, we're skipping the multipart creation

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		w := httptest.NewRecorder()

		UploadHandler(w, req)

		t.Skip("full multipart test requires multipart writer setup")
	})
}

// ========================================
// EXAMPLE 9: TESTING WITH TEST SERVER
// ========================================

func TestHTTPServer(t *testing.T) {
	// Create a test server
	handler := NewUserHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make HTTP requests to the test server
	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("create and get user", func(t *testing.T) {
		// Create user
		user := User{ID: 1, Name: "Eve", Email: "eve@example.com"}
		body, _ := json.Marshal(user)

		resp, err := client.Post(server.URL+"/user", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		// Get user
		resp, err = client.Get(server.URL + "/user?email=eve@example.com")
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var fetchedUser User
		if err := json.NewDecoder(resp.Body).Decode(&fetchedUser); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if fetchedUser.Name != "Eve" {
			t.Errorf("expected name Eve, got %s", fetchedUser.Name)
		}
	})
}

// ========================================
// EXAMPLE 10: TESTING HTTP/2
// ========================================

func TestHTTP2Support(t *testing.T) {
	// Create test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Proto == "HTTP/2.0" {
			w.Write([]byte("Using HTTP/2"))
		} else {
			w.Write([]byte("Not using HTTP/2"))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	// Make request
	client := &http.Client{}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response: %s (Protocol: %s)", string(body), resp.Proto)
}

// Helper function for string matching
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
