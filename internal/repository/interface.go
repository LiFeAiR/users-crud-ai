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
	GetUserPermissions(ctx context.Context, userID int) ([]*models.Permission, error)
	AddUserPermissions(ctx context.Context, userID int, permissionIDs []int) error
	DeleteUserPermissions(ctx context.Context, userID int, permissionIDs []int) error
	GetUserRoles(ctx context.Context, userID int) ([]*models.Role, error)
	AddUserRoles(ctx context.Context, userID int, roleIDs []int) error
	DeleteUserRoles(ctx context.Context, userID int, roleIDs []int) error
	InitDB() error
}

// OrganizationRepository интерфейс для работы с организациями
type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error)
	GetOrganizationByID(ctx context.Context, id int) (*models.Organization, error)
	UpdateOrganization(ctx context.Context, org *models.Organization) error
	DeleteOrganization(ctx context.Context, id int) error
	GetOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error)
	GetOrganizationPermissions(ctx context.Context, organizationID int) ([]*models.Permission, error)
	AddOrganizationPermissions(ctx context.Context, organizationID int, permissionIDs []int) error
	DeleteOrganizationPermissions(ctx context.Context, organizationID int, permissionIDs []int) error
	GetOrganizationRoles(ctx context.Context, organizationID int) ([]*models.Role, error)
	AddOrganizationRoles(ctx context.Context, organizationID int, roleIDs []int) error
	DeleteOrganizationRoles(ctx context.Context, organizationID int, roleIDs []int) error
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

// TariffRepository интерфейс для работы с тарифами
type TariffRepository interface {
	CreateTariff(ctx context.Context, tariff *models.Tariff) (*models.Tariff, error)
	GetTariffByID(ctx context.Context, id int) (*models.Tariff, error)
	UpdateTariff(ctx context.Context, tariff *models.Tariff) error
	DeleteTariff(ctx context.Context, id int) error
	GetTariffs(ctx context.Context, limit, offset int) ([]*models.Tariff, error)
	GetTariffWithRoles(ctx context.Context, id int) (*models.Tariff, error)
	AddTariffRoles(ctx context.Context, tariffID int, roleIDs []int) error
	DeleteTariffRoles(ctx context.Context, tariffID int, roleIDs []int) error
	InitDB() error
}
