package discovery

import (
	"fmt"
	"testing"
)

func TestRegisterAndDiscover(t *testing.T) {
	// Use unique service names to avoid conflicts
	serviceName := fmt.Sprintf("test-service-%d", 1)

	// Register a service
	Register(serviceName, "localhost:8080")
	defer Deregister(serviceName)

	// Discover the service
	addr, err := Discover(serviceName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if addr != "localhost:8080" {
		t.Errorf("Expected localhost:8080, got %s", addr)
	}
}

func TestDiscoverNonExistent(t *testing.T) {
	// Try to discover non-existent service
	_, err := Discover("non-existent-service-xyz")
	if err == nil {
		t.Fatal("Expected error for non-existent service")
	}
}

func TestDeregister(t *testing.T) {
	serviceName := fmt.Sprintf("test-service-%d", 2)

	// Register and then deregister
	Register(serviceName, "localhost:8080")
	Deregister(serviceName)

	// Should not be discoverable
	_, err := Discover(serviceName)
	if err == nil {
		t.Fatal("Expected error after deregistration")
	}
}

func TestList(t *testing.T) {
	// Register multiple services with unique names
	service1 := fmt.Sprintf("service1-%d", 3)
	service2 := fmt.Sprintf("service2-%d", 3)
	service3 := fmt.Sprintf("service3-%d", 3)

	Register(service1, "localhost:8081")
	Register(service2, "localhost:8082")
	Register(service3, "localhost:8083")

	defer func() {
		Deregister(service1)
		Deregister(service2)
		Deregister(service3)
	}()

	// List all services
	services := List()

	// Should have at least our 3 services
	if len(services) < 3 {
		t.Errorf("Expected at least 3 services, got %d", len(services))
	}

	if services[service1] != "localhost:8081" {
		t.Error("service1 not found or incorrect address")
	}
	if services[service2] != "localhost:8082" {
		t.Error("service2 not found or incorrect address")
	}
	if services[service3] != "localhost:8083" {
		t.Error("service3 not found or incorrect address")
	}
}

func TestConcurrentRegistration(t *testing.T) {
	// Register services concurrently
	done := make(chan bool)
	serviceNames := make([]string, 10)

	for i := 0; i < 10; i++ {
		serviceNames[i] = fmt.Sprintf("concurrent-service-%d", i)
		go func(id int) {
			Register(serviceNames[id], fmt.Sprintf("localhost:808%d", id))
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Cleanup
	defer func() {
		for _, name := range serviceNames {
			Deregister(name)
		}
	}()

	// Check all services registered
	services := List()
	registeredCount := 0
	for _, name := range serviceNames {
		if _, exists := services[name]; exists {
			registeredCount++
		}
	}

	if registeredCount != 10 {
		t.Errorf("Expected 10 services registered, got %d", registeredCount)
	}
}
