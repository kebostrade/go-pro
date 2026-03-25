// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package security provides authentication, authorization, and security middleware.
package security

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CSRFConfig holds CSRF protection configuration.
type CSRFConfig struct {
	// Enabled determines if CSRF protection is active
	Enabled bool `json:"enabled" yaml:"enabled"`
	// Key is the secret key used to generate tokens
	Key string `json:"key" yaml:"key"`
	// CookieName is the name of the CSRF cookie
	CookieName string `json:"cookie_name" yaml:"cookie_name"`
	// HeaderName is the name of the header containing the CSRF token
	HeaderName string `json:"header_name" yaml:"header_name"`
	// FieldName is the name of the form field containing the CSRF token
	FieldName string `json:"field_name" yaml:"field_name"`
	// CookiePath is the path for the CSRF cookie
	CookiePath string `json:"cookie_path" yaml:"cookie_path"`
	// CookieDomain is the domain for the CSRF cookie
	CookieDomain string `json:"cookie_domain" yaml:"cookie_domain"`
	// CookieSecure determines if the cookie should be secure (HTTPS only)
	CookieSecure bool `json:"cookie_secure" yaml:"cookie_secure"`
	// CookieHTTPOnly determines if the cookie should be HTTP-only
	CookieHTTPOnly bool `json:"cookie_http_only" yaml:"cookie_http_only"`
	// CookieSameSite determines the SameSite attribute
	CookieSameSite http.SameSite `json:"cookie_same_site" yaml:"cookie_same_site"`
	// TokenLength is the length of the CSRF token in bytes
	TokenLength int `json:"token_length" yaml:"token_length"`
	// MaxAge is the maximum age of the CSRF token in seconds
	MaxAge int `json:"max_age" yaml:"max_age"`
	// TrustedOrigins is a list of trusted origins for CSRF validation
	TrustedOrigins []string `json:"trusted_origins" yaml:"trusted_origins"`
}

// DefaultCSRFConfig returns the default CSRF configuration.
func DefaultCSRFConfig() *CSRFConfig {
	return &CSRFConfig{
		Enabled:        true,
		Key:            generateRandomKey(32),
		CookieName:     "csrf_token",
		HeaderName:     "X-CSRF-Token",
		FieldName:      "csrf_token",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
		TokenLength:    32,
		MaxAge:         3600, // 1 hour
		TrustedOrigins: []string{},
	}
}

// CSRFToken represents a CSRF token with metadata.
type CSRFToken struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CSRFManager manages CSRF token generation and validation.
type CSRFManager struct {
	config *CSRFConfig
	tokens sync.Map // thread-safe token storage
}

// NewCSRFManager creates a new CSRF manager.
func NewCSRFManager(config *CSRFConfig) *CSRFManager {
	if config == nil {
		config = DefaultCSRFConfig()
	}
	return &CSRFManager{
		config: config,
	}
}

// GenerateToken generates a new CSRF token.
func (m *CSRFManager) GenerateToken() (string, error) {
	bytes := make([]byte, m.config.TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	csrfToken := &CSRFToken{
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(m.config.MaxAge) * time.Second),
	}

	m.tokens.Store(token, csrfToken)
	return token, nil
}

// ValidateToken validates a CSRF token.
func (m *CSRFManager) ValidateToken(token string) bool {
	if token == "" {
		return false
	}

	value, ok := m.tokens.Load(token)
	if !ok {
		return false
	}

	csrfToken, ok := value.(*CSRFToken)
	if !ok {
		return false
	}

	// Check if token has expired
	if time.Now().After(csrfToken.ExpiresAt) {
		m.tokens.Delete(token)
		return false
	}

	return true
}

// DeleteToken removes a CSRF token.
func (m *CSRFManager) DeleteToken(token string) {
	m.tokens.Delete(token)
}

// CleanupExpiredTokens removes all expired tokens.
func (m *CSRFManager) CleanupExpiredTokens() {
	now := time.Now()
	m.tokens.Range(func(key, value interface{}) bool {
		if csrfToken, ok := value.(*CSRFToken); ok {
			if now.After(csrfToken.ExpiresAt) {
				m.tokens.Delete(key)
			}
		}
		return true
	})
}

