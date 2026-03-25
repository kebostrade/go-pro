// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-pro-backend/internal/domain"
)

func TestMemoryProgressRepository_Create(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()
	progress := &domain.Progress{
		ID:        "progress-1",
		UserID:    "user-1",
		LessonID:  "lesson-1",
		Status:    domain.StatusInProgress,
	}

	err := repo.Create(ctx, progress)
	if err != nil {
		t.Fatalf("Failed to create progress: %v", err)
	}

	// Verify progress was created
	retrieved, err := repo.GetByID(ctx, "progress-1")
	if err != nil {
		t.Fatalf("Failed to get progress: %v", err)
	}

	if retrieved.ID != "progress-1" {
		t.Errorf("Expected ID progress-1, got %s", retrieved.ID)
	}
	if retrieved.UserID != "user-1" {
		t.Errorf("Expected UserID user-1, got %s", retrieved.UserID)
	}
	if retrieved.Status != domain.StatusInProgress {
		t.Errorf("Expected Status in_progress, got %s", retrieved.Status)
	}
}

func TestMemoryProgressRepository_GetByUserAndLesson(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()
	progress := &domain.Progress{
		ID:     "progress-2",
		UserID: "user-1",
		LessonID: "lesson-2",
		Status: domain.StatusCompleted,
	}

	err := repo.Create(ctx, progress)
	if err != nil {
		t.Fatalf("Failed to create progress: %v", err)
	}

	// Get by user and lesson
	retrieved, err := repo.GetByUserAndLesson(ctx, "user-1", "lesson-2")
	if err != nil {
		t.Fatalf("Failed to get progress by user and lesson: %v", err)
	}

	if retrieved.ID != "progress-2" {
		t.Errorf("Expected ID progress-2, got %s", retrieved.ID)
	}
	if retrieved.Status != domain.StatusCompleted {
		t.Errorf("Expected Status completed, got %s", retrieved.Status)
	}

	// Test not found case
	_, err = repo.GetByUserAndLesson(ctx, "user-2", "lesson-2")
	if err == nil {
		t.Error("Expected error for non-existent progress, got nil")
	}
}

func TestMemoryProgressRepository_Update(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()
	progress := &domain.Progress{
		ID:       "progress-3",
		UserID:   "user-1",
		LessonID: "lesson-3",
		Status:   domain.StatusInProgress,
	}

	err := repo.Create(ctx, progress)
	if err != nil {
		t.Fatalf("Failed to create progress: %v", err)
	}

	// Update progress
	progress.Status = domain.StatusCompleted

	err = repo.Update(ctx, progress)
	if err != nil {
		t.Fatalf("Failed to update progress: %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByID(ctx, "progress-3")
	if err != nil {
		t.Fatalf("Failed to get progress: %v", err)
	}

	if retrieved.Status != domain.StatusCompleted {
		t.Errorf("Expected Status completed, got %s", retrieved.Status)
	}
}

func TestMemoryProgressRepository_GetByUserID(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()

	// Create multiple progress records for same user
	progress1 := &domain.Progress{
		ID:       "progress-4",
		UserID:   "user-2",
		LessonID: "lesson-1",
		Status:   domain.StatusCompleted,
	}
	progress2 := &domain.Progress{
		ID:       "progress-5",
		UserID:   "user-2",
		LessonID: "lesson-2",
		Status:   domain.StatusInProgress,
	}

	err := repo.Create(ctx, progress1)
	if err != nil {
		t.Fatalf("Failed to create progress1: %v", err)
	}
	err = repo.Create(ctx, progress2)
	if err != nil {
		t.Fatalf("Failed to create progress2: %v", err)
	}

	// Get all progress for user
	progresses, total, err := repo.GetByUserID(ctx, "user-2", &domain.PaginationRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to get progress by user ID: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
	if len(progresses) != 2 {
		t.Errorf("Expected 2 progress records, got %d", len(progresses))
	}
}

func TestMemoryProgressRepository_Pagination(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()

	// Create multiple progress records
	for i := 1; i <= 15; i++ {
		progress := &domain.Progress{
			ID:       fmt.Sprintf("progress-%d", i),
			UserID:   "user-pagination",
			LessonID: fmt.Sprintf("lesson-%d", i),
			Status:   domain.StatusInProgress,
		}
		if err := repo.Create(ctx, progress); err != nil {
			t.Fatalf("Failed to create progress %d: %v", i, err)
		}
	}

	// Test first page
	progresses, total, err := repo.GetByUserID(ctx, "user-pagination", &domain.PaginationRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to get first page: %v", err)
	}

	if total != 15 {
		t.Errorf("Expected total 15, got %d", total)
	}
	if len(progresses) != 10 {
		t.Errorf("Expected 10 records on first page, got %d", len(progresses))
	}

	// Test second page
	progresses, total, err = repo.GetByUserID(ctx, "user-pagination", &domain.PaginationRequest{
		Page:     2,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to get second page: %v", err)
	}

	if len(progresses) != 5 {
		t.Errorf("Expected 5 records on second page, got %d", len(progresses))
	}
}

func TestMemoryProgressRepository_Delete(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()
	progress := &domain.Progress{
		ID:       "progress-6",
		UserID:   "user-1",
		LessonID: "lesson-1",
		Status:   domain.StatusInProgress,
	}

	err := repo.Create(ctx, progress)
	if err != nil {
		t.Fatalf("Failed to create progress: %v", err)
	}

	// Delete progress
	err = repo.Delete(ctx, "progress-6")
	if err != nil {
		t.Fatalf("Failed to delete progress: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, "progress-6")
	if err == nil {
		t.Error("Expected error for deleted progress, got nil")
	}

	// Verify user-lesson index was also removed
	_, err = repo.GetByUserAndLesson(ctx, "user-1", "lesson-1")
	if err == nil {
		t.Error("Expected error for deleted progress via user-lesson, got nil")
	}
}

func TestMemoryProgressRepository_DuplicateUserLesson(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()

	progress1 := &domain.Progress{
		ID:       "progress-7",
		UserID:   "user-1",
		LessonID: "lesson-1",
		Status:   domain.StatusInProgress,
	}

	err := repo.Create(ctx, progress1)
	if err != nil {
		t.Fatalf("Failed to create progress1: %v", err)
	}

	// Try to create duplicate progress for same user and lesson
	progress2 := &domain.Progress{
		ID:       "progress-8",
		UserID:   "user-1",
		LessonID: "lesson-1",
		Status:   domain.StatusInProgress,
	}

	err = repo.Create(ctx, progress2)
	if err == nil {
		t.Error("Expected error for duplicate user-lesson combination, got nil")
	}
}

func TestMemoryProgressRepository_CompletedLesson(t *testing.T) {
	repo := NewMemoryProgressRepository()

	ctx := context.Background()
	now := time.Now()
	completedAt := now

	progress := &domain.Progress{
		ID:          "progress-9",
		UserID:      "user-1",
		LessonID:    "lesson-1",
		Status:      domain.StatusCompleted,
		CompletedAt: &completedAt,
	}

	err := repo.Create(ctx, progress)
	if err != nil {
		t.Fatalf("Failed to create completed progress: %v", err)
	}

	// Verify all fields
	retrieved, err := repo.GetByID(ctx, "progress-9")
	if err != nil {
		t.Fatalf("Failed to get progress: %v", err)
	}

	if retrieved.Status != domain.StatusCompleted {
		t.Errorf("Expected Status completed, got %s", retrieved.Status)
	}
	if retrieved.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set, got nil")
	}
}
