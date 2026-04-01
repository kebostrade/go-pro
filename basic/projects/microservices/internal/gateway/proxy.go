// Package gateway provides HTTP proxy functionality.
package gateway

import (
	"net/http"
)

// Proxy handles HTTP reverse proxying to backend services.
type Proxy struct {
	registry *Registry
}

// NewProxy creates a new proxy instance.
func NewProxy(registry *Registry) *Proxy {
	return &Proxy{
		registry: registry,
	}
}

// ServeUserService proxies requests to the user service.
func (p *Proxy) ServeUserService(w http.ResponseWriter, r *http.Request) {
	proxy := p.registry.GetUserServiceProxy()
	proxy.ServeHTTP(w, r)
}

// ServeOrderService proxies requests to the order service.
func (p *Proxy) ServeOrderService(w http.ResponseWriter, r *http.Request) {
	proxy := p.registry.GetOrderServiceProxy()
	proxy.ServeHTTP(w, r)
}
