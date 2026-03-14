package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// ServiceRegistry holds information about available services
type ServiceRegistry struct {
	userServiceURL  string
	orderServiceURL string
	log             *logrus.Logger
	client          *http.Client
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(userURL, orderURL string, log *logrus.Logger) *ServiceRegistry {
	return &ServiceRegistry{
		userServiceURL:  userURL,
		orderServiceURL: orderURL,
		log:             log,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyRequest forwards a request to the appropriate service
func (sr *ServiceRegistry) ProxyRequest(service, path string, method string, body []byte, headers http.Header) (*http.Response, error) {
	var url string
	switch service {
	case "users":
		url = sr.userServiceURL + path
	case "orders":
		url = sr.orderServiceURL + path
	default:
		return nil, fmt.Errorf("unknown service: %s", service)
	}

	sr.log.WithFields(logrus.Fields{
		"service": service,
		"method":  method,
		"path":    path,
		"url":     url,
	}).Debug("Proxying request")

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy relevant headers
	for key, values := range headers {
		// Skip headers that shouldn't be forwarded
		if key == "Content-Length" || key == "Host" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return sr.client.Do(req)
}

// CheckServiceHealth checks the health of a specific service
func (sr *ServiceRegistry) CheckServiceHealth(service string) error {
	var url string
	switch service {
	case "service-a":
		url = sr.userServiceURL + "/health"
	case "service-b":
		url = sr.orderServiceURL + "/health"
	default:
		return fmt.Errorf("unknown service: %s", service)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := sr.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service unhealthy: status %d", resp.StatusCode)
	}

	return nil
}

// ReadResponseBody reads and returns the response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}
