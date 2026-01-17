package handlers

import (
	"context"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	// Test 1: Успешное удаление пользователя
	t.Run("DeleteUserSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("DeleteUser", ctx, 1).Return(nil, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод DeleteUser
		_, err := baseHandler.DeleteUser(ctx, &grpc.Id{Id: 1})

		// Проверяем результат
		assert.NoError(t, err)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
