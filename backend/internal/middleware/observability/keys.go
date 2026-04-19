// GO-PRO Learning Platform Backend
// Package observability provides observability context keys.
package observability

import "context"

type contextKey string

var (
	TraceContextKey contextKey = "trace_context"
	SpanKey        contextKey = "span"
	LoggerKey      contextKey = "logger"
	MetricsKey    contextKey = "metrics"
)

func (c contextKey) String() string { return string(c) }
