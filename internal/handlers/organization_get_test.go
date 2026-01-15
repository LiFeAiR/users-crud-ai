package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepository для тестирования
type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) CreateOrganization(org *models.Organization) (*models.Organization, error) {
	args := m.Called(org)

	if o, ok := args.Get(0).(*models.Organization); ok {
		return o, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockOrganizationRepository) GetOrganizationByID(id int) (*models.Organization, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) UpdateOrganization(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) DeleteOrganization(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrganizationRepository) GetOrganizations(limit, offset int) ([]*models.Organization, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*models.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) InitDB() error {
	panic("implement me")
}

func TestGetOrganizationHandler(t *testing.T) {
	// Test 1: Получение организации
	t.Run("GetOrganizationSuccess", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("GET", "/api/organization?id=1", nil)

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Подготавливаем тестовую организацию
		testOrg := &models.Organization{
			ID:   1,
			Name: "Test Org",
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetOrganizationByID", 1).Return(testOrg, nil)

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.GetOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check the response body
		var responseOrg models.Organization
		err := json.Unmarshal(rr.Body.Bytes(), &responseOrg)
		assert.NoError(t, err)
		assert.Equal(t, "Test Org", responseOrg.Name)
		assert.Equal(t, 1, responseOrg.ID)
	})

	// Test 2: Отсутствует ID
	t.Run("GetOrganizationMissingID", func(t *testing.T) {
		// Create a test request without ID
		req := httptest.NewRequest("GET", "/api/organization", nil)

		// Создаем базовый обработчик (без мока, так как не должен вызываться)
		baseHandler := &BaseHandler{
			orgRepo: nil,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.GetOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test 3: Некорректный ID
	t.Run("GetOrganizationInvalidID", func(t *testing.T) {
		// Create a test request with invalid ID
		req := httptest.NewRequest("GET", "/api/organization?id=abc", nil)

		// Создаем базовый обработчик (без мока, так как не должен вызываться)
		baseHandler := &BaseHandler{
			orgRepo: nil,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.GetOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test 4: Организация не найдена
	t.Run("GetOrganizationNotFound", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("GET", "/api/organization?id=1", nil)

		// Создаем мок репозиторий
		mockRepo := new(MockOrganizationRepository)

		// Определяем ожидаемое поведение мока - возвращаем ошибку
		mockRepo.On("GetOrganizationByID", 1).Return((*models.Organization)(nil), errors.New("organization not found"))

		// Создаем базовый обработчик с моком
		baseHandler := &BaseHandler{
			orgRepo: mockRepo,
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		baseHandler.GetOrganization(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
