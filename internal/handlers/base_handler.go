package handlers

import (
	"github.com/LiFeAiR/crud-ai/internal/repository"
)

// BaseHandler базовый обработчик, который принимает репозитории
type BaseHandler struct {
	userRepo  repository.UserRepository
	orgRepo   repository.OrganizationRepository
	permRepo  repository.PermissionRepository
	secretKey string
}

// NewBaseHandler создает новый базовый обработчик
func NewBaseHandler(
	userRepo repository.UserRepository,
	orgRepo repository.OrganizationRepository,
	permRepo repository.PermissionRepository,
	secretKey string,
) *BaseHandler {
	return &BaseHandler{
		userRepo:  userRepo,
		orgRepo:   orgRepo,
		permRepo:  permRepo,
		secretKey: secretKey,
	}
}
