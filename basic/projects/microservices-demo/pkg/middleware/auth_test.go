package middleware

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	userID := "test-user-id"
	username := "testuser"
	email := "test@example.com"

	token, err := GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}
}

func TestValidateToken(t *testing.T) {
	userID := "test-user-id"
	username := "testuser"
	email := "test@example.com"

	// Generate a token
	token, err := GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check claims
	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestValidateInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.here"

	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestTokenExpiration(t *testing.T) {
	// This test would require mocking time or waiting 24 hours
	// For now, we just verify the expiration is set correctly
	userID := "test-user-id"
	username := "testuser"
	email := "test@example.com"

	token, err := GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check expiration is in the future
	expirationTime := time.Unix(claims.ExpiresAt.Unix(), 0)
	if expirationTime.Before(time.Now()) {
		t.Error("Token expiration is in the past")
	}

	// Check expiration is approximately 24 hours from now
	expectedExpiration := time.Now().Add(24 * time.Hour)
	diff := expirationTime.Sub(expectedExpiration)
	if diff < -time.Minute || diff > time.Minute {
		t.Errorf("Token expiration is not approximately 24 hours from now")
	}
}

