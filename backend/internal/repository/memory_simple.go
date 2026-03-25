// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package repository provides functionality for the GO-PRO Learning Platform.
package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
)

// MemoryCourseRepository implements CourseRepository using in-memory storage.
type MemoryCourseRepository struct {
	courses map[string]*domain.Course
	mu      sync.RWMutex
}

// NewMemoryCourseRepository creates a new in-memory course repository.
func NewMemoryCourseRepository() *MemoryCourseRepository {
	return &MemoryCourseRepository{
		courses: make(map[string]*domain.Course),
	}
}

// Create implements CourseRepository.Create.
func (r *MemoryCourseRepository) Create(ctx context.Context, course *domain.Course) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.courses[course.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("course with id %s already exists", course.ID))
	}

	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()
	r.courses[course.ID] = course

	return nil
}

// GetByID implements CourseRepository.GetByID.
func (r *MemoryCourseRepository) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	course, exists := r.courses[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", id))
	}

	// Return a copy to prevent modification.
	courseCopy := *course

	return &courseCopy, nil
}

// GetAll implements CourseRepository.GetAll.
func (r *MemoryCourseRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Course, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var courses []*domain.Course
	for _, course := range r.courses {
		courseCopy := *course
		courses = append(courses, &courseCopy)
	}

	// Sort by creation time (newest first)
	sort.Slice(courses, func(i, j int) bool {
		return courses[i].CreatedAt.After(courses[j].CreatedAt)
	})

	total := int64(len(courses))

	// Apply pagination if provided.
	if pagination != nil {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize

		if start >= len(courses) {
			return []*domain.Course{}, total, nil
		}

		if end > len(courses) {
			end = len(courses)
		}

		courses = courses[start:end]
	}

	return courses, total, nil
}

// Update implements CourseRepository.Update.
func (r *MemoryCourseRepository) Update(ctx context.Context, course *domain.Course) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.courses[course.ID]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", course.ID))
	}

	course.UpdatedAt = time.Now()
	r.courses[course.ID] = course

	return nil
}

// Delete implements CourseRepository.Delete.
func (r *MemoryCourseRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.courses[id]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("course with id %s not found", id))
	}

	delete(r.courses, id)

	return nil
}

// MemoryUserRepository implements UserRepository using in-memory storage.
type MemoryUserRepository struct {
	users           map[string]*domain.User // Keyed by user ID
	usersByFirebase map[string]*domain.User // Keyed by Firebase UID
	usersByEmail    map[string]*domain.User // Keyed by email
	mu              sync.RWMutex
}

// NewMemoryUserRepository creates a new in-memory user repository.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:           make(map[string]*domain.User),
		usersByFirebase: make(map[string]*domain.User),
		usersByEmail:    make(map[string]*domain.User),
	}
}

// Create implements UserRepository.Create.
func (r *MemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("user with id %s already exists", user.ID))
	}

	if _, exists := r.usersByFirebase[user.FirebaseUID]; exists {
		return errors.NewConflictError(fmt.Sprintf("user with Firebase UID %s already exists", user.FirebaseUID))
	}

	if _, exists := r.usersByEmail[user.Email]; exists {
		return errors.NewConflictError(fmt.Sprintf("user with email %s already exists", user.Email))
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	r.usersByFirebase[user.FirebaseUID] = user
	r.usersByEmail[user.Email] = user

	return nil
}

// GetByID implements UserRepository.GetByID.
func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	userCopy := *user
	return &userCopy, nil
}

// GetByFirebaseUID implements UserRepository.GetByFirebaseUID.
func (r *MemoryUserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByFirebase[firebaseUID]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with Firebase UID %s not found", firebaseUID))
	}

	userCopy := *user
	return &userCopy, nil
}

// GetByEmail implements UserRepository.GetByEmail.
func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByEmail[email]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
	}

	userCopy := *user
	return &userCopy, nil
}

// GetAll implements UserRepository.GetAll.
func (r *MemoryUserRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var users []*domain.User
	for _, user := range r.users {
		userCopy := *user
		users = append(users, &userCopy)
	}

	// Sort by creation time (newest first)
	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.After(users[j].CreatedAt)
	})

	total := int64(len(users))

	// Apply pagination
	if pagination != nil && pagination.PageSize > 0 {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize

		if start >= len(users) {
			return []*domain.User{}, total, nil
		}

		if end > len(users) {
			end = len(users)
		}

		users = users[start:end]
	}

	return users, total, nil
}

