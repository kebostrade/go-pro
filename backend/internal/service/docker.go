// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides Docker environment management functionality.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DockerService handles Docker environment operations
type DockerService struct {
	basePath string // e.g., "basic/projects"
	timeout  time.Duration
}

// NewDockerService creates a new Docker service
func NewDockerService(basePath string) *DockerService {
	return &DockerService{
		basePath: basePath,
		timeout:  2 * time.Minute,
	}
}

// ServiceStatus represents a single container's status
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Health string `json:"health"`
}

// DockerStatus represents the status of a topic's Docker environment
type DockerStatus struct {
	TopicID    string            `json:"topic_id"`
	Status     string            `json:"status"` // "running", "stopped", "not_created", "error"
	Services   []ServiceStatus   `json:"services"`
	Ports      map[string]string `json:"ports"`
	Error      string            `json:"error,omitempty"`
	LastUpdate time.Time         `json:"last_update"`
}

// StartEnvironment starts Docker compose for a topic
func (s *DockerService) StartEnvironment(ctx context.Context, topicID string) (*DockerStatus, error) {
	composePath := filepath.Join(s.basePath, topicID, "docker-compose.yml")

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "up", "-d", "--remove-orphans")
	cmd.Dir = filepath.Dir(composePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &DockerStatus{
			TopicID:    topicID,
			Status:     "error",
			Error:      fmt.Sprintf("Failed to start: %s", string(output)),
			LastUpdate: time.Now(),
		}, err
	}

	return s.GetStatus(ctx, topicID)
}

// StopEnvironment stops Docker compose for a topic
func (s *DockerService) StopEnvironment(ctx context.Context, topicID string) (*DockerStatus, error) {
	composePath := filepath.Join(s.basePath, topicID, "docker-compose.yml")

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "down")
	cmd.Dir = filepath.Dir(composePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &DockerStatus{
			TopicID:    topicID,
			Status:     "error",
			Error:      fmt.Sprintf("Failed to stop: %s", string(output)),
			LastUpdate: time.Now(),
		}, err
	}

	return &DockerStatus{
		TopicID:    topicID,
		Status:     "stopped",
		Services:   []ServiceStatus{},
		Ports:      map[string]string{},
		LastUpdate: time.Now(),
	}, nil
}

// GetStatus returns current status of Docker environment
func (s *DockerService) GetStatus(ctx context.Context, topicID string) (*DockerStatus, error) {
	composePath := filepath.Join(s.basePath, topicID, "docker-compose.yml")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "ps", "--format", "json")
	cmd.Dir = filepath.Dir(composePath)

	output, err := cmd.Output()
	if err != nil {
		// docker compose ps returns exit 1 when no containers exist
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return &DockerStatus{
				TopicID:    topicID,
				Status:     "stopped",
				Services:   []ServiceStatus{},
				Ports:      map[string]string{},
				LastUpdate: time.Now(),
			}, nil
		}
		return nil, err
	}

	// Parse JSON output from docker compose ps
	var rawServices []map[string]interface{}
	if err := json.Unmarshal(output, &rawServices); err != nil {
		return nil, err
	}

	services := make([]ServiceStatus, 0, len(rawServices))
	ports := make(map[string]string)
	allRunning := len(rawServices) > 0

	for _, raw := range rawServices {
		name, _ := raw["Service"].(string)
		status, _ := raw["State"].(string)
		portStr, _ := raw["Ports"].(string)

		health := "starting"
		if strings.Contains(strings.ToLower(status), "running") {
			health = "healthy"
		} else if strings.Contains(strings.ToLower(status), "exited") {
			health = "stopped"
			allRunning = false
		}

		services = append(services, ServiceStatus{
			Name:   name,
			Status: status,
			Health: health,
		})

		if portStr != "" {
			ports[name] = portStr
		}
	}

	overallStatus := "running"
	if !allRunning {
		overallStatus = "stopped"
	}

	return &DockerStatus{
		TopicID:    topicID,
		Status:     overallStatus,
		Services:   services,
		Ports:      ports,
		LastUpdate: time.Now(),
	}, nil
}
