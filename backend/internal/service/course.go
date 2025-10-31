// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/internal/repository"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// courseService implements the CourseService interface.
type courseService struct {
	repo      repository.CourseRepository
	logger    logger.Logger
	validator validator.Validator
	cache     cache.CacheManager
	messaging *messaging.Service
}

// NewCourseService creates a new course service.
func NewCourseService(
	repo repository.CourseRepository,
	config *Config,
) CourseService {
	return &courseService{
		repo:      repo,
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}
}

// CreateCourse creates a new course.
func (s *courseService) CreateCourse(ctx context.Context, req *domain.CreateCourseRequest) (*domain.Course, error) {
	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		logger.LogError(s.logger, ctx, err, "validation failed for create course request")
		return nil, err
	}

	// Generate course ID from title.
	courseID := generateSlug(req.Title)

	// Check if course with this ID already exists.
	existingCourse, err := s.repo.GetByID(ctx, courseID)
	if err == nil && existingCourse != nil {
		return nil, errors.NewConflictError(fmt.Sprintf("course with title '%s' already exists", req.Title))
	}

	// Create course.
	course := &domain.Course{
		ID:          courseID,
		Title:       req.Title,
		Description: req.Description,
		Lessons:     []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository.
	if err := s.repo.Create(ctx, course); err != nil {
		logger.LogError(s.logger, ctx, err, "failed to create course", "course_id", courseID)
		return nil, errors.NewInternalError("failed to create course", err)
	}

	s.logger.Info(ctx, "course created successfully",
		"course_id", course.ID,
		"title", course.Title,
	)

	return course, nil
}

// GetCourseByID retrieves a course by ID.
func (s *courseService) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("course ID is required")
	}

	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if apiErr, ok := errors.IsAPIError(err); ok {
			return nil, apiErr
		}
		logger.LogError(s.logger, ctx, err, "failed to get course", "course_id", id)

		return nil, errors.NewInternalError("failed to retrieve course", err)
	}

	s.logger.Debug(ctx, "course retrieved successfully", "course_id", id)

	return course, nil
}

// GetAllCourses retrieves all courses with pagination.
func (s *courseService) GetAllCourses(ctx context.Context, pagination *domain.PaginationRequest) (*domain.ListResponse, error) {
	// Set default pagination if not provided.
	if pagination == nil {
		pagination = &domain.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}
	}

	courses, total, err := s.repo.GetAll(ctx, pagination)
	if err != nil {
		logger.LogError(s.logger, ctx, err, "failed to get courses")
		return nil, errors.NewInternalError("failed to retrieve courses", err)
	}

	// Calculate pagination metadata.
	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))
	hasNext := pagination.Page < totalPages
	hasPrev := pagination.Page > 1

	paginationResponse := &domain.PaginationResponse{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	response := &domain.ListResponse{
		Items:      courses,
		Pagination: paginationResponse,
	}

	s.logger.Debug(ctx, "courses retrieved successfully",
		"total_courses", total,
		"page", pagination.Page,
		"page_size", pagination.PageSize,
	)

	return response, nil
}

// UpdateCourse updates an existing course.
func (s *courseService) UpdateCourse(ctx context.Context, id string, req *domain.UpdateCourseRequest) (*domain.Course, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("course ID is required")
	}

	// Validate request.
	if err := s.validator.Validate(req); err != nil {
		logger.LogError(s.logger, ctx, err, "validation failed for update course request", "course_id", id)
		return nil, err
	}

	// Get existing course.
	existingCourse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if apiErr, ok := errors.IsAPIError(err); ok {
			return nil, apiErr
		}
		logger.LogError(s.logger, ctx, err, "failed to get course for update", "course_id", id)

		return nil, errors.NewInternalError("failed to retrieve course for update", err)
	}

	// Update fields if provided.
	updatedCourse := *existingCourse
	if req.Title != nil {
		updatedCourse.Title = *req.Title
	}
	if req.Description != nil {
		updatedCourse.Description = *req.Description
	}
	updatedCourse.UpdatedAt = time.Now()

	// Save updated course.
	if err := s.repo.Update(ctx, &updatedCourse); err != nil {
		if apiErr, ok := errors.IsAPIError(err); ok {
			return nil, apiErr
		}
		logger.LogError(s.logger, ctx, err, "failed to update course", "course_id", id)

		return nil, errors.NewInternalError("failed to update course", err)
	}

	s.logger.Info(ctx, "course updated successfully",
		"course_id", id,
		"title", updatedCourse.Title,
	)

	return &updatedCourse, nil
}

// DeleteCourse deletes a course by ID.
func (s *courseService) DeleteCourse(ctx context.Context, id string) error {
	if id == "" {
		return errors.NewBadRequestError("course ID is required")
	}

	// Check if course exists.
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if apiErr, ok := errors.IsAPIError(err); ok {
			return apiErr
		}
		logger.LogError(s.logger, ctx, err, "failed to get course for deletion", "course_id", id)

		return errors.NewInternalError("failed to retrieve course for deletion", err)
	}

	// Delete course.
	if err := s.repo.Delete(ctx, id); err != nil {
		if apiErr, ok := errors.IsAPIError(err); ok {
			return apiErr
		}
		logger.LogError(s.logger, ctx, err, "failed to delete course", "course_id", id)

		return errors.NewInternalError("failed to delete course", err)
	}

	s.logger.Info(ctx, "course deleted successfully", "course_id", id)

	return nil
}

// Helper function to generate slug from title.
func generateSlug(title string) string {
	// Convert to lowercase.
	slug := strings.ToLower(title)

	// Replace spaces with hyphens.
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Remove multiple consecutive hyphens.
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end.
	slug = strings.Trim(slug, "-")

	// Ensure we have a valid slug.
	if slug == "" {
		slug = fmt.Sprintf("course-%d", time.Now().Unix())
	}

	return slug
}