// Update implements UserRepository.Update.
func (r *MemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.users[user.ID]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", user.ID))
	}

	// Update indexes if email or Firebase UID changed
	if existing.Email != user.Email {
		delete(r.usersByEmail, existing.Email)
		r.usersByEmail[user.Email] = user
	}

	if existing.FirebaseUID != user.FirebaseUID {
		delete(r.usersByFirebase, existing.FirebaseUID)
		r.usersByFirebase[user.FirebaseUID] = user
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

// UpdateLastLogin implements UserRepository.UpdateLastLogin.
func (r *MemoryUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", userID))
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now

	return nil
}

// UpdateLastActivity implements UserRepository.UpdateLastActivity.
func (r *MemoryUserRepository) UpdateLastActivity(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", userID))
	}

	user.UpdatedAt = time.Now()

	return nil
}

// Delete implements UserRepository.Delete.
func (r *MemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	delete(r.users, id)
	delete(r.usersByFirebase, user.FirebaseUID)
	delete(r.usersByEmail, user.Email)

	return nil
}

// MemoryProgressRepository implements ProgressRepository using in-memory storage.
type MemoryProgressRepository struct {
	progress     map[string]*domain.Progress
	userProgress map[string][]string
	userLesson   map[string]string
	mu           sync.RWMutex
}

// NewMemoryProgressRepository creates a new in-memory progress repository.
func NewMemoryProgressRepository() *MemoryProgressRepository {
	return &MemoryProgressRepository{
		progress:     make(map[string]*domain.Progress),
		userProgress: make(map[string][]string),
		userLesson:   make(map[string]string),
	}
}

func (r *MemoryProgressRepository) Create(ctx context.Context, progress *domain.Progress) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.progress[progress.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("progress with id %s already exists", progress.ID))
	}
	key := fmt.Sprintf("%s:%s", progress.UserID, progress.LessonID)
	if _, exists := r.userLesson[key]; exists {
		return errors.NewConflictError(fmt.Sprintf("progress already exists for user %s and lesson %s", progress.UserID, progress.LessonID))
	}
	progress.CreatedAt = time.Now()
	progress.UpdatedAt = time.Now()
	r.progress[progress.ID] = progress
	r.userProgress[progress.UserID] = append(r.userProgress[progress.UserID], progress.ID)
	r.userLesson[key] = progress.ID
	return nil
}

func (r *MemoryProgressRepository) GetByID(ctx context.Context, id string) (*domain.Progress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	progress, exists := r.progress[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("progress with id %s not found", id))
	}
	progressCopy := *progress
	return &progressCopy, nil
}

func (r *MemoryProgressRepository) GetByUserID(ctx context.Context, userID string, pagination *domain.PaginationRequest) ([]*domain.Progress, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	progressIDs := r.userProgress[userID]
	progressList := make([]*domain.Progress, 0, len(progressIDs))
	for _, id := range progressIDs {
		if progress, exists := r.progress[id]; exists {
			progressCopy := *progress
			progressList = append(progressList, &progressCopy)
		}
	}
	total := int64(len(progressList))
	if pagination != nil && pagination.PageSize > 0 {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize
		if start >= len(progressList) {
			return []*domain.Progress{}, total, nil
		}
		if end > len(progressList) {
			end = len(progressList)
		}
		progressList = progressList[start:end]
	}
	return progressList, total, nil
}

func (r *MemoryProgressRepository) GetByUserAndLesson(ctx context.Context, userID, lessonID string) (*domain.Progress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", userID, lessonID)
	progressID, exists := r.userLesson[key]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("progress not found for user %s and lesson %s", userID, lessonID))
	}
	progress, exists := r.progress[progressID]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("progress with id %s not found", progressID))
	}
	progressCopy := *progress
	return &progressCopy, nil
}

func (r *MemoryProgressRepository) Update(ctx context.Context, progress *domain.Progress) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.progress[progress.ID]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("progress with id %s not found", progress.ID))
	}
	progress.UpdatedAt = time.Now()
	r.progress[progress.ID] = progress
	return nil
}

func (r *MemoryProgressRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	progress, exists := r.progress[id]
	if !exists {
		return errors.NewNotFoundError(fmt.Sprintf("progress with id %s not found", id))
	}
	key := fmt.Sprintf("%s:%s", progress.UserID, progress.LessonID)
	delete(r.userLesson, key)
	progressIDs := r.userProgress[progress.UserID]
	for i, pid := range progressIDs {
		if pid == id {
			r.userProgress[progress.UserID] = append(progressIDs[:i], progressIDs[i+1:]...)
			break
		}
	}
	delete(r.progress, id)
	return nil
}

