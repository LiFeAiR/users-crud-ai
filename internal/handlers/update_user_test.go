package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
	"github.com/LiFeAiR/users-crud-ai/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUserHandler(t *testing.T) {
	// Test 1: Обновление пользователя с organization
	t.Run("UpdateUserWithOrganization", func(t *testing.T) {
		// Create a test request with valid JSON
		jsonData := `{"id":1,"name":"John Smith","email":"johnsmith@example.com","password":"newpassword123","organization":"Updated Org"}`
		req := httptest.NewRequest("PUT", "/api/user", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Создаем мок репозиторий
		mockRepo := new(MockUserRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("UpdateUser", mock.Anything).Return(nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler with nil server (for testing purposes)
		// In real application, server would be passed
		baseHandler.UpdateUser(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check the response body
		var responseUser models.User
		err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, "John Smith", responseUser.Name)
		assert.Equal(t, "johnsmith@example.com", responseUser.Email)
		assert.Equal(t, utils.Ptr("Updated Org"), responseUser.Organization)
		assert.Equal(t, 1, responseUser.ID) // Dummy ID
	})

	// Test 2: Обновление пользователя без organization (NULL)
	t.Run("UpdateUserWithoutOrganization", func(t *testing.T) {
		// Create a test request with valid JSON
		jsonData := `{"id":1,"name":"Jane Smith","email":"janesmith@example.com","password":"newpassword456"}`
		req := httptest.NewRequest("PUT", "/api/user", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Создаем мок репозиторий
		mockRepo := new(MockUserRepository)

		// Определяем ожидаемое поведение мока
		mockRepo.On("UpdateUser", mock.Anything).Return(nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler with nil server (for testing purposes)
		// In real application, server would be passed
		baseHandler.UpdateUser(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestUpdateUserHandler_InvalidJSON(t *testing.T) {
	// Create a test request with invalid JSON
	invalidJSON := `{"name":"John Doe","email":}`
	req := httptest.NewRequest("PUT", "/api/user", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	// Создаем мок репозиторий
	mockRepo := new(MockUserRepository)

	// Подготавливаем тестового пользователя
	testUser := &models.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Определяем ожидаемое поведение мока
	mockRepo.On("GetUserByID", 1).Return(testUser, nil)

	// Создаем базовый обработчик с моком
	baseHandler := &BaseHandler{
		userRepo: mockRepo,
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler with nil server (for testing purposes)
	// In real application, server would be passed
	baseHandler.UpdateUser(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
