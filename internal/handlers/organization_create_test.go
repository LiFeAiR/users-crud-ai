package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrganizationHandler(t *testing.T) {
	// Test 1: Создание организации
	t.Run("CreateOrganizationSuccess", func(t *testing.T) {
		// Create a test request with valid JSON
		jsonData := `{"name":"Test Org"}`
		req := httptest.NewRequest("POST", "/api/organization", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Подготавливаем тестовую организацию
		testOrg := &models.Organization{
			ID:   1,
			Name: "Test Org",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("CreateOrganization", mock.Anything).Return(testOrg, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.CreateOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Check the response body
		var responseOrg models.Organization
		err := json.Unmarshal(rr.Body.Bytes(), &responseOrg)
		assert.NoError(t, err)
		assert.Equal(t, "Test Org", responseOrg.Name)
		assert.Equal(t, 1, responseOrg.ID) // Dummy ID
	})

	// Test 2: Некорректный JSON
	t.Run("CreateOrganizationInvalidJSON", func(t *testing.T) {
		// Create a test request with invalid JSON
		invalidJSON := `{"name":}`
		req := httptest.NewRequest("POST", "/api/organization", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		// Создаем базовый обработчик (без мока, так как не должен вызываться)
		baseHandler := &BaseHandler{
			orgRepo: nil,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.CreateOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