// MemoryExerciseRepository implements ExerciseRepository using in-memory storage.
type MemoryExerciseRepository struct {
	exercises map[string]*domain.Exercise
	mu        sync.RWMutex
}

// NewMemoryExerciseRepository creates a new in-memory exercise repository.
func NewMemoryExerciseRepository() *MemoryExerciseRepository {
	repo := &MemoryExerciseRepository{
		exercises: make(map[string]*domain.Exercise),
	}
	repo.seedSampleExercises()
	return repo
}

func (r *MemoryExerciseRepository) seedSampleExercises() {
	now := time.Now()
	sampleExercises := []*domain.Exercise{
		{
			ID:          "1",
			LessonID:    "1",
			Title:       "FizzBuzz Challenge",
			Description: "Write a function that prints numbers from 1 to n. For multiples of 3 print 'Fizz', for multiples of 5 print 'Buzz', and for multiples of both print 'FizzBuzz'.",
			TestCases:   5,
			Difficulty:  domain.DifficultyBeginner,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "2",
			LessonID:    "1",
			Title:       "Reverse a String",
			Description: "Write a function that takes a string as input and returns the string reversed.",
			TestCases:   3,
			Difficulty:  domain.DifficultyBeginner,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "3",
			LessonID:    "2",
			Title:       "Binary Search Implementation",
			Description: "Implement a binary search algorithm that finds the index of a target value in a sorted array. Return -1 if the target is not found.",
			TestCases:   7,
			Difficulty:  domain.DifficultyIntermediate,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "4",
			LessonID:    "3",
			Title:       "Concurrent Web Scraper",
			Description: "Build a concurrent web scraper using goroutines that fetches URLs from a list and returns their content. Use proper synchronization and error handling.",
			TestCases:   10,
			Difficulty:  domain.DifficultyAdvanced,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	for _, exercise := range sampleExercises {
		r.exercises[exercise.ID] = exercise
	}
}

func (r *MemoryExerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.exercises[exercise.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("exercise with id %s already exists", exercise.ID))
	}
	exercise.CreatedAt = time.Now()
	exercise.UpdatedAt = time.Now()
	r.exercises[exercise.ID] = exercise
	return nil
}

func (r *MemoryExerciseRepository) GetByID(ctx context.Context, id string) (*domain.Exercise, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exercise, exists := r.exercises[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("exercise with id %s not found", id))
	}
	exerciseCopy := *exercise
	return &exerciseCopy, nil
}

func (r *MemoryExerciseRepository) GetByLessonID(ctx context.Context, lessonID string, pagination *domain.PaginationRequest) ([]*domain.Exercise, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var exercises []*domain.Exercise
	for _, exercise := range r.exercises {
		if exercise.LessonID == lessonID {
			exerciseCopy := *exercise
			exercises = append(exercises, &exerciseCopy)
		}
	}
	total := int64(len(exercises))
	if pagination != nil {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize
		if start >= len(exercises) {
			return []*domain.Exercise{}, total, nil
		}
		if end > len(exercises) {
			end = len(exercises)
		}
		exercises = exercises[start:end]
	}
	return exercises, total, nil
}

func (r *MemoryExerciseRepository) GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.Exercise, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var exercises []*domain.Exercise
	for _, exercise := range r.exercises {
		exerciseCopy := *exercise
		exercises = append(exercises, &exerciseCopy)
	}
	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].CreatedAt.After(exercises[j].CreatedAt)
	})
	total := int64(len(exercises))
	if pagination != nil {
		start := (pagination.Page - 1) * pagination.PageSize
		end := start + pagination.PageSize
		if start >= len(exercises) {
			return []*domain.Exercise{}, total, nil
		}
		if end > len(exercises) {
			end = len(exercises)
		}
		exercises = exercises[start:end]
	}
	return exercises, total, nil
}

func (r *MemoryExerciseRepository) Update(ctx context.Context, exercise *domain.Exercise) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.exercises[exercise.ID]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("exercise with id %s not found", exercise.ID))
	}
	exercise.UpdatedAt = time.Now()
	r.exercises[exercise.ID] = exercise
	return nil
}

func (r *MemoryExerciseRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.exercises[id]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("exercise with id %s not found", id))
	}
	delete(r.exercises, id)
	return nil
}

// NewRepositoriesSimple creates repository instances with simple approach.
func NewRepositoriesSimple() *Repositories {
	return &Repositories{
		Course:    NewMemoryCourseRepository(),
		User:      NewMemoryUserRepository(),
		Progress:  NewMemoryProgressRepository(),
		Exercise:  NewMemoryExerciseRepository(),
		Interview: NewMemoryInterviewRepository(),
	}
}
