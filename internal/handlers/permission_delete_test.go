package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

// TestBaseHandler_DeletePermission тестирует метод DeletePermission базового обработчика
func TestBaseHandler_DeletePermission(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное удаление права
	t.Run("DeletePermissionSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockPermissionRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("DeletePermission", ctx, 1).Return(nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод DeletePermission
		result, err := baseHandler.DeletePermission(ctx, &grpc.Id{
			Id: 1,
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при удалении в репозитории
	t.Run("DeletePermissionRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockPermissionRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("DeletePermission", ctx, 1).Return(errors.New("delete failed"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			permRepo: mockRepo,
		}

		// Вызываем метод DeletePermission
		result, err := baseHandler.DeletePermission(ctx, &grpc.Id{
			Id: 1,
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 3: Ошибка при передаче некорректных данных
	t.Run("DeletePermissionInvalidInput", func(t *testing.T) {
		// Создаем базовый обработчик без мока (тест на валидацию)
		baseHandler := &BaseHandler{}

		// Вызываем метод DeletePermission с nil параметром
		result, err := baseHandler.DeletePermission(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
