package discovery

import (
	"fmt"
	"sync"
)

// ServiceRegistry manages service discovery
type ServiceRegistry struct {
	services map[string]string
	mu       sync.RWMutex
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]string),
	}
}

// Register registers a service with its address
func (r *ServiceRegistry) Register(name, address string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[name] = address
}

// Discover returns the address of a service
func (r *ServiceRegistry) Discover(name string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	address, ok := r.services[name]
	if !ok {
		return "", fmt.Errorf("service %s not found", name)
	}
	return address, nil
}

// Deregister removes a service from the registry
func (r *ServiceRegistry) Deregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services, name)
}

// List returns all registered services
func (r *ServiceRegistry) List() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make(map[string]string, len(r.services))
	for k, v := range r.services {
		services[k] = v
	}
	return services
}

// Global registry instance
var globalRegistry = NewServiceRegistry()

// Register registers a service globally
func Register(name, address string) {
	globalRegistry.Register(name, address)
}

// Discover discovers a service globally
func Discover(name string) (string, error) {
	return globalRegistry.Discover(name)
}

// Deregister deregisters a service globally
func Deregister(name string) {
	globalRegistry.Deregister(name)
}

// List lists all services globally
func List() map[string]string {
	return globalRegistry.List()
}

