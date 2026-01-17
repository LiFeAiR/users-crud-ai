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

// TestBaseHandler_GetOrganization тестирует метод GetOrganization базового обработчика
func TestBaseHandler_GetOrganization(t *testing.T) {
	ctx := context.Background()
	// Test 1: Успешное получение организации
	t.Run("GetOrganizationSuccess", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockOrganizationRepository)

		// Подготавливаем тестовую организацию
		testOrg := &models.Organization{
			ID:   1,
			Name: "Test Org",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetOrganizationByID", ctx, 1).Return(testOrg, nil)
		mockRepo.On("GetOrganizationPermissions", ctx, 1).Return([]*models.Permission{{
			ID:          0,
			Name:        "1",
			Code:        "2",
			Description: "3",
		}}, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Вызываем метод GetOrganization
		org, err := baseHandler.GetOrganization(ctx, &grpc.Id{Id: 1})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, int32(1), org.Id)
		assert.Equal(t, "Test Org", org.Name)
		assert.Len(t, org.Permissions, 1)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})

	// Test 2: Некорректный аргумент (nil)
	t.Run("GetOrganizationNilArgument", func(t *testing.T) {
		// Создаем базовый обработчик
		baseHandler := &BaseHandler{
			orgRepo: nil,
		}

		// Вызываем метод GetOrganization с nil аргументом
		org, err := baseHandler.GetOrganization(ctx, nil)

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, org)
	})

	// Test 3: Организация не найдена
	t.Run("GetOrganizationNotFound", func(t *testing.T) {
		// Создаем мок репозиторий
		mockRepo := new(mocks.MockOrganizationRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("GetOrganizationByID", ctx, 1).Return((*models.Organization)(nil), errors.New("organization not found"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Вызываем метод GetOrganization
		org, err := baseHandler.GetOrganization(ctx, &grpc.Id{Id: 1})

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, org)

		// Проверяем, что мок был вызван правильно
		mockRepo.AssertExpectations(t)
	})
}
