package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	t.Run("successful fetch", func(t *testing.T) {
		// Create test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/users/123", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			user := User{ID: "123", Name: "TestUser", Email: "test@example.com"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		}))
		defer ts.Close()

		// Create client and call
		client := NewAPIClient(ts.URL)
		user, err := client.GetUser(context.Background(), "123")

		assert.NoError(t, err)
		assert.Equal(t, "123", user.ID)
		assert.Equal(t, "TestUser", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("user not found", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		client := NewAPIClient(ts.URL)
		user, err := client.GetUser(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("server error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		client := NewAPIClient(ts.URL)
		user, err := client.GetUser(context.Background(), "123")

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/users", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var reqBody map[string]string
			json.NewDecoder(r.Body).Decode(&reqBody)
			assert.Equal(t, "Alice", reqBody["name"])
			assert.Equal(t, "alice@example.com", reqBody["email"])

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			user := User{ID: "new-id", Name: "Alice", Email: "alice@example.com"}
			json.NewEncoder(w).Encode(user)
		}))
		defer ts.Close()

		client := NewAPIClient(ts.URL)
		user, err := client.CreateUser(context.Background(), "Alice", "alice@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "new-id", user.ID)
		assert.Equal(t, "Alice", user.Name)
	})

	t.Run("creation failed", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer ts.Close()

		client := NewAPIClient(ts.URL)
		user, err := client.CreateUser(context.Background(), "Bob", "bob@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
