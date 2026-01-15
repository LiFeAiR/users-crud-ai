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

func TestUpdateOrganizationHandler(t *testing.T) {
	// Test 1: Обновление организации
	t.Run("UpdateOrganizationSuccess", func(t *testing.T) {
		// Create a test request with valid JSON
		jsonData := `{"id":1,"name":"Updated Org"}`
		req := httptest.NewRequest("PUT", "/api/organization", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("UpdateOrganization", mock.Anything).Return(nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.UpdateOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check the response body
		var responseOrg models.Organization
		err := json.Unmarshal(rr.Body.Bytes(), &responseOrg)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Org", responseOrg.Name)
		assert.Equal(t, 1, responseOrg.ID)
	})

	// Test 2: Некорректный JSON
	t.Run("UpdateOrganizationInvalidJSON", func(t *testing.T) {
		// Create a test request with invalid JSON
		invalidJSON := `{"name":}`
		req := httptest.NewRequest("PUT", "/api/organization", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		// Создаем базовый обработчик (без мока, так как не должен вызываться)
		baseHandler := &BaseHandler{
			orgRepo: nil,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.UpdateOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
