package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// roleRepository реализация интерфейса RoleRepository
type roleRepository struct {
	db *DB
}

// NewRoleRepository создает новый репозиторий ролей
func NewRoleRepository(db *DB) RoleRepository {
	return &roleRepository{db: db}
}

// CreateRole создает новую роль
func (r *roleRepository) CreateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	query := `INSERT INTO roles (name, code, description) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.GetConnection().QueryRow(ctx, query, role.Name, role.Code, role.Description).Scan(&role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return role, nil
}

// GetRoleByID получает роль по ID
func (r *roleRepository) GetRoleByID(ctx context.Context, id int) (*models.Role, error) {
	query := `SELECT id, name, code, description FROM roles WHERE id = $1`
	row := r.db.GetConnection().QueryRow(ctx, query, id)

	role := &models.Role{}
	err := row.Scan(&role.ID, &role.Name, &role.Code, &role.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// UpdateRole обновляет информацию о роли
func (r *roleRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	query := `UPDATE roles SET name = $1, code = $2, description = $3 WHERE id = $4`
	_, err := r.db.GetConnection().Exec(ctx, query, role.Name, role.Code, role.Description, role.ID)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

// DeleteRole удаляет роль
func (r *roleRepository) DeleteRole(ctx context.Context, id int) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.GetConnection().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// GetRoles получает список ролей с ограничением и смещением
func (r *roleRepository) GetRoles(ctx context.Context, limit, offset int) ([]*models.Role, error) {
	query := `SELECT id, name, code, description FROM roles ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.GetConnection().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
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

// GetRoleWithPermissions получает роль с привязанными правами
func (r *roleRepository) GetRoleWithPermissions(ctx context.Context, id int) (*models.Role, error) {
	// Получаем роль
	role, err := r.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Получаем все права, связанные с этой ролью
	query := `SELECT p.id, p.name, p.code, p.description
	          FROM permissions p
	          JOIN role_permissions rp ON p.id = rp.permission_id
	          WHERE rp.role_id = $1`
	rows, err := r.db.GetConnection().Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		permission := models.Permission{}
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Code, &permission.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate permissions: %w", err)
	}

	role.Permissions = permissions
	return role, nil
}

// AddRolePermissions добавляет права к роли
func (r *roleRepository) AddRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	// Проверяем, что роль существует
	_, err := r.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Проверяем, что все права существуют
	for _, permissionID := range permissionIDs {
		_, err := r.db.GetConnection().Query(ctx, `SELECT id FROM permissions WHERE id = $1`, permissionID)
		if err != nil {
			return fmt.Errorf("permission not found: %w", err)
		}
	}

	// Добавляем права к роли
	for _, permissionID := range permissionIDs {
		query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := r.db.GetConnection().Exec(ctx, query, roleID, permissionID)
		if err != nil {
			return fmt.Errorf("failed to add permission to role: %w", err)
		}
	}

	return nil
}

// DeleteRolePermissions удаляет права из роли
func (r *roleRepository) DeleteRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	// Проверяем, что роль существует
	_, err := r.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Удаляем права из роли
	for _, permissionID := range permissionIDs {
		query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
		_, err := r.db.GetConnection().Exec(ctx, query, roleID, permissionID)
		if err != nil {
			return fmt.Errorf("failed to delete permission from role: %w", err)
		}
	}

	return nil
}

// InitDB инициализирует таблицы в БД для ролей
func (r *roleRepository) InitDB() error {
	// Создание таблицы ролей
	query := `
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
)`
	_, err := r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize roles table: %w", err)
	}

	// Создание таблицы связи ролей и прав (many-to-many)
	query = `
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
)`
	_, err = r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize role_permissions table: %w", err)
	}

	log.Println("RoleRepository initialized successfully")
	return nil
}