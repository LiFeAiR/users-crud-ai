package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// permissionRepository реализация интерфейса PermissionRepository
type permissionRepository struct {
	db *DB
}

// NewPermissionRepository создает новый репозиторий прав
func NewPermissionRepository(db *DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// CreatePermission создает новое право
func (r *permissionRepository) CreatePermission(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	query := `INSERT INTO permissions (name, code, description) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.GetConnection().QueryRow(ctx, query, permission.Name, permission.Code, permission.Description).Scan(&permission.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}
	return permission, nil
}

// GetPermissionByID получает право по ID
func (r *permissionRepository) GetPermissionByID(ctx context.Context, id int) (*models.Permission, error) {
	query := `SELECT id, name, code, description FROM permissions WHERE id = $1`
	row := r.db.GetConnection().QueryRow(ctx, query, id)

	permission := &models.Permission{}
	err := row.Scan(&permission.ID, &permission.Name, &permission.Code, &permission.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return permission, nil
}

// UpdatePermission обновляет информацию о праве
func (r *permissionRepository) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	query := `UPDATE permissions SET name = $1, code = $2, description = $3 WHERE id = $4`
	_, err := r.db.GetConnection().Exec(ctx, query, permission.Name, permission.Code, permission.Description, permission.ID)
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}
	return nil
}

// DeletePermission удаляет право
func (r *permissionRepository) DeletePermission(ctx context.Context, id int) error {
	query := `DELETE FROM permissions WHERE id = $1`
	_, err := r.db.GetConnection().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

// GetPermissions получает список прав с ограничением и смещением
func (r *permissionRepository) GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	query := `SELECT id, name, code, description FROM permissions ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.GetConnection().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
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

// InitDB инициализирует таблицы в БД для прав
func (r *permissionRepository) InitDB() error {
	query := `
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
)`
	_, err := r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Println("PermissionRepository initialized successfully")
	return nil
}