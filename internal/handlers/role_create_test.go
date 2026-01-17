package handlers

import (
	"context"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRole(t *testing.T) {
	// Создаем мок репозитория
	mockRoleRepo := new(mocks.MockRoleRepository)

	// Создаем базовый обработчик с моком
	handler := &BaseHandler{
		roleRepo: mockRoleRepo,
	}

	// Создаем запрос
	request := &api_pb.RoleCreateRequest{
		Name:        "Test Role",
		Code:        "test_role",
		Description: "Test role description",
	}

	// Определяем ожидаемое поведение мока
	expectedRole := &models.Role{
		ID:          1,
		Name:        "Test Role",
		Code:        "test_role",
		Description: "Test role description",
	}

	mockRoleRepo.On("CreateRole", mock.Anything, mock.MatchedBy(func(role *models.Role) bool {
		return role.Name == "Test Role" && role.Code == "test_role"
	})).Return(expectedRole, nil)

	// Выполняем метод
	response, err := handler.CreateRole(context.Background(), request)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int32(1), response.Id)
	assert.Equal(t, "Test Role", response.Name)
	assert.Equal(t, "test_role", response.Code)
	assert.Equal(t, "Test role description", response.Description)

	// Проверяем, что мок был вызван
	mockRoleRepo.AssertExpectations(t)
}
