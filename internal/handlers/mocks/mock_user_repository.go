package mocks

import (
	"context"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepository для тестирования
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) SetUserTariff(ctx context.Context, userID int, tariffID *int32) error {
	args := m.Called(ctx, userID, tariffID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserTariff(ctx context.Context, userID int) (*models.Tariff, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.Tariff), args.Error(1)
}

func (m *MockUserRepository) GetUserPermissions(ctx context.Context, userID int) ([]*models.Permission, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Permission), args.Error(1)
}

func (m *MockUserRepository) AddUserPermissions(ctx context.Context, userID int, permissionIDs []int) error {
	args := m.Called(ctx, userID, permissionIDs)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserPermissions(ctx context.Context, userID int, permissionIDs []int) error {
	args := m.Called(ctx, userID, permissionIDs)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserRoles(ctx context.Context, userID int) ([]*models.Role, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Role), args.Error(1)
}

func (m *MockUserRepository) AddUserRoles(ctx context.Context, userID int, roleIDs []int) error {
	args := m.Called(ctx, userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserRoles(ctx context.Context, userID int, roleIDs []int) error {
	args := m.Called(ctx, userID, roleIDs)
	return args.Error(0)
}

func (m *MockUserRepository) CheckPassword(ctx context.Context, userID int, password string) (bool, error) {
	args := m.Called(ctx, userID, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) InitDB() error {
	panic("implement me")
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)

	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
