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

func TestCreateUserHandler(t *testing.T) {
	// Test 1: Создание пользователя с organization
	t.Run("CreateUserWithOrganization", func(t *testing.T) {
		// Create a test request with valid JSON
		jsonData := `{"name":"John Doe","email":"john@example.com","password":"secret123","organization":"Test Org"}`
		req := httptest.NewRequest("POST", "/api/user", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Создаем мок репозиторий
		mockRepo := new(MockUserRepository)

		// Подготавливаем тестового пользователя
		testUser := &models.User{
			ID:           1,
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "secret123",
			Organization: utils.Ptr("Test Org"),
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("CreateUser", mock.Anything).Return(testUser, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler with nil server (for testing purposes)
		// In real application, server would be passed
		baseHandler.CreateUser(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Check the response body
		var responseUser models.User
		err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", responseUser.Name)
		assert.Equal(t, "john@example.com", responseUser.Email)
		assert.Equal(t, 1, responseUser.ID) // Dummy ID
	})

	// Test 2: Создание пользователя без organization (NULL)
	t.Run("CreateUserWithoutOrganization", func(t *testing.T) {
		// Create a test request with valid JSON
		jsonData := `{"name":"Jane Smith","email":"jane@example.com","password":"secret456"}`
		req := httptest.NewRequest("POST", "/api/user", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Создаем мок репозиторий
		mockRepo := new(MockUserRepository)

		// Подготавливаем тестового пользователя без organization
		testUser := &models.User{
			ID:           1,
			Name:         "Jane Smith",
			Email:        "jane@example.com",
			Password:     "secret456",
			Organization: nil,
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("CreateUser", mock.Anything).Return(testUser, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			userRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler with nil server (for testing purposes)
		// In real application, server would be passed
		baseHandler.CreateUser(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Check the response body
		var responseUser models.User
		err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, "Jane Smith", responseUser.Name)
		assert.Equal(t, "jane@example.com", responseUser.Email)
		assert.Equal(t, 1, responseUser.ID) // Dummy ID
	})
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	// Create a test request with invalid JSON
	invalidJSON := `{"name":"John Doe","email":}`
	req := httptest.NewRequest("POST", "/api/user", bytes.NewBufferString(invalidJSON))
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
	baseHandler.CreateUser(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
