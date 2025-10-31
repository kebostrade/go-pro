// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package logger provides structured logging functionality.
package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// Logger interface defines logging methods.
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

// slogLogger implements Logger using structured logging.
type slogLogger struct {
	logger *slog.Logger
}

// New creates a new structured logger.
func New(level, format string) Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	return &slogLogger{logger: logger}
}

// Debug logs a debug message with context.
func (l *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// Info logs an info message with context.
func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// Warn logs a warning message with context.
func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// Error logs an error message with context.
func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

// With returns a new logger with additional context.
func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{logger: l.logger.With(args...)}
}

// RequestIDKey is the context key for request IDs.
type RequestIDKey struct{}

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey{}, requestID)
}

// GetRequestID retrieves the request ID from context.
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return requestID
	}

	return ""
}

// LogHTTPRequest logs HTTP request details.
func LogHTTPRequest(logger Logger, ctx context.Context, method, path, userAgent string, duration time.Duration) {
	logger.Info(ctx, "HTTP request completed",
		"method", method,
		"path", path,
		"user_agent", userAgent,
		"duration_ms", duration.Milliseconds(),
		"request_id", GetRequestID(ctx),
	)
}

// LogError logs an error with additional context.
func LogError(logger Logger, ctx context.Context, err error, message string, args ...any) {
	allArgs := append([]any{"error", err.Error()}, args...)
	logger.Error(ctx, message, allArgs...)
}
