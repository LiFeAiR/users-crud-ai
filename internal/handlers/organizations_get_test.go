package handlers

import (
	"context"
	"testing"

	"github.com/LiFeAiR/crud-ai/internal/handlers/mocks"
	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"github.com/stretchr/testify/assert"
)

// TestBaseHandler_GetOrganizations тестирует метод GetOrganizations базового обработчика
func TestBaseHandler_GetOrganizations(t *testing.T) {
	ctx := context.Background()
	// Test 1: Успешное получение организации
	t.Run("GetOrganizationsSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockOrganizationRepository)

		// Подготавливаем тестовую организацию
		testOrg := &models.Organization{
			ID:   1,
			Name: "Test Org",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetOrganizations", ctx, 10, 0).Return([]*models.Organization{testOrg}, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Вызываем метод GetOrganizations
		orgs, err := baseHandler.GetOrganizations(ctx, &grpc.ListRequest{})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, orgs)
		assert.Equal(t, int32(1), orgs.Data[0].Id)
		assert.Equal(t, "Test Org", orgs.Data[0].Name)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
