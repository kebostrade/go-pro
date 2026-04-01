package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry("http://localhost:8001", "http://localhost:8002")

	if registry.GetUserServiceURL() != "http://localhost:8001" {
		t.Errorf("expected user service URL 'http://localhost:8001', got %s", registry.GetUserServiceURL())
	}
	if registry.GetOrderServiceURL() != "http://localhost:8002" {
		t.Errorf("expected order service URL 'http://localhost:8002', got %s", registry.GetOrderServiceURL())
	}
}

func TestRegistryGetUserServiceProxy(t *testing.T) {
	registry := NewRegistry("http://localhost:8001", "http://localhost:8002")

	proxy := registry.GetUserServiceProxy()
	if proxy == nil {
		t.Fatal("expected non-nil proxy for user service")
	}
}

func TestRegistryGetOrderServiceProxy(t *testing.T) {
	registry := NewRegistry("http://localhost:8001", "http://localhost:8002")

	proxy := registry.GetOrderServiceProxy()
	if proxy == nil {
		t.Fatal("expected non-nil proxy for order service")
	}
}

func TestProxyServeUserService(t *testing.T) {
	// Create a mock upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer upstream.Close()

	registry := NewRegistry(upstream.URL, "http://localhost:8002")
	proxy := NewProxy(registry)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	proxy.ServeUserService(w, req)

	// Should not panic and should complete
	if w.Code == 0 {
		t.Error("response code not set")
	}
}

func TestProxyServeOrderService(t *testing.T) {
	// Create a mock upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer upstream.Close()

	registry := NewRegistry("http://localhost:8001", upstream.URL)
	proxy := NewProxy(registry)

	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	w := httptest.NewRecorder()

	proxy.ServeOrderService(w, req)

	// Should not panic and should complete
	if w.Code == 0 {
		t.Error("response code not set")
	}
}
