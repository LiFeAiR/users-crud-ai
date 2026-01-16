package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/internal/utils"
	"github.com/jackc/pgx/v5"
)

// userRepository реализация интерфейса UserRepository
type userRepository struct {
	db *DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, organization_id) VALUES ($1, $2, $3) RETURNING id`
	var org interface{}
	if user.Organization != nil {
		org = user.Organization.ID
	}
	err := r.db.GetConnection().QueryRow(ctx, query, user.Name, user.Email, org).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUserByID получает пользователя по ID
func (r *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, organization_id FROM users WHERE id = $1`
	row := r.db.GetConnection().QueryRow(ctx, query, id)

	user := &models.User{}
	var org sql.NullInt32
	err := row.Scan(&user.ID, &user.Name, &user.Email, &org)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем, было ли значение NULL
	if org.Valid {
		user.Organization = &models.Organization{ID: int(org.Int32)}
	} else {
		user.Organization = nil
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2, organization_id = $3 WHERE id = $4`
	var orgId sql.NullInt32
	if user.Organization != nil {
		orgId = utils.NewNullInt32(int32(user.Organization.ID))
	}
	_, err := r.db.GetConnection().Exec(ctx, query, user.Name, user.Email, orgId, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser удаляет пользователя
func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.GetConnection().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetUsers получает список пользователей с ограничением и смещением
func (r *userRepository) GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `SELECT 
    				u.id, u.name, u.email, u.organization_id,
    				o.name AS organization
			  FROM users u
			  LEFT JOIN organizations o ON o.id = u.organization_id
			  ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.GetConnection().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var (
			orgId   sql.NullInt32
			orgName sql.NullString
		)
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &orgId, &orgName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		// Проверяем, было ли значение NULL
		if orgId.Valid && orgName.Valid {
			user.Organization = &models.Organization{ID: int(orgId.Int32), Name: orgName.String}
		} else {
			user.Organization = nil
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}

// InitDB инициализирует таблицы в БД
func (r *userRepository) InitDB() error {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
    organization_id integer
);
alter table users
    drop IF EXISTS organization;
alter table users
    add IF NOT EXISTS organization_id integer;
`
	_, err := r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Println("UserRepository initialized successfully")
	return nil
}
