package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestBaseHandler_CreateUser тестирует метод CreateUser базового обработчика
func TestBaseHandler_CreateUser(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное создание пользователя
	t.Run("CreateUserSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Подготавливаем тестового пользователя для возврата из репозитория
		expectedUser := &models.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("CreateUser", ctx, mock.Anything).Return(expectedUser, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод CreateUser
		result, err := baseHandler.CreateUser(ctx, &grpc.UserCreateRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Test User", result.Name)
		assert.Equal(t, "test@example.com", result.Email)
		//assert.Equal(t, nil, result.Organization)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при неудачном создании в репозитории
	t.Run("CreateUserRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("CreateUser", ctx, mock.Anything).
			Return((*models.User)(nil), errors.New("create failed"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод CreateUser
		result, err := baseHandler.CreateUser(ctx, &grpc.UserCreateRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