// CSRFMiddleware provides CSRF protection for HTTP handlers.
type CSRFMiddleware struct {
	manager *CSRFManager
	config  *CSRFConfig
}

// NewCSRFMiddleware creates a new CSRF middleware.
func NewCSRFMiddleware(config *CSRFConfig) *CSRFMiddleware {
	return &CSRFMiddleware{
		manager: NewCSRFManager(config),
		config:  config,
	}
}

// Handler returns the CSRF middleware handler.
func (m *CSRFMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF for safe methods
		if m.isSafeMethod(r.Method) {
			// Generate and set token for safe methods
			token, err := m.manager.GenerateToken()
			if err != nil {
				m.writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate CSRF token")
				return
			}
			m.setCSRFCookie(w, token)
			r.Header.Set(m.config.HeaderName, token)
			next.ServeHTTP(w, r)
			return
		}

		// For unsafe methods, validate the token
		if !m.config.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Get token from header or form
		token := m.getTokenFromRequest(r)
		if token == "" {
			m.writeErrorResponse(w, http.StatusForbidden, "CSRF token missing")
			return
		}

		// Validate token
		if !m.manager.ValidateToken(token) {
			m.writeErrorResponse(w, http.StatusForbidden, "Invalid CSRF token")
			return
		}

		// Validate origin if configured
		if !m.validateOrigin(r) {
			m.writeErrorResponse(w, http.StatusForbidden, "Invalid origin")
			return
		}

		// Generate new token for next request
		newToken, err := m.manager.GenerateToken()
		if err != nil {
			m.writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate CSRF token")
			return
		}
		m.setCSRFCookie(w, newToken)
		r.Header.Set(m.config.HeaderName, newToken)

		next.ServeHTTP(w, r)
	})
}

// isSafeMethod checks if the HTTP method is safe (doesn't modify state).
func (m *CSRFMiddleware) isSafeMethod(method string) bool {
	switch strings.ToUpper(method) {
	case "GET", "HEAD", "OPTIONS", "TRACE":
		return true
	default:
		return false
	}
}

// getTokenFromRequest extracts the CSRF token from the request.
func (m *CSRFMiddleware) getTokenFromRequest(r *http.Request) string {
	// Check header first
	token := r.Header.Get(m.config.HeaderName)
	if token != "" {
		return token
	}

	// Check form field
	token = r.FormValue(m.config.FieldName)
	if token != "" {
		return token
	}

	// Check cookie
	cookie, err := r.Cookie(m.config.CookieName)
	if err == nil {
		return cookie.Value
	}

	return ""
}

// setCSRFCookie sets the CSRF cookie in the response.
func (m *CSRFMiddleware) setCSRFCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.config.CookieName,
		Value:    token,
		Path:     m.config.CookiePath,
		Domain:   m.config.CookieDomain,
		MaxAge:   m.config.MaxAge,
		Secure:   m.config.CookieSecure,
		HttpOnly: m.config.CookieHTTPOnly,
		SameSite: m.config.CookieSameSite,
	})
}

// validateOrigin validates the request origin against trusted origins.
func (m *CSRFMiddleware) validateOrigin(r *http.Request) bool {
	if len(m.config.TrustedOrigins) == 0 {
		return true // No origin validation if not configured
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = r.Header.Get("Referer")
	}

	if origin == "" {
		return false // Require origin for unsafe methods
	}

	// Check if origin is in trusted list
	for _, trusted := range m.config.TrustedOrigins {
		if origin == trusted || strings.HasPrefix(origin, trusted) {
			return true
		}
	}

	return false
}

// writeErrorResponse writes a JSON error response.
func (m *CSRFMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   false,
		"error":     message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetToken returns the current CSRF token from the request.
func (m *CSRFMiddleware) GetToken(r *http.Request) string {
	return r.Header.Get(m.config.HeaderName)
}

// TokenHandler returns a handler that generates and returns a CSRF token.
func (m *CSRFMiddleware) TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := m.manager.GenerateToken()
		if err != nil {
			m.writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate CSRF token")
			return
		}

		m.setCSRFCookie(w, token)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"token":   token,
		})
	}
}

// Helper functions

// generateRandomKey generates a random key of specified length.
func generateRandomKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based key
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
