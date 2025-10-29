package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"go-pro-backend/internal/domain"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/pkg/logger"
)

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares to a handler
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// RequestID generates and adds a unique request ID to the context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to context
		ctx := logger.WithRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logging logs HTTP requests with structured logging
func Logging(log logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			logger.LogHTTPRequest(log, r.Context(),
				r.Method,
				r.URL.Path,
				r.UserAgent(),
				duration,
			)

			// Log additional details for errors or slow requests
			if wrapped.statusCode >= 400 || duration > 5*time.Second {
				log.Warn(r.Context(), "HTTP request attention required",
					"status_code", wrapped.statusCode,
					"duration_ms", duration.Milliseconds(),
					"remote_addr", r.RemoteAddr,
					"query_params", r.URL.RawQuery,
				)
			}
		})
	}
}

// CORS handles Cross-Origin Resource Sharing
func CORS(origins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Set default CORS headers
			if len(origins) == 0 || contains(origins, "*") {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if contains(origins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Request-ID, X-CSRF-Token")
			w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Total-Count, X-Page, X-Page-Size")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Recovery recovers from panics and returns a proper error response
func Recovery(log logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					logger.LogError(log, r.Context(),
						apierrors.NewInternalError("panic recovered", nil),
						"panic recovered",
						"panic_value", err,
						"stack_trace", string(debug.Stack()),
					)

					// Return error response
					WriteErrorResponse(w, r, apierrors.NewInternalError("internal server error", nil))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimit implements simple rate limiting (in production, use Redis or similar)
func RateLimit(requests int, window time.Duration) Middleware {
	type client struct {
		requests int
		window   time.Time
	}

	clients := make(map[string]*client)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			now := time.Now()

			// Clean up old entries (simplified cleanup)
			for k, v := range clients {
				if now.Sub(v.window) > window {
					delete(clients, k)
				}
			}

			// Check rate limit
			c, exists := clients[ip]
			if !exists {
				clients[ip] = &client{requests: 1, window: now}
				next.ServeHTTP(w, r)
				return
			}

			if now.Sub(c.window) > window {
				c.requests = 1
				c.window = now
				next.ServeHTTP(w, r)
				return
			}

			if c.requests >= requests {
				WriteErrorResponse(w, r, &apierrors.APIError{
					Type:       "RATE_LIMIT_EXCEEDED",
					Message:    "too many requests",
					StatusCode: http.StatusTooManyRequests,
				})
				return
			}

			c.requests++
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout adds a timeout to requests
func Timeout(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ContentType validates the Content-Type header for specific endpoints
func ContentType(contentType string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				ct := r.Header.Get("Content-Type")
				if ct != "" && !strings.HasPrefix(ct, contentType) {
					WriteErrorResponse(w, r, apierrors.NewBadRequestError("invalid content type"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Security adds security headers
func Security() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")

			next.ServeHTTP(w, r)
		})
	}
}

// Validation middleware for request validation
func Validation() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add validation context if needed
			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID generates a random request ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return strings.Split(xff, ",")[0]
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Use RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

// WriteErrorResponse writes a standardized error response
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apierrors.APIError
	var statusCode int

	if errors.As(err, &apiErr) {
		statusCode = apiErr.StatusCode
	} else {
		statusCode = http.StatusInternalServerError
		apiErr = apierrors.NewInternalError("internal server error", err)
	}

	response := &domain.APIResponse{
		Success: false,
		Error: &domain.APIError{
			Type:    apiErr.Type,
			Message: apiErr.Message,
		},
		RequestID: logger.GetRequestID(r.Context()),
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Use a simple JSON encoder to avoid import cycles
	jsonResponse := `{"success":false,"error":{"type":"` + apiErr.Type + `","message":"` + apiErr.Message + `"},"request_id":"` + response.RequestID + `","timestamp":"` + response.Timestamp.Format(time.RFC3339) + `"}`
	w.Write([]byte(jsonResponse))
}

// Pagination extracts pagination parameters from query string
func Pagination(defaultPageSize int, maxPageSize int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			page := 1
			pageSize := defaultPageSize

			if pageStr := r.URL.Query().Get("page"); pageStr != "" {
				if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
					page = p
				}
			}

			if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
				if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= maxPageSize {
					pageSize = ps
				}
			}

			pagination := &domain.PaginationRequest{
				Page:     page,
				PageSize: pageSize,
			}

			ctx := context.WithValue(r.Context(), "pagination", pagination)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
