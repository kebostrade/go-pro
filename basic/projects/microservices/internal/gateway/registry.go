// Package gateway provides API Gateway functionality.
package gateway

import (
	"net/http/httputil"
	"net/url"
)

// Registry holds service endpoint URLs.
type Registry struct {
	userServiceURL  string
	orderServiceURL string
}

// NewRegistry creates a new service registry.
func NewRegistry(userServiceURL, orderServiceURL string) *Registry {
	return &Registry{
		userServiceURL:  userServiceURL,
		orderServiceURL: orderServiceURL,
	}
}

// GetUserServiceURL returns the user service URL.
func (r *Registry) GetUserServiceURL() string {
	return r.userServiceURL
}

// GetOrderServiceURL returns the order service URL.
func (r *Registry) GetOrderServiceURL() string {
	return r.orderServiceURL
}

// GetUserServiceProxy returns a reverse proxy for the user service.
func (r *Registry) GetUserServiceProxy() *httputil.ReverseProxy {
	userURL, _ := url.Parse(r.userServiceURL)
	return httputil.NewSingleHostReverseProxy(userURL)
}

// GetOrderServiceProxy returns a reverse proxy for the order service.
func (r *Registry) GetOrderServiceProxy() *httputil.ReverseProxy {
	orderURL, _ := url.Parse(r.orderServiceURL)
	return httputil.NewSingleHostReverseProxy(orderURL)
}
