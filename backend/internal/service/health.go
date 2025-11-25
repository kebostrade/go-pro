// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"fmt"
	"time"

	"go-pro-backend/internal/domain"
)

// healthService implements the HealthService interface.
type healthService struct {
	startTime time.Time
	version   string
}

// NewHealthService creates a new health service.
func NewHealthService(version string) HealthService {
	return &healthService{
		startTime: time.Now(),
		version:   version,
	}
}

// GetHealthStatus returns the health status of the application.
func (s *healthService) GetHealthStatus(ctx context.Context) (*domain.HealthResponse, error) {
	uptime := time.Since(s.startTime)

	response := &domain.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   s.version,
		Uptime:    formatUptime(uptime),
	}

	return response, nil
}

// formatUptime formats duration into a human-readable string.
func formatUptime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}

	return fmt.Sprintf("%ds", seconds)
}
