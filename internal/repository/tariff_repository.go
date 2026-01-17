package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// tariffRepository реализация интерфейса TariffRepository
type tariffRepository struct {
	db *DB
}

// NewTariffRepository создает новый репозиторий тарифов
func NewTariffRepository(db *DB) TariffRepository {
	return &tariffRepository{db: db}
}

// CreateTariff создает новый тариф
func (r *tariffRepository) CreateTariff(ctx context.Context, tariff *models.Tariff) (*models.Tariff, error) {
	query := `INSERT INTO tariffs (name, description, price) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.GetConnection().QueryRow(ctx, query, tariff.Name, tariff.Description, tariff.Price).Scan(&tariff.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create tariff: %w", err)
	}
	return tariff, nil
}

// GetTariffByID получает тариф по ID
func (r *tariffRepository) GetTariffByID(ctx context.Context, id int) (*models.Tariff, error) {
	query := `SELECT id, name, description, price FROM tariffs WHERE id = $1`
	row := r.db.GetConnection().QueryRow(ctx, query, id)

	tariff := &models.Tariff{}
	err := row.Scan(&tariff.ID, &tariff.Name, &tariff.Description, &tariff.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("tariff not found")
		}
		return nil, fmt.Errorf("failed to get tariff: %w", err)
	}

	return tariff, nil
}

// UpdateTariff обновляет информацию о тарифе
func (r *tariffRepository) UpdateTariff(ctx context.Context, tariff *models.Tariff) error {
	query := `UPDATE tariffs SET name = $1, description = $2, price = $3 WHERE id = $4`
	_, err := r.db.GetConnection().Exec(ctx, query, tariff.Name, tariff.Description, tariff.Price, tariff.ID)
	if err != nil {
		return fmt.Errorf("failed to update tariff: %w", err)
	}
	return nil
}

// DeleteTariff удаляет тариф
func (r *tariffRepository) DeleteTariff(ctx context.Context, id int) error {
	query := `DELETE FROM tariffs WHERE id = $1`
	_, err := r.db.GetConnection().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tariff: %w", err)
	}
	return nil
}

// GetTariffs получает список тарифов с ограничением и смещением
func (r *tariffRepository) GetTariffs(ctx context.Context, limit, offset int) ([]*models.Tariff, error) {
	query := `SELECT id, name, description, price FROM tariffs ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.GetConnection().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tariffs: %w", err)
	}
	defer rows.Close()

	var tariffs []*models.Tariff
	for rows.Next() {
		tariff := &models.Tariff{}
		err := rows.Scan(&tariff.ID, &tariff.Name, &tariff.Description, &tariff.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tariff: %w", err)
		}
		tariffs = append(tariffs, tariff)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tariffs: %w", err)
	}

	return tariffs, nil
}

// GetTariffWithRoles получает тариф с привязанными ролями
func (r *tariffRepository) GetTariffWithRoles(ctx context.Context, id int) (*models.Tariff, error) {
	// Получаем тариф
	tariff, err := r.GetTariffByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Получаем все роли, связанные с этим тарифом
	query := `SELECT r.id, r.name, r.code, r.description
	          FROM roles r
	          JOIN tariff_roles tr ON r.id = tr.role_id
	          WHERE tr.tariff_id = $1`
	rows, err := r.db.GetConnection().Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tariff roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		role := models.Role{}
		err := rows.Scan(&role.ID, &role.Name, &role.Code, &role.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate roles: %w", err)
	}

	tariff.Roles = roles
	return tariff, nil
}

// AddTariffRoles добавляет роли к тарифу
func (r *tariffRepository) AddTariffRoles(ctx context.Context, tariffID int, roleIDs []int) error {
	// Проверяем, что тариф существует
	_, err := r.GetTariffByID(ctx, tariffID)
	if err != nil {
		return fmt.Errorf("tariff not found: %w", err)
	}

	// Проверяем, что все роли существуют
	for _, roleID := range roleIDs {
		_, err := r.db.GetConnection().Query(ctx, `SELECT id FROM roles WHERE id = $1`, roleID)
		if err != nil {
			return fmt.Errorf("role not found: %w", err)
		}
	}

	// Добавляем роли к тарифу
	for _, roleID := range roleIDs {
		query := `INSERT INTO tariff_roles (tariff_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := r.db.GetConnection().Exec(ctx, query, tariffID, roleID)
		if err != nil {
			return fmt.Errorf("failed to add role to tariff: %w", err)
		}
	}

	return nil
}

// DeleteTariffRoles удаляет роли из тарифа
func (r *tariffRepository) DeleteTariffRoles(ctx context.Context, tariffID int, roleIDs []int) error {
	// Проверяем, что тариф существует
	_, err := r.GetTariffByID(ctx, tariffID)
	if err != nil {
		return fmt.Errorf("tariff not found: %w", err)
	}

	// Удаляем роли из тарифа
	for _, roleID := range roleIDs {
		query := `DELETE FROM tariff_roles WHERE tariff_id = $1 AND role_id = $2`
		_, err := r.db.GetConnection().Exec(ctx, query, tariffID, roleID)
		if err != nil {
			return fmt.Errorf("failed to delete role from tariff: %w", err)
		}
	}

	return nil
}

// InitDB инициализирует таблицы в БД для тарифов
func (r *tariffRepository) InitDB() error {
	// Создание таблицы тарифов
	query := `
CREATE TABLE IF NOT EXISTS tariffs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL
)`
	_, err := r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize tariffs table: %w", err)
	}

	// Создание таблицы связи тарифов и ролей (many-to-many)
	query = `
CREATE TABLE IF NOT EXISTS tariff_roles (
    tariff_id INTEGER REFERENCES tariffs(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (tariff_id, role_id)
)`
	_, err = r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize tariff_roles table: %w", err)
	}

	log.Println("TariffRepository initialized successfully")
	return nil
}