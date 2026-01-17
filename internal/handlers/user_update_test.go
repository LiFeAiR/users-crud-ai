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

// TestBaseHandler_UpdateUser тестирует метод UpdateUser базового обработчика
func TestBaseHandler_UpdateUser(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное обновление пользователя
	t.Run("UpdateUserSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		userRepo := new(mocks.MockUserRepository)
		orgRepo := new(mocks.MockOrganizationRepository)

		// Определяем ожидаемое поведение мока
		userRepo.On("UpdateUser", ctx, mock.Anything).Return(nil)
		var org *models.Organization
		orgRepo.On("GetOrganizationByID", ctx, mock.Anything).Return(org, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: userRepo,
			orgRepo:  orgRepo,
		}

		// Вызываем метод UpdateUser
		result, err := baseHandler.UpdateUser(ctx, &grpc.UserUpdateRequest{
			Id:             1,
			Name:           "Updated User",
			Email:          "updated@example.com",
			Password:       "newpassword123",
			OrganizationId: 1,
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Updated User", result.Name)
		assert.Equal(t, "updated@example.com", result.Email)
		assert.Nil(t, result.Organization)

		// Проверяем, что мок был вызван правильно
		userRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при отсутствии ID в запросе
	t.Run("UpdateUserMissingID", func(t *testing.T) {
		// Создаем базовый обработчик без мока (не нужен для этого теста)
		baseHandler := &BaseHandler{}

		// Вызываем метод UpdateUser с пустым ID
		result, err := baseHandler.UpdateUser(ctx, &grpc.UserUpdateRequest{
			Id:             0,
			Name:           "Updated User",
			Email:          "updated@example.com",
			Password:       "newpassword123",
			OrganizationId: 1,
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test 3: Ошибка при отсутствии самого запроса
	t.Run("UpdateUserNilRequest", func(t *testing.T) {
		// Создаем базовый обработчик без мока (не нужен для этого теста)
		baseHandler := &BaseHandler{}

		// Вызываем метод UpdateUser с nil запросом
		result, err := baseHandler.UpdateUser(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test 4: Ошибка при неудачном обновлении в репозитории
	t.Run("UpdateUserRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockUserRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("UpdateUser", ctx, mock.Anything).
			Return(errors.New("update failed"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Вызываем метод UpdateUser
		result, err := baseHandler.UpdateUser(ctx, &grpc.UserUpdateRequest{
			Id:             1,
			Name:           "Updated User",
			Email:          "updated@example.com",
			Password:       "newpassword123",
			OrganizationId: 1,
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}

// Вспомогательная функция для создания указателя на строку
func stringPtr(s string) *string {
	return &s
}
