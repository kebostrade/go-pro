package testutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/domain"
)

// MockCacheManager is a mock implementation of cache.CacheManager
type MockCacheManager struct {
	mu    sync.RWMutex
	data  map[string]interface{}
	calls map[string]int
}

// NewMockCacheManager creates a new mock cache manager
func NewMockCacheManager() *MockCacheManager {
	return &MockCacheManager{
		data:  make(map[string]interface{}),
		calls: make(map[string]int),
	}
}

func (m *MockCacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.calls["Get"]++

	if val, ok := m.data[key]; ok {
		// In a real implementation, we would unmarshal into dest
		// For mock, we just check if the key exists
		_ = dest
		_ = val
		return nil
	}
	return cache.ErrCacheMiss
}

func (m *MockCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls["Set"]++

	m.data[key] = value
	return nil
}

func (m *MockCacheManager) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls["Delete"]++

	delete(m.data, key)
	return nil
}

func (m *MockCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.calls["Exists"]++

	_, ok := m.data[key]
	return ok, nil
}

func (m *MockCacheManager) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[method]
}

func (m *MockCacheManager) ResetMock() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]interface{})
	m.calls = make(map[string]int)
}

// Implement remaining CacheManager interface methods
// SessionStore interface implementation
func (m *MockCacheManager) CreateSession(ctx context.Context, sessionID string, data map[string]interface{}, expiration time.Duration) error {
	return m.Set(ctx, "session:"+sessionID, data, expiration)
}

func (m *MockCacheManager) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := "session:" + sessionID
	if val, ok := m.data[key]; ok {
		if data, ok := val.(map[string]interface{}); ok {
			return data, nil
		}
	}
	return nil, cache.ErrCacheMiss
}

func (m *MockCacheManager) UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	return m.Set(ctx, "session:"+sessionID, data, 0)
}

func (m *MockCacheManager) DeleteSession(ctx context.Context, sessionID string) error {
	return m.Delete(ctx, "session:"+sessionID)
}

func (m *MockCacheManager) RefreshSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	return nil
}

func (m *MockCacheManager) ListUserSessions(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (m *MockCacheManager) DeleteUserSessions(ctx context.Context, userID string) error {
	return nil
}

func (m *MockCacheManager) CleanupExpiredSessions(ctx context.Context) error {
	return nil
}

// DistributedLock interface implementation
func (m *MockCacheManager) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return true, nil
}

func (m *MockCacheManager) Unlock(ctx context.Context, key string) error {
	return nil
}

func (m *MockCacheManager) Extend(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (m *MockCacheManager) IsLocked(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m *MockCacheManager) GetLockTTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (m *MockCacheManager) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	return true, nil
}

func (m *MockCacheManager) AllowN(ctx context.Context, key string, n int64, limit int64, window time.Duration) (bool, error) {
	return true, nil
}

func (m *MockCacheManager) Remaining(ctx context.Context, key string, limit int64, window time.Duration) (int64, error) {
	return limit, nil
}

func (m *MockCacheManager) Reset(ctx context.Context, key string) error {
	return m.Delete(ctx, key)
}

// PubSub interface implementation
func (m *MockCacheManager) Publish(ctx context.Context, channel string, message interface{}) error {
	return nil
}

func (m *MockCacheManager) Subscribe(ctx context.Context, channels ...string) (<-chan cache.Message, error) {
	ch := make(chan cache.Message)
	close(ch)
	return ch, nil
}

func (m *MockCacheManager) Unsubscribe(ctx context.Context, channels ...string) error {
	return nil
}

func (m *MockCacheManager) PSubscribe(ctx context.Context, patterns ...string) (<-chan cache.Message, error) {
	ch := make(chan cache.Message)
	close(ch)
	return ch, nil
}

func (m *MockCacheManager) PUnsubscribe(ctx context.Context, patterns ...string) error {
	return nil
}

