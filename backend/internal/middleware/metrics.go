// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package middleware provides performance monitoring functionality.
package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Metrics holds Prometheus-compatible metrics for monitoring.
type Metrics struct {
	// Request duration histogram buckets (in seconds)
	requestDurations map[string][]float64
	// Response size histogram (in bytes)
	responseSizes map[string][]int
	// Error rate counter (status codes >= 400)
	errorCounts map[int]int
	// Total request counter per endpoint
	requestCounts map[string]int

	mu sync.RWMutex
}

// NewMetrics creates a new metrics collector.
func NewMetrics() *Metrics {
	return &Metrics{
		requestDurations: make(map[string][]float64),
		responseSizes:    make(map[string][]int),
		errorCounts:      make(map[int]int),
		requestCounts:    make(map[string]int),
	}
}

// MetricsMiddleware creates middleware that collects performance metrics.
func MetricsMiddleware(metrics *Metrics) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code and size
			mw := &metricsWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				bytesWritten:   0,
			}

			// Execute request
			next.ServeHTTP(mw, r)

			// Calculate metrics
			duration := time.Since(start).Seconds()
			path := r.URL.Path
			statusCode := mw.statusCode
			responseSize := mw.bytesWritten

			// Record metrics
			metrics.RecordRequest(path, duration, statusCode, responseSize)
		})
	}
}

// metricsWriter wraps http.ResponseWriter to capture metrics.
type metricsWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (mw *metricsWriter) WriteHeader(code int) {
	mw.statusCode = code
	mw.ResponseWriter.WriteHeader(code)
}

func (mw *metricsWriter) Write(b []byte) (int, error) {
	n, err := mw.ResponseWriter.Write(b)
	mw.bytesWritten += n
	return n, err
}

// RecordRequest records request metrics.
func (m *Metrics) RecordRequest(path string, duration float64, statusCode, responseSize int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record request duration
	m.requestDurations[path] = append(m.requestDurations[path], duration)

	// Record response size
	m.responseSizes[path] = append(m.responseSizes[path], responseSize)

	// Record error count
	if statusCode >= 400 {
		m.errorCounts[statusCode]++
	}

	// Record total requests
	m.requestCounts[path]++

	// Cleanup: Keep only last 1000 samples per endpoint to prevent unbounded growth
	if len(m.requestDurations[path]) > 1000 {
		m.requestDurations[path] = m.requestDurations[path][100:]
	}
	if len(m.responseSizes[path]) > 1000 {
		m.responseSizes[path] = m.responseSizes[path][100:]
	}
}

// GetMetrics returns current metrics snapshot.
func (m *Metrics) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]interface{})

	// Request duration statistics per endpoint
	durationStats := make(map[string]map[string]interface{})
	for path, durations := range m.requestDurations {
		if len(durations) == 0 {
			continue
		}
		durationStats[path] = map[string]interface{}{
			"count": len(durations),
			"avg":   avg(durations),
			"p50":   percentile(durations, 50),
			"p95":   percentile(durations, 95),
			"p99":   percentile(durations, 99),
			"min":   min(durations),
			"max":   max(durations),
		}
	}
	metrics["request_duration_seconds"] = durationStats

	// Response size statistics per endpoint
	sizeStats := make(map[string]map[string]interface{})
	for path, sizes := range m.responseSizes {
		if len(sizes) == 0 {
			continue
		}
		sizesFloat := intsToFloats(sizes)
		sizeStats[path] = map[string]interface{}{
			"count": len(sizes),
			"avg":   avg(sizesFloat),
			"p50":   percentile(sizesFloat, 50),
			"p95":   percentile(sizesFloat, 95),
			"p99":   percentile(sizesFloat, 99),
			"min":   min(sizesFloat),
			"max":   max(sizesFloat),
		}
	}
	metrics["response_size_bytes"] = sizeStats

	// Error counts by status code
	metrics["error_count_by_status"] = m.errorCounts

	// Total request counts
	metrics["request_count_by_endpoint"] = m.requestCounts

	// Calculate overall error rate
	totalRequests := 0
	totalErrors := 0
	for _, count := range m.requestCounts {
		totalRequests += count
	}
	for _, count := range m.errorCounts {
		totalErrors += count
	}
	errorRate := 0.0
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100
	}
	metrics["error_rate_percent"] = errorRate

	return metrics
}

// PrometheusMetrics returns metrics in Prometheus format.
// This is a simplified Prometheus exporter that returns metrics as plain text.
func (m *Metrics) PrometheusMetrics() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var output string

	// Request duration histogram
	output += "# HELP http_request_duration_seconds HTTP request duration in seconds\n"
	output += "# TYPE http_request_duration_seconds histogram\n"
	for path, durations := range m.requestDurations {
		if len(durations) == 0 {
			continue
		}
		output += "http_request_duration_seconds{path=\"" + path + "\",quantile=\"0.5\"} " +
			strconv.FormatFloat(percentile(durations, 50), 'f', 6, 64) + "\n"
		output += "http_request_duration_seconds{path=\"" + path + "\",quantile=\"0.95\"} " +
			strconv.FormatFloat(percentile(durations, 95), 'f', 6, 64) + "\n"
		output += "http_request_duration_seconds{path=\"" + path + "\",quantile=\"0.99\"} " +
			strconv.FormatFloat(percentile(durations, 99), 'f', 6, 64) + "\n"
		output += "http_request_duration_seconds_count{path=\"" + path + "\"} " +
			strconv.Itoa(len(durations)) + "\n"
	}

	// Response size histogram
	output += "# HELP http_response_size_bytes HTTP response size in bytes\n"
	output += "# TYPE http_response_size_bytes histogram\n"
	for path, sizes := range m.responseSizes {
		if len(sizes) == 0 {
			continue
		}
		sizesFloat := intsToFloats(sizes)
		output += "http_response_size_bytes{path=\"" + path + "\",quantile=\"0.5\"} " +
			strconv.FormatFloat(percentile(sizesFloat, 50), 'f', 0, 64) + "\n"
		output += "http_response_size_bytes{path=\"" + path + "\",quantile=\"0.95\"} " +
			strconv.FormatFloat(percentile(sizesFloat, 95), 'f', 0, 64) + "\n"
		output += "http_response_size_bytes{path=\"" + path + "\",quantile=\"0.99\"} " +
			strconv.FormatFloat(percentile(sizesFloat, 99), 'f', 0, 64) + "\n"
		output += "http_response_size_bytes_count{path=\"" + path + "\"} " +
			strconv.Itoa(len(sizes)) + "\n"
	}

	// Error count
	output += "# HELP http_request_errors_total Total number of HTTP errors by status code\n"
	output += "# TYPE http_request_errors_total counter\n"
	for statusCode, count := range m.errorCounts {
		output += "http_request_errors_total{status_code=\"" + strconv.Itoa(statusCode) + "\"} " +
			strconv.Itoa(count) + "\n"
	}

	// Request count
	output += "# HELP http_requests_total Total number of HTTP requests by endpoint\n"
	output += "# TYPE http_requests_total counter\n"
	for path, count := range m.requestCounts {
		output += "http_requests_total{path=\"" + path + "\"} " + strconv.Itoa(count) + "\n"
	}

	return output
}

// Helper functions for statistics

func avg(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	// Simple percentile calculation (sorted data not required for rough estimate)
	sorted := make([]float64, len(values))
	copy(sorted, values)
	// Bubble sort for simplicity (production code should use sort.Float64s)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	index := int(float64(len(sorted)-1) * p / 100.0)
	return sorted[index]
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func intsToFloats(values []int) []float64 {
	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = float64(v)
	}
	return result
}
