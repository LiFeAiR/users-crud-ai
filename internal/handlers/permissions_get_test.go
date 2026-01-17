package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

// TestBaseHandler_GetPermissions тестирует метод GetPermissions базового обработчика
func TestBaseHandler_GetPermissions(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное получение списка прав
	t.Run("GetPermissionsSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockPermissionRepository)

		// Подготавливаем тестовые права для возврата из репозитория
		expectedPermissions := []*models.Permission{
			{
				ID:          1,
				Name:        "Permission 1",
				Code:        "permission_1",
				Description: "Description 1",
			},
			{
				ID:          2,
				Name:        "Permission 2",
				Code:        "permission_2",
				Description: "Description 2",
			},
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetPermissions", ctx, 10, 0).Return(expectedPermissions, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод GetPermissions
		result, err := baseHandler.GetPermissions(ctx, &grpc.ListRequest{
			Limit:  10,
			Offset: 0,
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, int32(1), result.Data[0].Id)
		assert.Equal(t, "Permission 1", result.Data[0].Name)
		assert.Equal(t, "permission_1", result.Data[0].Code)
		assert.Equal(t, "Description 1", result.Data[0].Description)
		assert.Equal(t, int32(2), result.Data[1].Id)
		assert.Equal(t, "Permission 2", result.Data[1].Name)
		assert.Equal(t, "permission_2", result.Data[1].Code)
		assert.Equal(t, "Description 2", result.Data[1].Description)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при получении списка прав из репозитория
	t.Run("GetPermissionsRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockPermissionRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("GetPermissions", ctx, 10, 0).Return(([]*models.Permission)(nil), errors.New("get permissions failed"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод GetPermissions
		result, err := baseHandler.GetPermissions(ctx, &grpc.ListRequest{
			Limit:  10,
			Offset: 0,
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 3: Ошибка при передаче некорректных данных
	t.Run("GetPermissionsInvalidInput", func(t *testing.T) {
		// Создаем базовый обработчик без мока (тест на валидацию)
		baseHandler := &BaseHandler{}

		// Вызываем метод GetPermissions с nil параметром
		result, err := baseHandler.GetPermissions(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
