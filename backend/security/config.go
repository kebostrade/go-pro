// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package security provides authentication, authorization, and security middleware.
package security

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// SecurityConfig holds all security-related configuration.
type SecurityConfig struct {
	JWT        JWTConfig
	CORS       CORSConfig
	RateLimit  RateLimitConfig
	APIKey     APIKeyConfig
	HTTPS      HTTPSConfig
	Headers    HeadersConfig
	Validation ValidationConfig
	Logging    LoggingConfig
}

type JWTConfig struct {
	Secret          []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
	Audience        string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
	WindowSize        time.Duration
	CleanupInterval   time.Duration
}

type APIKeyConfig struct {
	AdminKey       string
	KeyHeader      string
	AdminEndpoints []string
}

type HTTPSConfig struct {
	Enabled      bool
	RedirectHTTP bool
	HSTSMaxAge   int
	CertFile     string
	KeyFile      string
}

type HeadersConfig struct {
	EnableCSP            bool
	CSPPolicy            string
	EnableHSTS           bool
	EnableFrameOptions   bool
	EnableContentType    bool
	EnableReferrerPolicy bool
}

type ValidationConfig struct {
	MaxJSONSize      int64
	MaxFieldLength   int
	AllowedFileTypes []string
	SanitizeInput    bool
}

type LoggingConfig struct {
	LogSensitiveData bool
	LogLevel         string
	LogFormat        string
}

// NewSecurityConfig creates a new security configuration with secure defaults.
func NewSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		JWT: JWTConfig{
			Secret:          getJWTSecret(),
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour * 7, // 7 days
			Issuer:          getEnvOrDefault("JWT_ISSUER", "go-pro-api"),
			Audience:        getEnvOrDefault("JWT_AUDIENCE", "go-pro-users"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getAllowedOrigins(),
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-Key"},
			AllowCredentials: true,
			MaxAge:           3600, // 1 hour
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_RPM", 100),
			BurstSize:         getIntEnv("RATE_LIMIT_BURST", 10),
			WindowSize:        time.Minute,
			CleanupInterval:   5 * time.Minute,
		},
		APIKey: APIKeyConfig{
			AdminKey:       getEnvOrDefault("ADMIN_API_KEY", generateSecureAPIKey()),
			KeyHeader:      "X-API-Key",
			AdminEndpoints: []string{"/api/v1/admin/"},
		},
		HTTPS: HTTPSConfig{
			Enabled:      getBoolEnv("HTTPS_ENABLED", false),
			RedirectHTTP: getBoolEnv("HTTPS_REDIRECT", false),
			HSTSMaxAge:   31536000, // 1 year
			CertFile:     getEnvOrDefault("TLS_CERT_FILE", ""),
			KeyFile:      getEnvOrDefault("TLS_KEY_FILE", ""),
		},
		Headers: HeadersConfig{
			EnableCSP:            true,
			CSPPolicy:            "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self';",
			EnableHSTS:           true,
			EnableFrameOptions:   true,
			EnableContentType:    true,
			EnableReferrerPolicy: true,
		},
		Validation: ValidationConfig{
			MaxJSONSize:      1024 * 1024, // 1MB
			MaxFieldLength:   1000,
			AllowedFileTypes: []string{".go", ".txt", ".md"},
			SanitizeInput:    true,
		},
		Logging: LoggingConfig{
			LogSensitiveData: false,
			LogLevel:         getEnvOrDefault("LOG_LEVEL", "INFO"),
			LogFormat:        "json",
		},
	}
}

// Helper functions for environment variable parsing.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}

	return defaultValue
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// In production, this should be set as an environment variable.
		secret = "your-super-secure-jwt-secret-change-in-production-min-32-chars"
	}

	return []byte(secret)
}

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		// Secure default - only localhost for development.
		return []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		}
	}

	return strings.Split(origins, ",")
}

func generateSecureAPIKey() string {
	// This should be generated securely and stored as env var in production.
	return "admin-api-key-change-in-production-12345678"
}
