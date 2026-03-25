// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package repository

import (
	"context"
	"testing"

	"go-pro-backend/internal/domain"
)

func TestMemoryCourseRepository_CreateAndGet(t *testing.T) {
	repo := NewMemoryCourseRepository()

	ctx := context.Background()
	course := &domain.Course{
		ID:          "test-course-1",
		Title:       "Test Course",
		Description: "A test course for unit testing",
	}

	err := repo.Create(ctx, course)
	if err != nil {
		t.Fatalf("Failed to create course: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, "test-course-1")
	if err != nil {
		t.Fatalf("Failed to get course: %v", err)
	}

	if retrieved.ID != "test-course-1" {
		t.Errorf("Expected ID test-course-1, got %s", retrieved.ID)
	}
	if retrieved.Title != "Test Course" {
		t.Errorf("Expected title 'Test Course', got %s", retrieved.Title)
	}
}

func TestMemoryUserRepository_CreateAndGet(t *testing.T) {
	repo := NewMemoryUserRepository()

	ctx := context.Background()
	user := &domain.User{
		ID:          "test-user-1",
		FirebaseUID: "firebase-123",
		Username:    "testuser",
		Email:       "test@example.com",
		Role:        domain.RoleStudent,
		IsActive:    true,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, "test-user-1")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrieved.ID != "test-user-1" {
		t.Errorf("Expected ID test-user-1, got %s", retrieved.ID)
	}
	if retrieved.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", retrieved.Email)
	}
}

func TestMemoryUserRepository_GetByFirebaseUID(t *testing.T) {
	repo := NewMemoryUserRepository()

	ctx := context.Background()
	user := &domain.User{
		ID:          "test-user-2",
		FirebaseUID: "firebase-456",
		Username:    "testuser2",
		Email:       "test2@example.com",
		Role:        domain.RoleStudent,
		IsActive:    true,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := repo.GetByFirebaseUID(ctx, "firebase-456")
	if err != nil {
		t.Fatalf("Failed to get user by Firebase UID: %v", err)
	}

	if retrieved.ID != "test-user-2" {
		t.Errorf("Expected ID test-user-2, got %s", retrieved.ID)
	}
}