// Cache interface additional methods
func (m *MockCacheManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (m *MockCacheManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (m *MockCacheManager) Increment(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockCacheManager) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (m *MockCacheManager) Decrement(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockCacheManager) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (m *MockCacheManager) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return true, nil
}

func (m *MockCacheManager) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return "", nil
}

func (m *MockCacheManager) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return []interface{}{}, nil
}

func (m *MockCacheManager) MSet(ctx context.Context, pairs ...interface{}) error {
	return nil
}

func (m *MockCacheManager) Keys(ctx context.Context, pattern string) ([]string, error) {
	return []string{}, nil
}

func (m *MockCacheManager) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return []string{}, 0, nil
}

func (m *MockCacheManager) FlushDB(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]interface{})
	return nil
}

func (m *MockCacheManager) HealthCheck(ctx context.Context) error {
	return nil
}

func (m *MockCacheManager) Close() error {
	return nil
}

// MockCourseRepository is a mock implementation of repository.CourseRepository
type MockCourseRepository struct {
	mu      sync.RWMutex
	courses map[string]*domain.Course
	calls   map[string]int
}

// NewMockCourseRepository creates a new mock course repository
func NewMockCourseRepository() *MockCourseRepository {
	return &MockCourseRepository{
		courses: make(map[string]*domain.Course),
		calls:   make(map[string]int),
	}
}

func (m *MockCourseRepository) Create(ctx context.Context, course *domain.Course) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls["Create"]++

	m.courses[course.ID] = course
	return nil
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.calls["GetByID"]++

	if course, ok := m.courses[id]; ok {
		return course, nil
	}
	return nil, fmt.Errorf("course not found")
}

func (m *MockCourseRepository) GetAll(ctx context.Context, req *domain.PaginationRequest) ([]*domain.Course, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.calls["GetAll"]++

	courses := make([]*domain.Course, 0, len(m.courses))
	for _, course := range m.courses {
		courses = append(courses, course)
	}

	return courses, int64(len(courses)), nil
}

func (m *MockCourseRepository) Update(ctx context.Context, course *domain.Course) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls["Update"]++

	if _, ok := m.courses[course.ID]; !ok {
		return fmt.Errorf("course not found")
	}

	m.courses[course.ID] = course
	return nil
}

func (m *MockCourseRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls["Delete"]++

	if _, ok := m.courses[id]; !ok {
		return fmt.Errorf("course not found")
	}

	delete(m.courses, id)
	return nil
}

func (m *MockCourseRepository) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[method]
}

func (m *MockCourseRepository) ResetCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = make(map[string]int)
}

func (m *MockCourseRepository) AddCourse(course *domain.Course) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.courses[course.ID] = course
}

// MockMessagingService is a mock implementation of messaging service
type MockMessagingService struct {
	mu        sync.RWMutex
	published []MockMessage
	calls     map[string]int
}

type MockMessage struct {
	Topic   string
	Key     string
	Payload interface{}
}

func NewMockMessagingService() *MockMessagingService {
	return &MockMessagingService{
		published: make([]MockMessage, 0),
		calls:     make(map[string]int),
	}
}

func (m *MockMessagingService) PublishUserEvent(ctx context.Context, eventType string, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls["PublishUserEvent"]++

	m.published = append(m.published, MockMessage{
		Topic:   "user-events",
		Key:     user.ID,
		Payload: user,
	})
	return nil
}

func (m *MockMessagingService) PublishCourseEvent(ctx context.Context, eventType string, course *domain.Course) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls["PublishCourseEvent"]++

	m.published = append(m.published, MockMessage{
		Topic:   "course-events",
		Key:     course.ID,
		Payload: course,
	})
	return nil
}

func (m *MockMessagingService) GetPublishedMessages() []MockMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]MockMessage{}, m.published...)
}

func (m *MockMessagingService) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[method]
}

func (m *MockMessagingService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.published = make([]MockMessage, 0)
	m.calls = make(map[string]int)
}
