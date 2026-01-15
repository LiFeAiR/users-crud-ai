package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteOrganizationHandler(t *testing.T) {
	// Test 1: Удаление организации
	t.Run("DeleteOrganizationSuccess", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("DELETE", "/api/organization?id=1", nil)

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("DeleteOrganization", 1).Return(nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.DeleteOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check the response body
		var response map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Organization deleted successfully", response["message"])
	})

	// Test 2: Некорректный ID
	t.Run("DeleteOrganizationInvalidID", func(t *testing.T) {
		// Create a test request with invalid ID
		req := httptest.NewRequest("DELETE", "/api/organization?id=abc", nil)

		// Создаем базовый обработчик (без мока, так как не должен вызываться)
		baseHandler := &BaseHandler{
			orgRepo: nil,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.DeleteOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test 3: Организация не найдена
	t.Run("DeleteOrganizationNotFound", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("DELETE", "/api/organization?id=1", nil)

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("DeleteOrganization", 1).Return(errors.New("organization not found"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.DeleteOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
