package handlers

import (
	"github.com/LiFeAiR/crud-ai/internal/repository"
)

// BaseHandler базовый обработчик, который принимает репозитории
type BaseHandler struct {
	userRepo   repository.UserRepository
	orgRepo    repository.OrganizationRepository
	permRepo   repository.PermissionRepository
	roleRepo   repository.RoleRepository
	tariffRepo repository.TariffRepository
	secretKey  string
}

// NewBaseHandler создает новый базовый обработчик
func NewBaseHandler(
	userRepo repository.UserRepository,
	orgRepo repository.OrganizationRepository,
	permRepo repository.PermissionRepository,
	roleRepo repository.RoleRepository,
	tariffRepo repository.TariffRepository,
	secretKey string,
) *BaseHandler {
	return &BaseHandler{
		userRepo:   userRepo,
		orgRepo:    orgRepo,
		permRepo:   permRepo,
		roleRepo:   roleRepo,
		tariffRepo: tariffRepo,
		secretKey:  secretKey,
	}
}

// convertInt32SliceToInt конвертирует slice int32 в slice int
func convertInt32SliceToInt(slice []int32) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[i] = int(v)
	}
	return result
}
