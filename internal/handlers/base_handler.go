package handlers

import (
	"github.com/LiFeAiR/crud-ai/internal/repository"
	"github.com/LiFeAiR/crud-ai/internal/utils"
)

// BaseHandler базовый обработчик, который принимает репозитории
type BaseHandler struct {
	userRepo   repository.UserRepository
	orgRepo    repository.OrganizationRepository
	permRepo   repository.PermissionRepository
	roleRepo   repository.RoleRepository
	tariffRepo repository.TariffRepository
	secretKey  string

	jwtFunc func(string, int, string, string, []string) (string, error)
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
		jwtFunc:    utils.GenerateJWT,
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
