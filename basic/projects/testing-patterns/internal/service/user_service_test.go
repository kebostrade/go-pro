package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil {
		user.ID = "generated-id"
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*service.User")).Return(nil)

		user, err := svc.Create(context.Background(), "Alice", "alice@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "Alice", user.Name)
		assert.Equal(t, "alice@example.com", user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		user, err := svc.Create(context.Background(), "", "alice@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("empty email", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		user, err := svc.Create(context.Background(), "Bob", "")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrValidation)
	})
}

func TestUserService_GetByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		expectedUser := &User{ID: "123", Name: "Test", Email: "test@example.com"}
		mockRepo.On("GetByID", mock.Anything, "123").Return(expectedUser, nil)

		user, err := svc.GetByID(context.Background(), "123")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, ErrNotFound)

		user, err := svc.GetByID(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_List(t *testing.T) {
	t.Run("list users", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		users := []*User{
			{ID: "1", Name: "User1", Email: "user1@example.com"},
			{ID: "2", Name: "User2", Email: "user2@example.com"},
		}
		mockRepo.On("List", mock.Anything).Return(users, nil)

		result, err := svc.List(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		mockRepo.On("Delete", mock.Anything, "123").Return(nil)

		err := svc.Delete(context.Background(), "123")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewUserService(mockRepo)

		mockRepo.On("Delete", mock.Anything, "non-existent").Return(ErrNotFound)

		err := svc.Delete(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		mockRepo.AssertExpectations(t)
	})
}
