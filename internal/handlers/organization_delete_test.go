package handlers

import (
	"context"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOrganization(t *testing.T) {
	ctx := context.Background()
	// Test 1: Успешное удаление организации
	t.Run("DeleteOrganizationSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockOrganizationRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("DeleteOrganization", ctx, 1).Return(nil, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Вызываем метод DeleteOrganization
		_, err := baseHandler.DeleteOrganization(ctx, &grpc.Id{Id: 1})

		// Проверяем результат
		assert.NoError(t, err)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
