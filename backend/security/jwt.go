// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package security provides authentication, authorization, and security middleware.
package security

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Custom JWT claims.
type Claims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations.
type JWTManager struct {
	config JWTConfig
}

// NewJWTManager creates a new JWT manager.
func NewJWTManager(config JWTConfig) *JWTManager {
	return &JWTManager{
		config: config,
	}
}

// GenerateTokens generates both access and refresh tokens.
func (j *JWTManager) GenerateTokens(userID, email string, roles []string) (accessToken, refreshToken string, err error) {
	// Generate access token.
	accessToken, err = j.generateToken(userID, email, roles, "access", j.config.AccessTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token.
	refreshToken, err = j.generateToken(userID, email, roles, "refresh", j.config.RefreshTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateToken creates a JWT token with the specified parameters.
func (j *JWTManager) generateToken(userID, email string, roles []string, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        generateTokenID(),
			Subject:   userID,
			Audience:  []string{j.config.Audience},
			Issuer:    j.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.config.Secret)
}

// ValidateToken validates a JWT token and returns the claims.
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// Remove Bearer prefix if present.
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.config.Secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Additional validation.
	if err := j.validateClaims(claims); err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return claims, nil
}

// validateClaims performs additional validation on the claims.
func (j *JWTManager) validateClaims(claims *Claims) error {
	// Check if token is expired.
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return errors.New("token has expired")
	}

	// Check if token is not valid yet.
	if claims.NotBefore != nil && claims.NotBefore.After(time.Now()) {
		return errors.New("token not valid yet")
	}

	// Check issuer.
	if claims.Issuer != j.config.Issuer {
		return errors.New("invalid token issuer")
	}

	// Check audience.
	if len(claims.Audience) == 0 || claims.Audience[0] != j.config.Audience {
		return errors.New("invalid token audience")
	}

	// Validate user ID.
	if claims.UserID == "" {
		return errors.New("missing user ID in token")
	}

	// Validate token type.
	if claims.TokenType != "access" && claims.TokenType != "refresh" {
		return errors.New("invalid token type")
	}

	return nil
}

// RefreshTokens validates a refresh token and generates new tokens.
func (j *JWTManager) RefreshTokens(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Verify this is actually a refresh token.
	if claims.TokenType != "refresh" {
		return "", "", errors.New("provided token is not a refresh token")
	}

	// Generate new tokens.
	return j.GenerateTokens(claims.UserID, claims.Email, claims.Roles)
}

// ExtractTokenFromHeader extracts JWT token from Authorization header.
func ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Remove "Bearer " prefix.
	const bearerPrefix = "Bearer "
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return authHeader[len(bearerPrefix):]
	}

	return authHeader
}

// HasRole checks if the user has a specific role.
func (c *Claims) HasRole(role string) bool {
	for _, userRole := range c.Roles {
		if userRole == role {
			return true
		}
	}

	return false
}

// HasAnyRole checks if the user has any of the specified roles.
func (c *Claims) HasAnyRole(roles []string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}

	return false
}

// IsAccessToken checks if the token is an access token.
func (c *Claims) IsAccessToken() bool {
	return c.TokenType == "access"
}

// IsRefreshToken checks if the token is a refresh token.
func (c *Claims) IsRefreshToken() bool {
	return c.TokenType == "refresh"
}

// generateTokenID creates a unique token ID.
func generateTokenID() string {
	return fmt.Sprintf("tkn_%d_%d", time.Now().UnixNano(), time.Now().Unix())
}
