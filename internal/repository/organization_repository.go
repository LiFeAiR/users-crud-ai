package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// organizationRepository реализация интерфейса OrganizationRepository
type organizationRepository struct {
	db *DB
}

// NewOrganizationRepository создает новый репозиторий организаций
func NewOrganizationRepository(db *DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

// CreateOrganization создает новую организацию
func (r *organizationRepository) CreateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	query := `INSERT INTO organizations (name) VALUES ($1) RETURNING id`
	err := r.db.GetConnection().QueryRow(ctx, query, org.Name).Scan(&org.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	return org, nil
}

// GetOrganizationByID получает организацию по ID
func (r *organizationRepository) GetOrganizationByID(ctx context.Context, id int) (*models.Organization, error) {
	query := `SELECT id, name FROM organizations WHERE id = $1`
	row := r.db.GetConnection().QueryRow(ctx, query, id)

	org := &models.Organization{}
	err := row.Scan(&org.ID, &org.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}

// UpdateOrganization обновляет информацию об организации
func (r *organizationRepository) UpdateOrganization(ctx context.Context, org *models.Organization) error {
	query := `UPDATE organizations SET name = $1 WHERE id = $2`
	_, err := r.db.GetConnection().Exec(ctx, query, org.Name, org.ID)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}
	return nil
}

// DeleteOrganization удаляет организацию
func (r *organizationRepository) DeleteOrganization(ctx context.Context, id int) error {
	query := `DELETE FROM organizations WHERE id = $1`
	_, err := r.db.GetConnection().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

// GetOrganizations получает список организаций с ограничением и смещением
func (r *organizationRepository) GetOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error) {
	query := `SELECT id, name FROM organizations ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.GetConnection().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}
	defer rows.Close()

	var organizations []*models.Organization
	for rows.Next() {
		org := &models.Organization{}
		err := rows.Scan(&org.ID, &org.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}
		organizations = append(organizations, org)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate organizations: %w", err)
	}

	return organizations, nil
}

// GetOrganizationPermissions получает права организации
func (r *organizationRepository) GetOrganizationPermissions(
	ctx context.Context,
	organizationID int,
) ([]*models.Permission, error) {
	// TODO Сложно!!
	query := `SELECT DISTINCT p.id, p.name, p.code, p.description
			  FROM permissions p
				LEFT JOIN organization_permissions op ON p.id = op.permission_id and op.organization_id = $1
				LEFT JOIN organization_roles ro ON ro.organization_id = $1
				LEFT JOIN role_permissions rp ON rp.role_id = ro.role_id and p.id = op.permission_id
			  WHERE ro.organization_id = $1 or op.organization_id = $1
			  ORDER BY p.id`
	rows, err := r.db.GetConnection().Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization permissions: %w", err)
	}
	defer rows.Close()

	var permissions []*models.Permission
	for rows.Next() {
		permission := &models.Permission{}
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Code, &permission.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate permissions: %w", err)
	}

	return permissions, nil
}

// AddOrganizationPermissions добавляет права к организации
func (r *organizationRepository) AddOrganizationPermissions(ctx context.Context, organizationID int, permissionIDs []int) error {
	// Проверяем, что организация существует
	_, err := r.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	// Проверяем, что все права существуют
	for _, permissionID := range permissionIDs {
		_, err := r.db.GetConnection().Query(ctx, `SELECT id FROM permissions WHERE id = $1`, permissionID)
		if err != nil {
			return fmt.Errorf("permission not found: %w", err)
		}
	}

	// Добавляем права к организации
	for _, permissionID := range permissionIDs {
		query := `INSERT INTO organization_permissions (organization_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := r.db.GetConnection().Exec(ctx, query, organizationID, permissionID)
		if err != nil {
			return fmt.Errorf("failed to add permission to organization: %w", err)
		}
	}

	return nil
}

// DeleteOrganizationPermissions удаляет права из организации
func (r *organizationRepository) DeleteOrganizationPermissions(ctx context.Context, organizationID int, permissionIDs []int) error {
	// Проверяем, что организация существует
	_, err := r.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	// Удаляем права из организации
	for _, permissionID := range permissionIDs {
		query := `DELETE FROM organization_permissions WHERE organization_id = $1 AND permission_id = $2`
		_, err := r.db.GetConnection().Exec(ctx, query, organizationID, permissionID)
		if err != nil {
			return fmt.Errorf("failed to delete permission from organization: %w", err)
		}
	}

	return nil
}

// GetOrganizationRoles получает роли организации
func (r *organizationRepository) GetOrganizationRoles(ctx context.Context, organizationID int) ([]*models.Role, error) {
	query := `SELECT r.id, r.name, r.code, r.description
	          FROM roles r
	          JOIN organization_roles orl ON r.id = orl.role_id
	          WHERE orl.organization_id = $1`
	rows, err := r.db.GetConnection().Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization roles: %w", err)
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.ID, &role.Name, &role.Code, &role.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate roles: %w", err)
	}

	return roles, nil
}

// AddOrganizationRoles добавляет роли к организации
func (r *organizationRepository) AddOrganizationRoles(ctx context.Context, organizationID int, roleIDs []int) error {
	// Проверяем, что организация существует
	_, err := r.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	// Проверяем, что все роли существуют
	for _, roleID := range roleIDs {
		_, err := r.db.GetConnection().Query(ctx, `SELECT id FROM roles WHERE id = $1`, roleID)
		if err != nil {
			return fmt.Errorf("role not found: %w", err)
		}
	}

	// Добавляем роли к организации
	for _, roleID := range roleIDs {
		query := `INSERT INTO organization_roles (organization_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := r.db.GetConnection().Exec(ctx, query, organizationID, roleID)
		if err != nil {
			return fmt.Errorf("failed to add role to organization: %w", err)
		}
	}

	return nil
}

// DeleteOrganizationRoles удаляет роли из организации
func (r *organizationRepository) DeleteOrganizationRoles(ctx context.Context, organizationID int, roleIDs []int) error {
	// Проверяем, что организация существует
	_, err := r.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	// Удаляем роли из организации
	for _, roleID := range roleIDs {
		query := `DELETE FROM organization_roles WHERE organization_id = $1 AND role_id = $2`
		_, err := r.db.GetConnection().Exec(ctx, query, organizationID, roleID)
		if err != nil {
			return fmt.Errorf("failed to delete role from organization: %w", err)
		}
	}

	return nil
}

// InitDB инициализирует таблицы в БД для организаций
func (r *organizationRepository) InitDB() error {
	query := `
CREATE TABLE IF NOT EXISTS organizations (
id SERIAL PRIMARY KEY,
name VARCHAR(255) NOT NULL
)`
	_, err := r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Создание таблицы связи организаций и прав (many-to-many)
	query = `
CREATE TABLE IF NOT EXISTS organization_permissions (
organization_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
PRIMARY KEY (organization_id, permission_id)
)`
	_, err = r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize organization_permissions table: %w", err)
	}

	// Создание таблицы связи организаций и ролей (many-to-many)
	query = `
CREATE TABLE IF NOT EXISTS organization_roles (
organization_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
PRIMARY KEY (organization_id, role_id)
)`
	_, err = r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize organization_roles table: %w", err)
	}

	log.Println("OrganizationRepository initialized successfully")
	return nil
}
