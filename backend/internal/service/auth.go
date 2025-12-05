// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package service

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"go-pro-backend/internal/domain"
	apierrors "go-pro-backend/internal/errors"
	"go-pro-backend/internal/repository"
)

// AuthService defines business logic for authentication and user management.
type AuthService interface {
	// Initialize Firebase Admin SDK
	Initialize(ctx context.Context) error

	// Verify Firebase token and return claims
	VerifyToken(ctx context.Context, idToken string) (*domain.FirebaseClaims, error)

	// Get or create user from Firebase token
	GetOrCreateUser(ctx context.Context, firebaseUID, email, displayName, photoURL string) (*domain.User, error)

	// Verify token and sync user (main endpoint logic)
	VerifyAndSyncUser(ctx context.Context, idToken string) (*domain.VerifyTokenResponse, error)

	// Get user profile
	GetUserProfile(ctx context.Context, userID string) (*domain.UserProfileResponse, error)

	// Update user role (admin only)
	UpdateUserRole(ctx context.Context, adminUserID, targetUserID string, role domain.UserRole) error

	// Update user profile
	UpdateUserProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.User, error)
}

// firebaseAuthService implements AuthService using Firebase Admin SDK.
type firebaseAuthService struct {
	repo         repository.UserRepository
	config       *Config
	firebaseAuth *auth.Client
	initOnce     sync.Once
	initErr      error
}

// NewAuthService creates a new Firebase-based authentication service.
func NewAuthService(userRepo repository.UserRepository, config *Config) AuthService {
	return &firebaseAuthService{
		repo:   userRepo,
		config: config,
	}
}

// Initialize initializes the Firebase Admin SDK.
func (s *firebaseAuthService) Initialize(ctx context.Context) error {
	s.initOnce.Do(func() {
		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
		devMode := os.Getenv("DEV_MODE") == "true"

		if projectID == "" {
			if devMode {
				s.config.Logger.Warn(ctx, "Firebase not configured (DEV_MODE=true), skipping Firebase initialization")
				return
			}
			s.initErr = fmt.Errorf("FIREBASE_PROJECT_ID not set")
			return
		}

		// Initialize Firebase app
		var opt option.ClientOption
		if credentialsPath != "" {
			opt = option.WithCredentialsFile(credentialsPath)
		} else {
			// Use default credentials (for Cloud Run/GCE)
			opt = option.WithCredentialsFile("")
		}

		config := &firebase.Config{
			ProjectID: projectID,
		}

		app, err := firebase.NewApp(ctx, config, opt)
		if err != nil {
			s.initErr = fmt.Errorf("error initializing Firebase app: %w", err)
			return
		}

		// Initialize Auth client
		authClient, err := app.Auth(ctx)
		if err != nil {
			s.initErr = fmt.Errorf("error initializing Firebase Auth client: %w", err)
			return
		}

		s.firebaseAuth = authClient
		s.config.Logger.Info(ctx, "Firebase Admin SDK initialized successfully")
	})

	return s.initErr
}

// VerifyToken verifies a Firebase ID token and returns the claims.
func (s *firebaseAuthService) VerifyToken(ctx context.Context, idToken string) (*domain.FirebaseClaims, error) {
	devMode := os.Getenv("DEV_MODE") == "true"

	if s.firebaseAuth == nil {
		if err := s.Initialize(ctx); err != nil {
			if devMode && os.Getenv("FIREBASE_PROJECT_ID") == "" {
				// In dev mode without Firebase, create a mock token response
				s.config.Logger.Warn(ctx, "DEV_MODE: Using mock token verification")
				return &domain.FirebaseClaims{
					UserID:    "dev-user-id",
					Email:     "developer@local.dev",
					DisplayName: "Developer",
					Picture:    "",
				}, nil
			}
			return nil, apierrors.NewInternalError("Firebase not initialized", err)
		}
	}

	// Skip Firebase verification in dev mode if Firebase is not initialized
	if s.firebaseAuth == nil && devMode {
		s.config.Logger.Warn(ctx, "DEV_MODE: Using mock token verification")
		return &domain.FirebaseClaims{
			UserID:    "dev-user-id",
			Email:     "developer@local.dev",
			DisplayName: "Developer",
			Picture:    "",
		}, nil
	}

	// Verify the ID token
	token, err := s.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		s.config.Logger.Error(ctx, "Failed to verify Firebase token", "error", err)
		return nil, apierrors.NewUnauthorizedError("Invalid or expired Firebase token")
	}

	// Extract claims
	claims := &domain.FirebaseClaims{
		UserID:    token.UID,
		IssuedAt:  time.Unix(token.IssuedAt, 0),
		ExpiresAt: time.Unix(token.Expires, 0),
	}

	// Extract email
	if email, ok := token.Claims["email"].(string); ok {
		claims.Email = email
	}

	// Extract display name
	if name, ok := token.Claims["name"].(string); ok {
		claims.DisplayName = name
	}

	// Extract picture URL
	if picture, ok := token.Claims["picture"].(string); ok {
		claims.Picture = picture
	}

	return claims, nil
}

