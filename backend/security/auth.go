// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package security provides authentication, authorization, and security middleware.
package security

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never serialize password hash
	Roles        []string   `json:"roles"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	IsActive     bool       `json:"is_active"`
	IsVerified   bool       `json:"is_verified"`
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents a registration request.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenResponse represents a token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthService handles authentication operations.
type AuthService struct {
	jwtManager *JWTManager
	userStore  UserStore // In production, this would be a database
}

// UserStore interface for user storage operations.
type UserStore interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
}

// InMemoryUserStore is a simple in-memory implementation of UserStore.
type InMemoryUserStore struct {
	users  map[string]*User
	emails map[string]string // email -> id mapping
}

// NewInMemoryUserStore creates a new in-memory user store.
func NewInMemoryUserStore() *InMemoryUserStore {
	store := &InMemoryUserStore{
		users:  make(map[string]*User),
		emails: make(map[string]string),
	}

	// Create a default admin user for testing.
	adminUser := &User{
		ID:           "admin-user-001",
		Email:        "admin@gopro.dev",
		PasswordHash: hashPassword("admin123"),
		Roles:        []string{"admin", "user"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
		IsVerified:   true,
	}

	// Create a default regular user for testing.
	regularUser := &User{
		ID:           "demo-user-001",
		Email:        "demo@gopro.dev",
		PasswordHash: hashPassword("demo123"),
		Roles:        []string{"user"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
		IsVerified:   true,
	}

	if err := store.CreateUser(adminUser); err != nil {
		// Log error but continue - this is for testing only
		_ = err
	}
	if err := store.CreateUser(regularUser); err != nil {
		// Log error but continue - this is for testing only
		_ = err
	}

	return store
}

func (s *InMemoryUserStore) CreateUser(user *User) error {
	// Check if email already exists.
	if _, exists := s.emails[user.Email]; exists {
		return errors.New("email already exists")
	}

	s.users[user.ID] = user
	s.emails[user.Email] = user.ID

	return nil
}

func (s *InMemoryUserStore) GetUserByEmail(email string) (*User, error) {
	userID, exists := s.emails[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *InMemoryUserStore) GetUserByID(id string) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *InMemoryUserStore) UpdateUser(user *User) error {
	_, exists := s.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	user.UpdatedAt = time.Now()
	s.users[user.ID] = user

	return nil
}

func (s *InMemoryUserStore) DeleteUser(id string) error {
	user, exists := s.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
	delete(s.emails, user.Email)

	return nil
}

// NewAuthService creates a new authentication service.
func NewAuthService(jwtManager *JWTManager, userStore UserStore) *AuthService {
	return &AuthService{
		jwtManager: jwtManager,
		userStore:  userStore,
	}
}

// Register creates a new user account.
func (as *AuthService) Register(req RegisterRequest) (*User, error) {
	// Validate input.
	if err := as.validateRegistration(req); err != nil {
		return nil, err
	}

	// Check if user already exists.
	if _, err := as.userStore.GetUserByEmail(req.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password.
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user.
	user := &User{
		ID:           generateUserID(),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(passwordHash),
		Roles:        []string{"user"}, // Default role
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
		IsVerified:   false, // In production, implement email verification
	}

	if err := as.userStore.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns tokens.
func (as *AuthService) Login(req LoginRequest) (*TokenResponse, error) {
	// Validate input.
	if err := as.validateLogin(req); err != nil {
		return nil, err
	}

	// Get user.
	user, err := as.userStore.GetUserByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active.
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login.
	now := time.Now()
	user.LastLoginAt = &now
	if err := as.userStore.UpdateUser(user); err != nil {
		// Log error but don't fail login
		_ = err
	}

	// Generate tokens.
	accessToken, refreshToken, err := as.jwtManager.GenerateTokens(user.ID, user.Email, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(as.jwtManager.config.AccessTokenTTL.Seconds()),
	}, nil
}

// RefreshToken generates new tokens using a refresh token.
func (as *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	accessToken, newRefreshToken, err := as.jwtManager.RefreshTokens(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(as.jwtManager.config.AccessTokenTTL.Seconds()),
	}, nil
}

// GetCurrentUser returns the current authenticated user.
func (as *AuthService) GetCurrentUser(userID string) (*User, error) {
	user, err := as.userStore.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	return user, nil
}

// ChangePassword changes a user's password.
func (as *AuthService) ChangePassword(userID, currentPassword, newPassword string) error {
	user, err := as.userStore.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password.
	if err := as.validatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password.
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update user.
	user.PasswordHash = string(newPasswordHash)
	user.UpdatedAt = time.Now()

	return as.userStore.UpdateUser(user)
}

// Validation functions.

func (as *AuthService) validateRegistration(req RegisterRequest) error {
	if err := as.validateEmail(req.Email); err != nil {
		return err
	}

	return as.validatePassword(req.Password)
}

func (as *AuthService) validateLogin(req LoginRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (as *AuthService) validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if len(email) > 254 {
		return errors.New("email is too long")
	}
	if !ValidateEmail(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (as *AuthService) validatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return errors.New("password is too long")
	}

	// Check for at least one uppercase, lowercase, and number.
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasNumber := strings.ContainsAny(password, "0123456789")

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}

// HTTP Handlers.

// RegisterHandler handles user registration.
func (as *AuthService) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := as.Register(req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
		"data": map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"roles":     user.Roles,
			"is_active": user.IsActive,
		},
	})
}

// LoginHandler handles user login.
func (as *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokenResponse, err := as.Login(req)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data":    tokenResponse,
	})
}

// RefreshHandler handles token refresh.
func (as *AuthService) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokenResponse, err := as.RefreshToken(req.RefreshToken)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Token refreshed successfully",
		"data":    tokenResponse,
	})
}

// ProfileHandler returns the current user's profile.
func (as *AuthService) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserContextKey).(string)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := as.GetCurrentUser(userID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// Utility functions.

func generateUserID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("usr_%d", time.Now().UnixNano())
	}

	return "usr_" + hex.EncodeToString(bytes)
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   false,
		"error":     message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
