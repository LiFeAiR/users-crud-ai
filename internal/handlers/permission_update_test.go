package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestBaseHandler_UpdatePermission тестирует метод UpdatePermission базового обработчика
func TestBaseHandler_UpdatePermission(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное обновление права
	t.Run("UpdatePermissionSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(MockPermissionRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("UpdatePermission", ctx, mock.Anything).Return(nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод UpdatePermission
		result, err := baseHandler.UpdatePermission(ctx, &grpc.PermissionUpdateRequest{
			Id:          1,
			Name:        "Updated Permission",
			Code:        "updated_permission",
			Description: "Updated description",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Updated Permission", result.Name)
		assert.Equal(t, "updated_permission", result.Code)
		assert.Equal(t, "Updated description", result.Description)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при обновлении в репозитории
	t.Run("UpdatePermissionRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(MockPermissionRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("UpdatePermission", ctx, mock.Anything).Return(errors.New("update failed"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод UpdatePermission
		result, err := baseHandler.UpdatePermission(ctx, &grpc.PermissionUpdateRequest{
			Id:          1,
			Name:        "Updated Permission",
			Code:        "updated_permission",
			Description: "Updated description",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 3: Ошибка при передаче некорректных данных
	t.Run("UpdatePermissionInvalidInput", func(t *testing.T) {
		// Создаем базовый обработчик без мока (тест на валидацию)
		baseHandler := &BaseHandler{}

		// Вызываем метод UpdatePermission с некорректными данными
		result, err := baseHandler.UpdatePermission(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Вызываем метод UpdatePermission с нулевым ID
		result, err = baseHandler.UpdatePermission(ctx, &grpc.PermissionUpdateRequest{
			Id: 0,
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
