package repository

import (
	"context"

	"github.com/LiFeAiR/crud-ai/internal/models"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CheckPassword(ctx context.Context, userID int, password string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
	InitDB() error
}

// OrganizationRepository интерфейс для работы с организациями
type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error)
	GetOrganizationByID(ctx context.Context, id int) (*models.Organization, error)
	UpdateOrganization(ctx context.Context, org *models.Organization) error
	DeleteOrganization(ctx context.Context, id int) error
	GetOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error)
	InitDB() error
}

// PermissionRepository интерфейс для работы с правами
type PermissionRepository interface {
	CreatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	GetPermissionByID(ctx context.Context, id int) (*models.Permission, error)
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	DeletePermission(ctx context.Context, id int) error
	GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error)
	InitDB() error
}

// RoleRepository интерфейс для работы с ролями
type RoleRepository interface {
	CreateRole(ctx context.Context, role *models.Role) (*models.Role, error)
	GetRoleByID(ctx context.Context, id int) (*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id int) error
	GetRoles(ctx context.Context, limit, offset int) ([]*models.Role, error)
	GetRoleWithPermissions(ctx context.Context, id int) (*models.Role, error)
	AddRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error
	DeleteRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error
	InitDB() error
}
