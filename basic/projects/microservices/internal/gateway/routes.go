// Package gateway provides route configuration.
package gateway

import (
	"github.com/go-chi/chi/v5"
)

// MountRoutes registers all gateway routes on the chi router.
func MountRoutes(r *chi.Mux, proxy *Proxy) {
	// User service routes - proxy to service-a.
	r.Route("/api/users", func(r chi.Router) {
		r.Get("/*", proxy.ServeUserService)
		r.Get("", proxy.ServeUserService)
	})

	// Order service routes - proxy to service-b.
	r.Route("/api/orders", func(r chi.Router) {
		r.Get("/*", proxy.ServeOrderService)
		r.Get("", proxy.ServeOrderService)
	})
}
