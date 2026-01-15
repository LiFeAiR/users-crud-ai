package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
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
func (r *organizationRepository) CreateOrganization(org *models.Organization) (*models.Organization, error) {
	query := `INSERT INTO organizations (name) VALUES ($1) RETURNING id`
	err := r.db.GetConnection().QueryRow(context.Background(), query, org.Name).Scan(&org.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	return org, nil
}

// GetOrganizationByID получает организацию по ID
func (r *organizationRepository) GetOrganizationByID(id int) (*models.Organization, error) {
	query := `SELECT id, name FROM organizations WHERE id = $1`
	row := r.db.GetConnection().QueryRow(context.Background(), query, id)

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
func (r *organizationRepository) UpdateOrganization(org *models.Organization) error {
	query := `UPDATE organizations SET name = $1 WHERE id = $2`
	_, err := r.db.GetConnection().Exec(context.Background(), query, org.Name, org.ID)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}
	return nil
}

// DeleteOrganization удаляет организацию
func (r *organizationRepository) DeleteOrganization(id int) error {
	query := `DELETE FROM organizations WHERE id = $1`
	_, err := r.db.GetConnection().Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

// GetOrganizations получает список организаций с ограничением и смещением
func (r *organizationRepository) GetOrganizations(limit, offset int) ([]*models.Organization, error) {
	query := `SELECT id, name FROM organizations ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.GetConnection().Query(context.Background(), query, limit, offset)
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

	log.Println("OrganizationRepository initialized successfully")
	return nil
}
