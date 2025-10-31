// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package security provides authentication, authorization, and security middleware.
package security

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Middleware key types for context.
type contextKey string

const (
	UserContextKey   contextKey = "user"
	ClaimsContextKey contextKey = "claims"
)

// SecurityMiddleware holds all security middleware functions.
type SecurityMiddleware struct {
	config         *SecurityConfig
	jwtManager     *JWTManager
	rateLimitStore *RateLimitStore
}

// NewSecurityMiddleware creates a new security middleware instance.
func NewSecurityMiddleware(config *SecurityConfig) *SecurityMiddleware {
	return &SecurityMiddleware{
		config:         config,
		jwtManager:     NewJWTManager(config.JWT),
		rateLimitStore: NewRateLimitStore(config.RateLimit),
	}
}

// SecurityHeaders middleware adds security headers to all responses.
func (sm *SecurityMiddleware) SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := sm.config.Headers

		// Content Security Policy.
		if headers.EnableCSP {
			w.Header().Set("Content-Security-Policy", headers.CSPPolicy)
		}

		// HTTP Strict Transport Security.
		if headers.EnableHSTS && sm.config.HTTPS.Enabled {
			w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d; includeSubDomains; preload", sm.config.HTTPS.HSTSMaxAge))
		}

		// X-Frame-Options.
		if headers.EnableFrameOptions {
			w.Header().Set("X-Frame-Options", "DENY")
		}

		// X-Content-Type-Options.
		if headers.EnableContentType {
			w.Header().Set("X-Content-Type-Options", "nosniff")
		}

		// Referrer Policy.
		if headers.EnableReferrerPolicy {
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		}

		// Additional security headers.
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Remove server info.
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")

		next.ServeHTTP(w, r)
	})
}

// CORS middleware with secure configuration.
func (sm *SecurityMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed.
		if sm.isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Set CORS headers.
		if sm.config.CORS.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(sm.config.CORS.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(sm.config.CORS.AllowedHeaders, ", "))
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(sm.config.CORS.MaxAge))

		// Handle preflight requests.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isOriginAllowed checks if the origin is in the allowed list.
func (sm *SecurityMiddleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range sm.config.CORS.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}

	return false
}

// RateLimit middleware implements token bucket rate limiting.
func (sm *SecurityMiddleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		if !sm.rateLimitStore.Allow(clientIP) {
			sm.writeErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JWTAuth middleware validates JWT tokens.
func (sm *SecurityMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ExtractTokenFromHeader(r.Header.Get("Authorization"))
		if token == "" {
			sm.writeErrorResponse(w, http.StatusUnauthorized, "Missing or invalid token", nil)
			return
		}

		claims, err := sm.jwtManager.ValidateToken(token)
		if err != nil {
			sm.writeErrorResponse(w, http.StatusUnauthorized, "Invalid token", map[string]string{"error": err.Error()})
			return
		}

		// Only allow access tokens for authentication.
		if !claims.IsAccessToken() {
			sm.writeErrorResponse(w, http.StatusUnauthorized, "Invalid token type", nil)
			return
		}

		// Add claims to request context.
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		ctx = context.WithValue(ctx, UserContextKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// APIKeyAuth middleware validates API keys for admin endpoints.
func (sm *SecurityMiddleware) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(sm.config.APIKey.KeyHeader)
		if apiKey == "" {
			sm.writeErrorResponse(w, http.StatusUnauthorized, "Missing API key", nil)
			return
		}

		// Use constant time comparison to prevent timing attacks.
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(sm.config.APIKey.AdminKey)) != 1 {
			sm.writeErrorResponse(w, http.StatusUnauthorized, "Invalid API key", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRoles middleware checks if user has required roles.
func (sm *SecurityMiddleware) RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
			if !ok {
				sm.writeErrorResponse(w, http.StatusUnauthorized, "No authentication context", nil)
				return
			}

			if !claims.HasAnyRole(roles) {
				sm.writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// InputValidation middleware validates and sanitizes input.
func (sm *SecurityMiddleware) InputValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limit request body size.
		if r.ContentLength > sm.config.Validation.MaxJSONSize {
			sm.writeErrorResponse(w, http.StatusRequestEntityTooLarge, "Request body too large", nil)
			return
		}

		// Validate content type for POST/PUT requests.
		if r.Method == "POST" || r.Method == "PUT" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
				sm.writeErrorResponse(w, http.StatusUnsupportedMediaType, "Unsupported content type", nil)
				return
			}
		}

		// Sanitize query parameters.
		if sm.config.Validation.SanitizeInput {
			sm.sanitizeQueryParams(r)
		}

		next.ServeHTTP(w, r)
	})
}

// sanitizeQueryParams sanitizes query parameters to prevent XSS.
func (sm *SecurityMiddleware) sanitizeQueryParams(r *http.Request) {
	query := r.URL.Query()
	for key, values := range query {
		for i, value := range values {
			query[key][i] = html.EscapeString(value)
		}
	}
	r.URL.RawQuery = query.Encode()
}

// SecureLogging middleware logs requests securely.
func (sm *SecurityMiddleware) SecureLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code.
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Log request details (excluding sensitive information)
		logData := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": wrapped.statusCode,
			"duration_ms": duration.Milliseconds(),
			"client_ip":   getClientIP(r),
			"user_agent":  r.Header.Get("User-Agent"),
			"timestamp":   start.UTC().Format(time.RFC3339),
		}

		// Add user context if available (don't log sensitive data)
		if userID := r.Context().Value(UserContextKey); userID != nil {
			logData["user_id"] = userID
		}

		// Don't log sensitive headers or query parameters.
		if !sm.config.Logging.LogSensitiveData {
			// Remove sensitive query parameters.
			safeQuery := r.URL.Query()
			for key := range safeQuery {
				if sm.isSensitiveParam(key) {
					safeQuery.Set(key, "[REDACTED]")
				}
			}
			logData["query"] = safeQuery.Encode()
		}

		logJSON, _ := json.Marshal(logData)
		log.Printf("REQUEST: %s", string(logJSON))
	})
}

