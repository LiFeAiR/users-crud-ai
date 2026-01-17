package mocks

import (
	"context"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockRoleRepository имитация репозитория ролей для тестирования
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) CreateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	args := m.Called(ctx, role)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleByID(ctx context.Context, id int) (*models.Role, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) DeleteRole(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) GetRoles(ctx context.Context, limit, offset int) ([]*models.Role, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleWithPermissions(ctx context.Context, id int) (*models.Role, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) InitDB() error {
	args := m.Called()
	return args.Error(0)
}
