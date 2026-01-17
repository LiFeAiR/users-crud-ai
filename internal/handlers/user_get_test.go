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

// TestBaseHandler_GetUser тестирует метод GetUser базового обработчика
func TestBaseHandler_GetUser(t *testing.T) {
	ctx := context.Background()
	// Test 1: Успешное получение пользователя
	t.Run("GetUserSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Подготавливаем тестовую организацию
		testUser := &models.User{
			ID:   1,
			Name: "Test Org",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetUserByID", ctx, 1).Return(testUser, nil)
		mockRepo.On("GetUserPermissions", ctx, 1).Return([]*models.Permission{{
			ID:          0,
			Name:        "1",
			Code:        "2",
			Description: "3",
		}}, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод GetUser
		u, err := baseHandler.GetUser(ctx, &grpc.Id{Id: 1})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, int32(1), u.Id)
		assert.Equal(t, "Test Org", u.Name)
		assert.Len(t, u.Permissions, 1)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Некорректный аргумент (nil)
	t.Run("GetUserNilArgument", func(t *testing.T) {
		// Создаем базовый обработчик
		baseHandler := &BaseHandler{
			orgRepo: nil,
		}

		// Вызываем метод GetUser с nil аргументом
		org, err := baseHandler.GetUser(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, org)
	})

	// Test 3: Пользователь не найден
	t.Run("GetUserNotFound", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("GetUserByID", ctx, 1).Return((*models.User)(nil), errors.New("user not found"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод GetUser
		org, err := baseHandler.GetUser(ctx, &grpc.Id{Id: 1})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, org)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
