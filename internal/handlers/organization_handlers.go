package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
)

// CreateOrganization создает новую организацию
func (bh *BaseHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var org models.Organization

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Используем репозиторий для создания организации
	dbOrg, err := bh.orgRepo.CreateOrganization(&org)
	if err != nil {
		http.Error(w, "Failed to create organization", http.StatusInternalServerError)
		return
	}

	org.ID = dbOrg.ID

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Send response
	if err := json.NewEncoder(w).Encode(org); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetOrganization получает организацию по ID
func (bh *BaseHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из query параметров
	iDStr := r.URL.Query().Get("id")
	if iDStr == "" {
		http.Error(w, "Missing organization ID in query parameters", http.StatusBadRequest)
		return
	}

	// Конвертируем строку в целое число
	id, err := strconv.Atoi(iDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	// Используем репозиторий для получения организации
	org, err := bh.orgRepo.GetOrganizationByID(id)
	if err != nil {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Send response
	if err := json.NewEncoder(w).Encode(org); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateOrganization обновляет информацию об организации
func (bh *BaseHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из query параметров
	iDStr := r.URL.Query().Get("id")
	if iDStr == "" {
		http.Error(w, "Missing organization ID in query parameters", http.StatusBadRequest)
		return
	}

	// Конвертируем строку в целое число
	id, err := strconv.Atoi(iDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	var org models.Organization

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Устанавливаем ID из URL
	org.ID = id

	// Используем репозиторий для обновления организации
	err = bh.orgRepo.UpdateOrganization(&org)
	if err != nil {
		http.Error(w, "Failed to update organization", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Send response
	if err := json.NewEncoder(w).Encode(org); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteOrganization удаляет организацию
func (bh *BaseHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из query параметров
	iDStr := r.URL.Query().Get("id")
	if iDStr == "" {
		http.Error(w, "Missing organization ID in query parameters", http.StatusBadRequest)
		return
	}

	// Конвертируем строку в целое число
	id, err := strconv.Atoi(iDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	// Используем репозиторий для удаления организации
	err = bh.orgRepo.DeleteOrganization(id)
	if err != nil {
		http.Error(w, "Failed to delete organization", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send response
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Organization deleted successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetOrganizations получает список организаций с пагинацией
func (bh *BaseHandler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// Используем репозиторий для получения организаций
	organizations, err := bh.orgRepo.GetOrganizations(limit, offset)
	if err != nil {
		http.Error(w, "Failed to get organizations", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Send response
	if err := json.NewEncoder(w).Encode(organizations); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
