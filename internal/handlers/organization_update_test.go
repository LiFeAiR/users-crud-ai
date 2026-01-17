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

// TestBaseHandler_UpdateOrganization тестирует метод UpdateOrganization базового обработчика
func TestBaseHandler_UpdateOrganization(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное обновление организации
	t.Run("UpdateOrganizationSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockOrganizationRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("UpdateOrganization", ctx, &models.Organization{ID: 1, Name: "Updated Org"}).Return(nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Вызываем метод UpdateOrganization
		result, err := baseHandler.UpdateOrganization(ctx, &grpc.OrganizationUpdateRequest{
			Id:   1,
			Name: "Updated Org",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Updated Org", result.Name)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при пустом ID организации
	t.Run("UpdateOrganizationEmptyId", func(t *testing.T) {
		// Создаем базовый обработчик
		baseHandler := &BaseHandler{}

		// Вызываем метод UpdateOrganization с пустым ID
		result, err := baseHandler.UpdateOrganization(ctx, &grpc.OrganizationUpdateRequest{
			Id:   0,
			Name: "Updated Org",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test 3: Ошибка при nil запросе
	t.Run("UpdateOrganizationNilRequest", func(t *testing.T) {
		// Создаем базовый обработчик
		baseHandler := &BaseHandler{}

		// Вызываем метод UpdateOrganization с nil запросом
		result, err := baseHandler.UpdateOrganization(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test 4: Ошибка при неудачном обновлении в репозитории
	t.Run("UpdateOrganizationRepositoryError", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockOrganizationRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("UpdateOrganization", ctx, &models.Organization{ID: 1, Name: "Updated Org"}).Return(errors.New("update failed"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Вызываем метод UpdateOrganization
		result, err := baseHandler.UpdateOrganization(ctx, &grpc.OrganizationUpdateRequest{
			Id:   1,
			Name: "Updated Org",
		})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
