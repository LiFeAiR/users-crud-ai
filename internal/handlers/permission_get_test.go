package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

// TestBaseHandler_GetPermission тестирует метод GetPermission базового обработчика
func TestBaseHandler_GetPermission(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное получение права
	t.Run("GetPermissionSuccess", func(t *testing.T) {
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
		mockRepo.On("GetPermissionByID", ctx, 1).Return(expectedPermission, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод GetPermission
		result, err := baseHandler.GetPermission(ctx, &grpc.Id{
			Id: 1,
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

	// Test 2: Ошибка при получении несуществующего права
	t.Run("GetPermissionNotFound", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(MockPermissionRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("GetPermissionByID", ctx, 999).Return((*models.Permission)(nil), errors.New("permission not found"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод GetPermission
		result, err := baseHandler.GetPermission(ctx, &grpc.Id{
			Id: 999,
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 3: Ошибка при передаче некорректных данных
	t.Run("GetPermissionInvalidInput", func(t *testing.T) {
		// Создаем базовый обработчик без мока (тест на валидацию)
		baseHandler := &BaseHandler{}

		// Вызываем метод GetPermission с nil параметром
		result, err := baseHandler.GetPermission(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
