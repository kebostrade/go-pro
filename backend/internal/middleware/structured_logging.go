// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package middleware provides HTTP middleware for the GO-PRO API.
package middleware

import (
	"encoding/json"
	"time"
)

// StructuredLogEntry represents a structured log entry.
type StructuredLogEntry struct {
	Timestamp string            `json:"timestamp"`
	Level     string            `json:"level"`
	Message  string            `json:"message"`
	Request  string            `json:"request_id,omitempty"`
	Duration time.Duration     `json:"duration,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
}

// ToJSON converts the log entry to JSON.
func (s *StructuredLogEntry) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// NewStructuredLogEntry creates a new structured log entry.
func NewStructuredLogEntry(level, message string) *StructuredLogEntry {
	return &StructuredLogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:    level,
		Message:  message,
		Fields:   make(map[string]string),
	}
}
