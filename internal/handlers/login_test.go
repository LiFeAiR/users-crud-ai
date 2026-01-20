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

// TestBaseHandler_Login тестирует метод Login базового обработчика
func TestBaseHandler_Login(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешная авторизация
	t.Run("LoginSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Подготавливаем тестового пользователя
		expectedUser := &models.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetUserByEmail", ctx, "test@example.com").Return(expectedUser, nil)
		mockRepo.On("CheckPassword", ctx, 1, "password123").Return(true, nil)
		mockRepo.On("GetUserPermissions", ctx, 1).Return([]*models.Permission(nil), nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод Login
		result, err := baseHandler.Login(ctx, &grpc.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Token)
		assert.NotNil(t, result.User)
		assert.Equal(t, int32(1), result.User.Id)
		assert.Equal(t, "Test User", result.User.Name)
		assert.Equal(t, "test@example.com", result.User.Email)

		// Проверяем, что моки были вызваны правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при неверных учетных данных
	t.Run("LoginInvalidCredentials", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Определяем ожидаемое поведение мока - пользователь не найден
		mockRepo.On("GetUserByEmail", ctx, "test@example.com").Return((*models.User)(nil), nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод Login
		result, err := baseHandler.Login(ctx, &grpc.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = Invalid credentials", err.Error())

		// Проверяем, что моки были вызваны правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 3: Ошибка при проверке пароля
	t.Run("LoginWrongPassword", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Подготавливаем тестового пользователя
		expectedUser := &models.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetUserByEmail", ctx, "test@example.com").Return(expectedUser, nil)
		mockRepo.On("CheckPassword", ctx, 1, "wrongpassword").Return(false, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод Login
		result, err := baseHandler.Login(ctx, &grpc.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = Invalid credentials", err.Error())

		// Проверяем, что моки были вызваны правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 4: Ошибка при получении пользователя из репозитория
	t.Run("LoginRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Определяем ожидаемое поведение мока - ошибка при получении пользователя
		mockRepo.On("GetUserByEmail", ctx, "test@example.com").Return((*models.User)(nil), errors.New("database error"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод Login
		result, err := baseHandler.Login(ctx, &grpc.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "rpc error: code = Unauthenticated desc = Invalid credentials", err.Error())

		// Проверяем, что моки были вызваны правильно
		mockRepo.AssertExpectations(t)
	})
}
