package mocks

import (
	"context"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockPermissionRepository мок для репозитория прав
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) CreatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	args := m.Called(ctx, permission)
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetPermissionByID(ctx context.Context, id int) (*models.Permission, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) DeletePermission(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) InitDB() error {
	args := m.Called()
	return args.Error(0)
}
