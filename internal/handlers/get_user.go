package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// GetUser общий метод для получения пользователя по ID из query параметров
func (bh *BaseHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из query параметров
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID in query parameters", http.StatusBadRequest)
		return
	}

	// Конвертируем строку в целое число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Используем репозиторий для получения пользователя
	user, err := bh.userRepo.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			log.Printf("Failed to get user: %v", err)
		}
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send response
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
