// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package middleware provides distributed tracing middleware.
package middleware

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
	
	"go-pro-backend/internal/middleware/observability"
)

// TraceID represents a distributed trace ID.
type TraceID struct {
	Value    string
	Started time.Time
}

// TraceContext holds trace information.
type TraceContext struct {
	TraceID *TraceID
	SpanID  string
}

// NewTraceContext creates a new trace context.
func NewTraceContext() *TraceContext {
	return &TraceContext{
		TraceID: generateTraceID(),
		SpanID:  generateSpanID(),
	}
}

// generateTraceID creates a new trace ID.
func generateTraceID() *TraceID {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(rand.Int63()))
	return &TraceID{
		Value:    fmt.Sprintf("trace-%x", b),
		Started: time.Now(),
	}
}

// generateSpanID creates a new span ID.
func generateSpanID() string {
	rand.Seed(time.Now().UnixNano() + 1)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(rand.Int31()))
	return fmt.Sprintf("span-%x", b)
}

// WithTrace adds trace context to a request.
func WithTrace(ctx context.Context, tc *TraceContext) context.Context {
	return context.WithValue(ctx, observability.TraceContextKey, tc)
}

// GetTrace extracts trace context from a request.
func GetTrace(ctx context.Context) (*TraceContext, bool) {
	tc, ok := ctx.Value(observability.TraceContextKey).(*TraceContext)
	return tc, ok
}