// GetOrCreateUser gets or creates a user in the database from Firebase information.
func (s *firebaseAuthService) GetOrCreateUser(ctx context.Context, firebaseUID, email, displayName, photoURL string) (*domain.User, error) {
	// Try to find existing user by Firebase UID
	user, err := s.repo.GetByFirebaseUID(ctx, firebaseUID)
	if err == nil {
		// User exists, update last login
		if updateErr := s.repo.UpdateLastLogin(ctx, user.ID); updateErr != nil {
			s.config.Logger.Warn(ctx, "Failed to update last login", "user_id", user.ID, "error", updateErr)
		}
		return user, nil
	}

	// Check if this is the first user (should become admin)
	_, totalUsers, err := s.repo.GetAll(ctx, &domain.PaginationRequest{Page: 1, PageSize: 1})
	if err != nil {
		return nil, apierrors.NewInternalError("Failed to check user count", err)
	}

	// Determine role
	role := domain.RoleStudent
	if totalUsers == 0 {
		role = domain.RoleAdmin
		s.config.Logger.Info(ctx, "First user registration, assigning admin role", "email", email)
	}

	// Generate username from email (before @)
	username := email
	if len(email) > 0 {
		for i, c := range email {
			if c == '@' {
				username = email[:i]
				break
			}
		}
	}

	// Create new user
	now := time.Now()
	newUser := &domain.User{
		ID:          generateID("user"),
		FirebaseUID: firebaseUID,
		Email:       email,
		Username:    username,
		DisplayName: displayName,
		PhotoURL:    photoURL,
		Role:        role,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &now,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, apierrors.NewInternalError("Failed to create user", err)
	}

	s.config.Logger.Info(ctx, "Created new user", "user_id", newUser.ID, "email", email, "role", role)
	return newUser, nil
}

// VerifyAndSyncUser verifies a Firebase token and syncs the user to the database.
func (s *firebaseAuthService) VerifyAndSyncUser(ctx context.Context, idToken string) (*domain.VerifyTokenResponse, error) {
	// Verify the Firebase token
	claims, err := s.VerifyToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	// Check if user exists by Firebase UID
	existingUser, err := s.repo.GetByFirebaseUID(ctx, claims.UserID)
	isNewUser := err != nil

	// Get or create user
	user, err := s.GetOrCreateUser(ctx, claims.UserID, claims.Email, claims.DisplayName, claims.Picture)
	if err != nil {
		return nil, err
	}

	// If existing user, update display name and photo if changed
	if !isNewUser && existingUser != nil {
		updated := false

		if claims.DisplayName != "" && claims.DisplayName != existingUser.DisplayName {
			existingUser.DisplayName = claims.DisplayName
			updated = true
		}

		if claims.Picture != "" && claims.Picture != existingUser.PhotoURL {
			existingUser.PhotoURL = claims.Picture
			updated = true
		}

		if updated {
			existingUser.UpdatedAt = time.Now()
			if err := s.repo.Update(ctx, existingUser); err != nil {
				s.config.Logger.Warn(ctx, "Failed to update user profile from Firebase", "user_id", existingUser.ID, "error", err)
			} else {
				user = existingUser
			}
		}
	}

	return &domain.VerifyTokenResponse{
		User:          user.ToProfileResponse(),
		IsNewUser:     isNewUser,
		TokenVerified: true,
	}, nil
}

// GetUserProfile retrieves a user's profile.
func (s *firebaseAuthService) GetUserProfile(ctx context.Context, userID string) (*domain.UserProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, apierrors.NewNotFoundError("User not found")
	}

	return user.ToProfileResponse(), nil
}

// UpdateUserRole updates a user's role (admin only).
func (s *firebaseAuthService) UpdateUserRole(ctx context.Context, adminUserID, targetUserID string, role domain.UserRole) error {
	// Verify admin user
	adminUser, err := s.repo.GetByID(ctx, adminUserID)
	if err != nil {
		return apierrors.NewUnauthorizedError("Admin user not found")
	}

	if adminUser.Role != domain.RoleAdmin {
		return apierrors.NewForbiddenError("Only admins can update user roles")
	}

	// Get target user
	targetUser, err := s.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return apierrors.NewNotFoundError("Target user not found")
	}

	// Validate role
	if !role.IsValid() {
		return apierrors.NewBadRequestError("Invalid role")
	}

	// Prevent admin from demoting themselves
	if adminUserID == targetUserID && role != domain.RoleAdmin {
		return apierrors.NewBadRequestError("Admins cannot demote themselves")
	}

	// Update role
	targetUser.Role = role
	targetUser.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, targetUser); err != nil {
		return apierrors.NewInternalError("Failed to update user role", err)
	}

	s.config.Logger.Info(ctx, "User role updated", "admin_id", adminUserID, "target_user_id", targetUserID, "new_role", role)
	return nil
}

// UpdateUserProfile updates a user's profile information.
func (s *firebaseAuthService) UpdateUserProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, apierrors.NewNotFoundError("User not found")
	}

	// Validate request
	if err := s.config.Validator.Validate(req); err != nil {
		return nil, apierrors.NewValidationError("Invalid update request", err)
	}

	// Update fields
	updated := false

	if req.DisplayName != nil && *req.DisplayName != user.DisplayName {
		user.DisplayName = *req.DisplayName
		updated = true
	}

	if req.PhotoURL != nil && *req.PhotoURL != user.PhotoURL {
		user.PhotoURL = *req.PhotoURL
		updated = true
	}

	if !updated {
		return user, nil
	}

	// Save changes
	user.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, apierrors.NewInternalError("Failed to update user profile", err)
	}

	s.config.Logger.Info(ctx, "User profile updated", "user_id", userID)
	return user, nil
}
