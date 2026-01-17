package handlers

import (
	"context"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrganization(t *testing.T) {
	ctx := context.Background()

	// Test 1: Успешное создание организации
	t.Run("CreateOrganizationSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockOrganizationRepository)

		// Определяем ожидаемое поведение мока
		expectedOrg := &models.Organization{
			ID:   1,
			Name: "Test Organization",
		}
		mockRepo.On("CreateOrganization", ctx, &models.Organization{Name: "Test Organization"}).Return(expectedOrg, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Вызываем метод CreateOrganization
		result, err := baseHandler.CreateOrganization(ctx, &grpc.OrganizationCreateRequest{Name: "Test Organization"})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "Test Organization", result.Name)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Ошибка при пустом имени организации
	t.Run("CreateOrganizationEmptyName", func(t *testing.T) {
		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{}

		// Вызываем метод CreateOrganization с пустым именем
		result, err := baseHandler.CreateOrganization(ctx, &grpc.OrganizationCreateRequest{Name: ""})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// Test 3: Ошибка при nil запросе
	t.Run("CreateOrganizationNilRequest", func(t *testing.T) {
		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{}

		// Вызываем метод CreateOrganization с nil запросом
		result, err := baseHandler.CreateOrganization(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
