package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
	"github.com/LiFeAiR/users-crud-ai/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersHandler(t *testing.T) {
	// Создаем мок репозиторий
	mockRepo := new(MockUserRepository)

	// Создаем базовый обработчик с моком
	baseHandler := &BaseHandler{
		userRepo: mockRepo,
	}

	// Тест 1: Получение списка пользователей
	t.Run("GetUsers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?limit=10&offset=0", nil)

		// Подготавливаем тестовых пользователей
		testUserWithOrg := &models.User{
			ID:           1,
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "secret123",
			Organization: utils.Ptr("Test Org"),
		}

		testUserWithoutOrg := &models.User{
			ID:           2,
			Name:         "Jane Smith",
			Email:        "jane@example.com",
			Password:     "secret456",
			Organization: nil,
		}

		// Определяем ожидаемое поведение мока
		mockRepo.On("GetUsers", 10, 0).Return([]*models.User{testUserWithOrg, testUserWithoutOrg}, nil)

		// Создаем recorder для ответа
		rr := httptest.NewRecorder()

		// Вызываем метод GetUser
		baseHandler.GetUsers(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
