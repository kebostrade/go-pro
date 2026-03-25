// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package observability provides OpenTelemetry tracing and metrics.
// Note: This package is not yet fully implemented.
package observability

// Config holds OpenTelemetry configuration.
type Config struct {
	ServiceName    string
	Environment    string
	OTLPEndpoint   string
	StdoutExporter bool
	SamplingRate   float64
}

// DefaultConfig returns the default OpenTelemetry configuration.
func DefaultConfig() *Config {
	return &Config{
		ServiceName:    "go-pro-backend",
		Environment:    "development",
		OTLPEndpoint:   "http://localhost:4318",
		StdoutExporter: true,
		SamplingRate:   1.0,
	}
}
