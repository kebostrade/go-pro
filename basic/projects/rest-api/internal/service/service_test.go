package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/internal/domain"
	"github.com/DimaJoyti/go-pro/basic/projects/rest-api/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// mockUserRepository is a mock implementation for testing
type mockUserRepository struct {
	users  map[string]*domain.User
	nextID int
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[string]*domain.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	id := fmt.Sprintf("test-id-%d", m.nextID)
	m.nextID++
	user.ID = id
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, errors.ErrNotFound
}

func (m *mockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *mockUserRepository) Update(ctx context.Context, id, name, email string) (*domain.User, error) {
	if user, ok := m.users[id]; ok {
		user.Name = name
		user.Email = email
		return user, nil
	}
	return nil, errors.ErrNotFound
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.users[id]; ok {
		delete(m.users, id)
		return nil
	}
	return errors.ErrNotFound
}

func TestUserService_Create(t *testing.T) {
	t.Run("valid user creation", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		user, err := svc.Create(context.Background(), "Alice", "alice@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "Alice", user.Name)
		assert.Equal(t, "alice@example.com", user.Email)
		assert.NotEmpty(t, user.ID)
	})

	t.Run("invalid email format", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		user, err := svc.Create(context.Background(), "Bob", "invalid-email")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})
}

func TestUserService_GetByID(t *testing.T) {
	t.Run("existing user", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		// Create a user first
		created, _ := svc.Create(context.Background(), "TestUser", "test@example.com")

		// Get the user
		user, err := svc.GetByID(context.Background(), created.ID)

		assert.NoError(t, err)
		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, "TestUser", user.Name)
	})

	t.Run("non-existing user", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		user, err := svc.GetByID(context.Background(), "non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})
}

func TestUserService_List(t *testing.T) {
	t.Run("list users", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		// Create some users
		_, _ = svc.Create(context.Background(), "User1", "user1@example.com")
		_, _ = svc.Create(context.Background(), "User2", "user2@example.com")

		users, err := svc.List(context.Background())

		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})
}

func TestUserService_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		// Create a user first
		created, _ := svc.Create(context.Background(), "OriginalName", "original@example.com")

		// Update the user
		updated, err := svc.Update(context.Background(), created.ID, "UpdatedName", "updated@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "UpdatedName", updated.Name)
		assert.Equal(t, "updated@example.com", updated.Email)
	})

	t.Run("update non-existing user", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		user, err := svc.Update(context.Background(), "non-existent-id", "Name", "email@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("update with invalid email", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		// Create a user first
		created, _ := svc.Create(context.Background(), "TestUser", "test@example.com")

		// Try to update with invalid email
		user, err := svc.Update(context.Background(), created.ID, "NewName", "invalid-email")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})
}

func TestUserService_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		// Create a user first
		created, _ := svc.Create(context.Background(), "TestUser", "test@example.com")

		// Delete the user
		err := svc.Delete(context.Background(), created.ID)

		assert.NoError(t, err)

		// Verify user is deleted
		_, err = svc.GetByID(context.Background(), created.ID)
		assert.Error(t, err)
	})

	t.Run("delete non-existing user", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		svc := NewUserService(mockRepo)

		err := svc.Delete(context.Background(), "non-existent-id")

		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})
}
