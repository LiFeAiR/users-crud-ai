package handlers

import (
	"github.com/LiFeAiR/users-crud-ai/internal/repository"
)

// BaseHandler базовый обработчик, который принимает репозитории
type BaseHandler struct {
	userRepo  repository.UserRepository
	orgRepo   repository.OrganizationRepository
}

// NewBaseHandler создает новый базовый обработчик
func NewBaseHandler(userRepo repository.UserRepository, orgRepo repository.OrganizationRepository) *BaseHandler {
	return &BaseHandler{
		userRepo:  userRepo,
		orgRepo:   orgRepo,
	}
}
