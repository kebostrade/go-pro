// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package cache provides functionality for the GO-PRO Learning Platform.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisSessionStore implements session storage using Redis.
type RedisSessionStore struct {
	client *redis.Client
	prefix string
}

// SessionData represents session data structure.
type SessionData struct {
	UserID    string                 `json:"user_id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	Roles     []string               `json:"roles"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	ExpiresAt time.Time              `json:"expires_at"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
}

// NewRedisSessionStore creates a new Redis session store.
func NewRedisSessionStore(client *redis.Client, prefix string) *RedisSessionStore {
	if prefix == "" {
		prefix = "session:"
	}

	return &RedisSessionStore{
		client: client,
		prefix: prefix,
	}
}

// sessionKey generates a session key with prefix.
func (r *RedisSessionStore) sessionKey(sessionID string) string {
	return r.prefix + sessionID
}

// userSessionsKey generates a user sessions key.
func (r *RedisSessionStore) userSessionsKey(userID string) string {
	return r.prefix + "user:" + userID
}

// CreateSession creates a new session.
func (r *RedisSessionStore) CreateSession(ctx context.Context, sessionID string, data map[string]interface{}, expiration time.Duration) error {
	sessionData := &SessionData{
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expiration),
	}

	// Extract common fields from data.
	if userID, ok := data["user_id"].(string); ok {
		sessionData.UserID = userID
	}
	if username, ok := data["username"].(string); ok {
		sessionData.Username = username
	}
	if email, ok := data["email"].(string); ok {
		sessionData.Email = email
	}
	if roles, ok := data["roles"].([]string); ok {
		sessionData.Roles = roles
	}
	if ipAddress, ok := data["ip_address"].(string); ok {
		sessionData.IPAddress = ipAddress
	}
	if userAgent, ok := data["user_agent"].(string); ok {
		sessionData.UserAgent = userAgent
	}

	// Serialize session data.
	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	// Use pipeline for atomic operations.
	pipe := r.client.TxPipeline()

	// Set session data.
	pipe.Set(ctx, r.sessionKey(sessionID), sessionJSON, expiration)

	// Add session to user's session list if user_id is present.
	if sessionData.UserID != "" {
		pipe.SAdd(ctx, r.userSessionsKey(sessionData.UserID), sessionID)
		pipe.Expire(ctx, r.userSessionsKey(sessionData.UserID), expiration)
	}

	// Execute pipeline.
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSession retrieves session data.
func (r *RedisSessionStore) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	sessionJSON, err := r.client.Get(ctx, r.sessionKey(sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheNotFound
		}

		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var sessionData SessionData
	if err := json.Unmarshal([]byte(sessionJSON), &sessionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	// Check if session is expired.
	if time.Now().After(sessionData.ExpiresAt) {
		// Clean up expired session.
		r.DeleteSession(ctx, sessionID)
		return nil, ErrCacheNotFound
	}

	// Return session data.
	result := make(map[string]interface{})
	result["user_id"] = sessionData.UserID
	result["username"] = sessionData.Username
	result["email"] = sessionData.Email
	result["roles"] = sessionData.Roles
	result["created_at"] = sessionData.CreatedAt
	result["updated_at"] = sessionData.UpdatedAt
	result["expires_at"] = sessionData.ExpiresAt
	result["ip_address"] = sessionData.IPAddress
	result["user_agent"] = sessionData.UserAgent

	// Add custom data.
	for k, v := range sessionData.Data {
		result[k] = v
	}

	return result, nil
}

// UpdateSession updates existing session data.
func (r *RedisSessionStore) UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	// Get existing session.
	existingData, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get existing session: %w", err)
	}

	// Merge data.
	for k, v := range data {
		existingData[k] = v
	}
	existingData["updated_at"] = time.Now()

	// Get current TTL.
	ttl, err := r.client.TTL(ctx, r.sessionKey(sessionID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get session TTL: %w", err)
	}

	// Create updated session.
	return r.CreateSession(ctx, sessionID, existingData, ttl)
}

// DeleteSession removes a session.
func (r *RedisSessionStore) DeleteSession(ctx context.Context, sessionID string) error {
	// Get session to find user ID.
	sessionData, err := r.GetSession(ctx, sessionID)
	if err != nil && err != ErrCacheNotFound {
		return fmt.Errorf("failed to get session for deletion: %w", err)
	}

	// Use pipeline for atomic operations.
	pipe := r.client.TxPipeline()

	// Delete session.
	pipe.Del(ctx, r.sessionKey(sessionID))

	// Remove from user's session list if user_id is present.
	if sessionData != nil {
		if userID, ok := sessionData["user_id"].(string); ok && userID != "" {
			pipe.SRem(ctx, r.userSessionsKey(userID), sessionID)
		}
	}

	// Execute pipeline.
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// RefreshSession extends session expiration.
func (r *RedisSessionStore) RefreshSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	// Check if session exists.
	exists, err := r.client.Exists(ctx, r.sessionKey(sessionID)).Result()
	if err != nil {
		return fmt.Errorf("failed to check session existence: %w", err)
	}
	if exists == 0 {
		return ErrCacheNotFound
	}

	// Update expiration.
	err = r.client.Expire(ctx, r.sessionKey(sessionID), expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	// Update session data with new expiration time.
	sessionData, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session for refresh: %w", err)
	}

	sessionData["expires_at"] = time.Now().Add(expiration)
	sessionData["updated_at"] = time.Now()

	return r.UpdateSession(ctx, sessionID, sessionData)
}

// ListUserSessions returns all session IDs for a user.
func (r *RedisSessionStore) ListUserSessions(ctx context.Context, userID string) ([]string, error) {
	sessions, err := r.client.SMembers(ctx, r.userSessionsKey(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}

		return nil, fmt.Errorf("failed to list user sessions: %w", err)
	}

	// Filter out expired sessions.
	validSessions := make([]string, 0, len(sessions))
	for _, sessionID := range sessions {
		exists, err := r.client.Exists(ctx, r.sessionKey(sessionID)).Result()
		if err != nil {
			continue
		}
		if exists > 0 {
			validSessions = append(validSessions, sessionID)
		} else {
			// Remove expired session from user's session list.
			r.client.SRem(ctx, r.userSessionsKey(userID), sessionID)
		}
	}

	return validSessions, nil
}

// DeleteUserSessions removes all sessions for a user.
func (r *RedisSessionStore) DeleteUserSessions(ctx context.Context, userID string) error {
	// Get all user sessions.
	sessions, err := r.ListUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to list user sessions: %w", err)
	}

	if len(sessions) == 0 {
		return nil
	}

	// Use pipeline for batch deletion.
	pipe := r.client.TxPipeline()

	// Delete all session keys.
	sessionKeys := make([]string, len(sessions))
	for i, sessionID := range sessions {
		sessionKeys[i] = r.sessionKey(sessionID)
	}
	pipe.Del(ctx, sessionKeys...)

	// Delete user sessions set.
	pipe.Del(ctx, r.userSessionsKey(userID))

	// Execute pipeline.
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

// CleanupExpiredSessions removes expired sessions (should be run periodically).
func (r *RedisSessionStore) CleanupExpiredSessions(ctx context.Context) error {
	// Use SCAN to iterate through all session keys.
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, r.prefix+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan session keys: %w", err)
		}

		// Check each key's TTL.
		for _, key := range keys {
			ttl, err := r.client.TTL(ctx, key).Result()
			if err != nil {
				continue
			}

			// If TTL is -1 (no expiration) or -2 (key doesn't exist), skip.
			if ttl == -1 || ttl == -2 {
				continue
			}

			// If TTL is very small (less than 1 second), consider it expired.
			if ttl < time.Second {
				r.client.Del(ctx, key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// GetSessionCount returns the total number of active sessions.
func (r *RedisSessionStore) GetSessionCount(ctx context.Context) (int64, error) {
	keys, err := r.client.Keys(ctx, r.prefix+"*").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	// Filter out user session sets (they contain "user:" in the key)
	count := int64(0)
	for _, key := range keys {
		if !contains(key, "user:") {
			count++
		}
	}

	return count, nil
}

// GetUserSessionCount returns the number of active sessions for a user.
func (r *RedisSessionStore) GetUserSessionCount(ctx context.Context, userID string) (int64, error) {
	sessions, err := r.ListUserSessions(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user session count: %w", err)
	}

	return int64(len(sessions)), nil
}

// Helper function to check if string contains substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			(len(s) > len(substr) && s[1:len(substr)+1] == substr))))
}
