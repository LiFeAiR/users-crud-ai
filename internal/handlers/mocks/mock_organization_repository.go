package mocks

import (
	"context"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/stretchr/testify/mock"
)

// Mock OrganizationRepository для тестирования
type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) CreateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	args := m.Called(ctx, org)

	if o, ok := args.Get(0).(*models.Organization); ok {
		return o, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockOrganizationRepository) GetOrganizationByID(ctx context.Context, id int) (*models.Organization, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) UpdateOrganization(ctx context.Context, org *models.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) DeleteOrganization(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrganizationRepository) GetOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) InitDB() error {
	panic("implement me")
}
