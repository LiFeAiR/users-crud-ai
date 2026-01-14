package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LiFeAiR/users-crud-ai/internal/models"
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
func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, organization) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.GetConnection().QueryRow(context.Background(), query, user.Name, user.Email, user.Organization).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUserByID получает пользователя по ID
func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, name, email, organization FROM users WHERE id = $1`
	row := r.db.GetConnection().QueryRow(context.Background(), query, id)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Organization)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func (r *userRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2, organization = $3 WHERE id = $4`
	_, err := r.db.GetConnection().Exec(context.Background(), query, user.Name, user.Email, user.Organization, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser удаляет пользователя
func (r *userRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.GetConnection().Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetUsers получает список пользователей с ограничением и смещением
func (r *userRepository) GetUsers(limit, offset int) ([]*models.User, error) {
	query := `SELECT id, name, email, organization FROM users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.GetConnection().Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Organization)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
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
	organization VARCHAR(255)
)`
	_, err := r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}
