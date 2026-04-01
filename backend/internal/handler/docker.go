// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package handler provides HTTP request handlers for the Docker API.
package handler

import (
	"encoding/json"
	"net/http"

	"go-pro-backend/internal/service"
)

// DockerHandler handles Docker environment HTTP requests
type DockerHandler struct {
	dockerService *service.DockerService
}

// NewDockerHandler creates a new Docker handler
func NewDockerHandler(dockerService *service.DockerService) *DockerHandler {
	return &DockerHandler{
		dockerService: dockerService,
	}
}

// DockerRequest represents incoming Docker API request
type DockerRequest struct {
	TopicID string `json:"topic_id"`
}

// DockerResponse wraps the DockerStatus in standard API response
type DockerResponse struct {
	Success bool                  `json:"success"`
	Data    *service.DockerStatus `json:"data,omitempty"`
	Error   string                `json:"error,omitempty"`
}

// handleDockerUp starts Docker environment for a topic
// POST /api/docker/up
func (h *DockerHandler) handleDockerUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DockerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.TopicID == "" {
		http.Error(w, "topic_id is required", http.StatusBadRequest)
		return
	}

	status, err := h.dockerService.StartEnvironment(r.Context(), req.TopicID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, &DockerResponse{
			Success: false,
			Data:    status,
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, &DockerResponse{
		Success: true,
		Data:    status,
	})
}

// handleDockerDown stops Docker environment for a topic
// POST /api/docker/down
func (h *DockerHandler) handleDockerDown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DockerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.TopicID == "" {
		http.Error(w, "topic_id is required", http.StatusBadRequest)
		return
	}

	status, err := h.dockerService.StopEnvironment(r.Context(), req.TopicID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, &DockerResponse{
			Success: false,
			Data:    status,
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, &DockerResponse{
		Success: true,
		Data:    status,
	})
}

// handleDockerStatus gets Docker environment status for a topic
// GET /api/docker/status?topic_id=xxx
func (h *DockerHandler) handleDockerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topicID := r.URL.Query().Get("topic_id")
	if topicID == "" {
		http.Error(w, "topic_id query parameter is required", http.StatusBadRequest)
		return
	}

	status, err := h.dockerService.GetStatus(r.Context(), topicID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, &DockerResponse{
			Success: false,
			Data:    status,
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, &DockerResponse{
		Success: true,
		Data:    status,
	})
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, resp *DockerResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
