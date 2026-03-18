// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package repository provides in-memory interview repository implementation.
package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"go-pro-backend/internal/errors"
)

// InterviewSession represents an interview session (same as handler).
type InterviewSession struct {
	ID           string      `json:"id"`
	UserID       string      `json:"user_id"`
	Type         string      `json:"type"`
	Difficulty   string      `json:"difficulty"`
	Questions    []Question  `json:"questions"`
	CurrentIndex int         `json:"current_index"`
	Answers      []Answer    `json:"answers"`
	Status       string      `json:"status"`
	Score        *int        `json:"score,omitempty"`
	CreatedAt    int64       `json:"created_at"`
	CompletedAt  *int64      `json:"completed_at,omitempty"`
}

// Question represents an interview question (same as handler).
type Question struct {
	ID             string   `json:"id"`
	Content        string   `json:"content"`
	Type           string   `json:"type"`
	Difficulty     string   `json:"difficulty"`
	ExpectedPoints []string `json:"expected_points,omitempty"`
	TimeLimit      int      `json:"time_limit"`
}

// Answer represents a user's answer (same as handler).
type Answer struct {
	QuestionID string  `json:"question_id"`
	Content    string  `json:"content"`
	Score      *int    `json:"score,omitempty"`
	Feedback   string  `json:"feedback,omitempty"`
	CreatedAt  int64   `json:"created_at"`
}

// MemoryInterviewRepository implements InterviewRepository using in-memory storage.
type MemoryInterviewRepository struct {
	sessions map[string]*InterviewSession
	mu       sync.RWMutex
}

// NewMemoryInterviewRepository creates a new in-memory interview repository.
func NewMemoryInterviewRepository() *MemoryInterviewRepository {
	return &MemoryInterviewRepository{
		sessions: make(map[string]*InterviewSession),
	}
}

// Create implements InterviewRepository.Create.
func (r *MemoryInterviewRepository) Create(ctx context.Context, session interface{}) error {
	s, ok := session.(*InterviewSession)
	if !ok {
		return errors.NewBadRequestError("invalid session type")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[s.ID]; exists {
		return errors.NewConflictError(fmt.Sprintf("session with id %s already exists", s.ID))
	}

	sessionCopy := *s
	r.sessions[s.ID] = &sessionCopy

	return nil
}

// GetByID implements InterviewRepository.GetByID.
func (r *MemoryInterviewRepository) GetByID(ctx context.Context, id string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("session with id %s not found", id))
	}

	sessionCopy := *session
	return &sessionCopy, nil
}

// GetByUserID implements InterviewRepository.GetByUserID.
func (r *MemoryInterviewRepository) GetByUserID(ctx context.Context, userID string) ([]interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sessions []interface{}
	for _, session := range r.sessions {
		if session.UserID == userID {
			sessionCopy := *session
			sessions = append(sessions, &sessionCopy)
		}
	}

	sort.Slice(sessions, func(i, j int) bool {
		si, _ := sessions[i].(*InterviewSession)
		sj, _ := sessions[j].(*InterviewSession)
		return si.CreatedAt > sj.CreatedAt
	})

	return sessions, nil
}

// Update implements InterviewRepository.Update.
func (r *MemoryInterviewRepository) Update(ctx context.Context, session interface{}) error {
	s, ok := session.(*InterviewSession)
	if !ok {
		return errors.NewBadRequestError("invalid session type")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[s.ID]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("session with id %s not found", s.ID))
	}

	sessionCopy := *s
	r.sessions[s.ID] = &sessionCopy

	return nil
}

// Delete implements InterviewRepository.Delete.
func (r *MemoryInterviewRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[id]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("session with id %s not found", id))
	}

	delete(r.sessions, id)

	return nil
}