// HTTPSRedirect middleware redirects HTTP to HTTPS.
func (sm *SecurityMiddleware) HTTPSRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sm.config.HTTPS.RedirectHTTP && r.Header.Get("X-Forwarded-Proto") != "https" {
			httpsURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
			http.Redirect(w, r, httpsURL, http.StatusPermanentRedirect)

			return
		}
		next.ServeHTTP(w, r)
	})
}

// Helper functions.

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// getClientIP extracts the real client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (most common)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header.
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to remote address.
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return ip
}

// isSensitiveParam checks if a parameter name is sensitive.
func (sm *SecurityMiddleware) isSensitiveParam(param string) bool {
	sensitive := []string{
		"password", "pwd", "secret", "token", "key", "auth",
		"api_key", "apikey", "session", "csrf", "signature",
	}

	param = strings.ToLower(param)
	for _, s := range sensitive {
		if strings.Contains(param, s) {
			return true
		}
	}

	return false
}

// writeErrorResponse writes a standardized error response.
func (sm *SecurityMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, message string, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success":   false,
		"error":     message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if details != nil && len(details) > 0 {
		response["details"] = details
	}

	json.NewEncoder(w).Encode(response)
}

// RateLimitStore implements a token bucket rate limiter.
type RateLimitStore struct {
	config  RateLimitConfig
	clients map[string]*tokenBucket
	mutex   sync.RWMutex
}

type tokenBucket struct {
	tokens     int
	lastRefill time.Time
}

// NewRateLimitStore creates a new rate limit store.
func NewRateLimitStore(config RateLimitConfig) *RateLimitStore {
	store := &RateLimitStore{
		config:  config,
		clients: make(map[string]*tokenBucket),
	}

	// Start cleanup routine.
	go store.cleanup()

	return store
}

// Allow checks if the client is allowed to make a request.
func (rls *RateLimitStore) Allow(clientIP string) bool {
	rls.mutex.Lock()
	defer rls.mutex.Unlock()

	bucket, exists := rls.clients[clientIP]
	if !exists {
		bucket = &tokenBucket{
			tokens:     rls.config.RequestsPerMinute,
			lastRefill: time.Now(),
		}
		rls.clients[clientIP] = bucket
	}

	// Refill tokens.
	rls.refillTokens(bucket)

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// refillTokens refills the token bucket based on elapsed time.
func (rls *RateLimitStore) refillTokens(bucket *tokenBucket) {
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)

	if elapsed >= rls.config.WindowSize {
		bucket.tokens = rls.config.RequestsPerMinute
		bucket.lastRefill = now
	}
}

// cleanup removes old client entries.
func (rls *RateLimitStore) cleanup() {
	ticker := time.NewTicker(rls.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rls.mutex.Lock()
		now := time.Now()
		for clientIP, bucket := range rls.clients {
			if now.Sub(bucket.lastRefill) > rls.config.CleanupInterval*2 {
				delete(rls.clients, clientIP)
			}
		}
		rls.mutex.Unlock()
	}
}

// Input validation helpers.

var (
	// Regex patterns for validation.
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	uuidRegex     = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	alphanumRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ValidateEmail validates email format.
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidateUUID validates UUID format.
func ValidateUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

// ValidateAlphanumeric validates alphanumeric strings.
func ValidateAlphanumeric(str string) bool {
	return alphanumRegex.MatchString(str)
}

// SanitizeString removes potentially dangerous characters.
func SanitizeString(input string, maxLength int) string {
	// Remove HTML/JS dangerous characters.
	input = html.EscapeString(input)

	// Truncate if too long.
	if len(input) > maxLength {
		input = input[:maxLength]
	}

	return strings.TrimSpace(input)
}
