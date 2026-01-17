package handlers

import (
	"context"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

// TestBaseHandler_GetUsers тестирует метод GetUsers базового обработчика
func TestBaseHandler_GetUsers(t *testing.T) {
	ctx := context.Background()
	// Test 1: Успешное получение пользователя
	t.Run("GetUsersSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Подготавливаем тестовую пользователя
		testOrg := &models.User{
			ID:   1,
			Name: "Test Org",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetUsers", ctx, 10, 0).Return([]*models.User{testOrg}, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод GetUsers
		orgs, err := baseHandler.GetUsers(ctx, &grpc.ListRequest{})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, orgs)
		assert.Equal(t, int32(1), orgs.Data[0].Id)
		assert.Equal(t, "Test Org", orgs.Data[0].Name)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
