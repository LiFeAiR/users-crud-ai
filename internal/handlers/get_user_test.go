package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
	"github.com/LiFeAiR/users-crud-ai/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepository для тестирования
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsers(limit, offset int) ([]*models.User, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) InitDB() error {
	panic("implement me")
}

func (m *MockUserRepository) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)

	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestGetUserHandler_Success(t *testing.T) {
	// Создаем мок репозиторий
	mockRepo := new(MockUserRepository)

	// Подготавливаем тестового пользователя с organization
	testUser := &models.User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		Organization: utils.Ptr("Test Org"),
	}

	// Определяем ожидаемое поведение мока
	mockRepo.On("GetUserByID", 1).Return(testUser, nil)

	// Создаем базовый обработчик с моком
	baseHandler := &BaseHandler{
		userRepo: mockRepo,
	}

	// Создаем тестовый запрос с query параметром id
	req := httptest.NewRequest("GET", "/api/user?id=1", nil)

	// Создаем recorder для ответа
	rr := httptest.NewRecorder()

	// Вызываем метод GetUser
	baseHandler.GetUser(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем заголовки
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Проверяем тело ответа
	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, responseUser.ID)
	assert.Equal(t, testUser.Name, responseUser.Name)
	assert.Equal(t, testUser.Email, responseUser.Email)
	assert.Equal(t, testUser.Organization, responseUser.Organization)

	// Проверяем, что мок был вызван правильно
	mockRepo.AssertExpectations(t)
}

func TestGetUserHandler_SuccessWithNullOrganization(t *testing.T) {
	// Создаем мок репозиторий
	mockRepo := new(MockUserRepository)

	// Подготавливаем тестового пользователя без organization (nil)
	testUser := &models.User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		Organization: nil,
	}

	// Определяем ожидаемое поведение мока
	mockRepo.On("GetUserByID", 1).Return(testUser, nil)

	// Создаем базовый обработчик с моком
	baseHandler := &BaseHandler{
		userRepo: mockRepo,
	}

	// Создаем тестовый запрос с query параметром id
	req := httptest.NewRequest("GET", "/api/user?id=1", nil)

	// Создаем recorder для ответа
	rr := httptest.NewRecorder()

	// Вызываем метод GetUser
	baseHandler.GetUser(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем заголовки
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Проверяем тело ответа
	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, responseUser.ID)
	assert.Equal(t, testUser.Name, responseUser.Name)
	assert.Equal(t, testUser.Email, responseUser.Email)
	assert.Nil(t, responseUser.Organization)

	// Проверяем, что мок был вызван правильно
	mockRepo.AssertExpectations(t)
}

func TestGetUserHandler_MissingID(t *testing.T) {
	// Создаем базовый обработчик
	baseHandler := &BaseHandler{
		userRepo: nil, // Не нужен для этого теста
	}

	// Создаем тестовый запрос без query параметра id
	req := httptest.NewRequest("GET", "/api/user", nil)

	// Создаем recorder для ответа
	rr := httptest.NewRecorder()

	// Вызываем метод GetUser
	baseHandler.GetUser(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Проверяем сообщение об ошибке
	assert.Contains(t, rr.Body.String(), "Missing user ID in query parameters")
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	// Создаем базовый обработчик
	baseHandler := &BaseHandler{
		userRepo: nil, // Не нужен для этого теста
	}

	// Создаем тестовый запрос с некорректным ID
	req := httptest.NewRequest("GET", "/api/user?id=abc", nil)

	// Создаем recorder для ответа
	rr := httptest.NewRecorder()

	// Вызываем метод GetUser
	baseHandler.GetUser(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Проверяем сообщение об ошибке
	assert.Contains(t, rr.Body.String(), "Invalid user ID")
}

func TestGetUserHandler_UserNotFound(t *testing.T) {
	// Создаем мок репозиторий
	mockRepo := new(MockUserRepository)

	// Определяем ожидаемое поведение мока - пользователь не найден
	mockRepo.On("GetUserByID", 999).Return((*models.User)(nil), fmt.Errorf("user not found"))

	// Создаем базовый обработчик с моком
	baseHandler := &BaseHandler{
		userRepo: mockRepo,
	}

	// Создаем тестовый запрос с существующим ID, но несуществующим пользователем
	req := httptest.NewRequest("GET", "/api/user?id=999", nil)

	// Создаем recorder для ответа
	rr := httptest.NewRecorder()

	// Вызываем метод GetUser
	baseHandler.GetUser(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Проверяем сообщение об ошибке
	assert.Contains(t, rr.Body.String(), "User not found")

	// Проверяем, что мок был вызван правильно
	mockRepo.AssertExpectations(t)
}
