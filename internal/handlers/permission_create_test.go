package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
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

// TestBaseHandler_CreatePermission тестирует метод CreatePermission базового обработчика
func TestBaseHandler_CreatePermission(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное создание права
	t.Run("CreatePermissionSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(MockPermissionRepository)

		// Подготавливаем тестовое право для возврата из репозитория
		expectedPermission := &models.Permission{
			ID:          1,
			Name:        "Test Permission",
			Code:        "test_permission",
			Description: "Test description",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("CreatePermission", ctx, mock.Anything).Return(expectedPermission, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод CreatePermission
		result, err := baseHandler.CreatePermission(ctx, &grpc.PermissionCreateRequest{
			Name:        "Test Permission",
			Code:        "test_permission",
			Description: "Test description",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Test Permission", result.Name)
		assert.Equal(t, "test_permission", result.Code)
		assert.Equal(t, "Test description", result.Description)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при отсутствии обязательных полей
	t.Run("CreatePermissionMissingFields", func(t *testing.T) {
		// Создаем базовый обработчик без мока (тест на валидацию)
		baseHandler := &BaseHandler{}

		// Вызываем метод CreatePermission с неполными данными
		result, err := baseHandler.CreatePermission(ctx, &grpc.PermissionCreateRequest{
			Name: "",
			Code: "",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test 3: Ошибка при неудачном создании в репозитории
	t.Run("CreatePermissionRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(MockPermissionRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("CreatePermission", ctx, mock.Anything).
			Return((*models.Permission)(nil), errors.New("create failed"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод CreatePermission
		result, err := baseHandler.CreatePermission(ctx, &grpc.PermissionCreateRequest{
			Name:        "Test Permission",
			Code:        "test_permission",
			Description: "Test description",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}