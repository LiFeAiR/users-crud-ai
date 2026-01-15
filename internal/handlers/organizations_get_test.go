package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetOrganizationsHandler(t *testing.T) {
	// Test 1: Получение списка организаций
	t.Run("GetOrganizationsSuccess", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("GET", "/api/organizations?limit=10&offset=0", nil)

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Подготавливаем тестовые организации
		testOrg1 := &models.Organization{
			ID:   1,
			Name: "Test Org 1",
		}
		testOrg2 := &models.Organization{
			ID:   2,
			Name: "Test Org 2",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetOrganizations", 10, 0).Return([]*models.Organization{testOrg1, testOrg2}, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.GetOrganizations(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test 2: Получение списка организаций с параметрами по умолчанию
	t.Run("GetOrganizationsDefaultParams", func(t *testing.T) {
		// Create a test request without parameters
		req := httptest.NewRequest("GET", "/api/organizations", nil)

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Подготавливаем тестовые организации
		testOrg := &models.Organization{
			ID:   1,
			Name: "Test Org",
		}

		// Определяем ожидаемое поведение мока - по умолчанию limit=10, offset=0
		mockRepo.On("GetOrganizations", 10, 0).Return([]*models.Organization{testOrg}, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.GetOrganizations(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test 3: Ошибка при получении организаций
	t.Run("GetOrganizationsError", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("GET", "/api/organizations?limit=10&offset=0", nil)

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("GetOrganizations", 10, 0).Return([]*models.Organization{}, errors.New("database error"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.GetOrganizations(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
