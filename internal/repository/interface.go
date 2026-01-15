package repository

import "github.com/LiFeAiR/users-crud-ai/internal/models"

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	GetUsers(limit, offset int) ([]*models.User, error)
	InitDB() error
}

// OrganizationRepository интерфейс для работы с организациями
type OrganizationRepository interface {
	CreateOrganization(org *models.Organization) (*models.Organization, error)
	GetOrganizationByID(id int) (*models.Organization, error)
	UpdateOrganization(org *models.Organization) error
	DeleteOrganization(id int) error
	GetOrganizations(limit, offset int) ([]*models.Organization, error)
	InitDB() error
}
