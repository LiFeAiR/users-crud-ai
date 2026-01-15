package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// GetUsersHandler обработчик для получения списка пользователей
func (h *BaseHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры limit и offset из запроса
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Устанавливаем значения по умолчанию
	limit := 10
	offset := 0

	// Парсим limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	// Парсим offset
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	// Получаем список пользователей из репозитория
	users, err := h.userRepo.GetUsers(limit, offset)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Response{
		Data: users,
	})
}

// Response структура для ответа API
type Response struct {
	Data interface{} `json:"data"`
}
