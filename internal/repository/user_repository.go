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

// CheckPassword проверяет пароль пользователя по хешу
func (r *userRepository) CheckPassword(ctx context.Context, userID int, password string) (bool, error) {
	query := `SELECT password_hash FROM users WHERE id = $1`
	var hash string
	err := r.db.GetConnection().QueryRow(ctx, query, userID).Scan(&hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user password hash: %w", err)
	}

	// Используем функцию из utils для сравнения
	return utils.CheckPassword(password, hash), nil
}

// GetUserByEmail получает пользователя по email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, organization_id FROM users WHERE email = $1`
	row := r.db.GetConnection().QueryRow(ctx, query, email)

	user := &models.User{}
	var org sql.NullInt32
	err := row.Scan(&user.ID, &user.Name, &user.Email, &org)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Проверяем, было ли значение NULL
	if org.Valid {
		user.Organization = &models.Organization{ID: int(org.Int32)}
	} else {
		user.Organization = nil
	}

	return user, nil
}

// CreateUser создает нового пользователя
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, password_hash, organization_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var org interface{}
	if user.Organization != nil {
		org = user.Organization.ID
	}
	err := r.db.GetConnection().QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash, org).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUserByID получает пользователя по ID
func (r *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, organization_id, tariff_id FROM users WHERE id = $1`
	row := r.db.GetConnection().QueryRow(ctx, query, id)

	user := &models.User{}
	var (
		org    sql.NullInt32
		tariff sql.NullInt32
	)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &org, &tariff)
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

	// Проверяем, было ли значение NULL
	if tariff.Valid {
		user.TariffID = utils.Ptr(int(tariff.Int32))
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	// Если пароль не пустой, обновляем и хеш пароля
	if user.PasswordHash != "" {
		query := `UPDATE users SET name = $1, email = $2, password_hash = $3, organization_id = $4 WHERE id = $5`
		var orgId sql.NullInt32
		if user.Organization != nil {
			orgId = utils.NewNullInt32(int32(user.Organization.ID))
		}
		_, err := r.db.GetConnection().Exec(ctx, query, user.Name, user.Email, user.PasswordHash, orgId, user.ID)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		// Если пароль не указан, обновляем только остальные поля
		query := `UPDATE users SET name = $1, email = $2, organization_id = $3 WHERE id = $4`
		var orgId sql.NullInt32
		if user.Organization != nil {
			orgId = utils.NewNullInt32(int32(user.Organization.ID))
		}
		_, err := r.db.GetConnection().Exec(ctx, query, user.Name, user.Email, orgId, user.ID)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
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

// GetUserPermissions получает права пользователя
func (r *userRepository) GetUserPermissions(ctx context.Context, userID int) ([]*models.Permission, error) {
	// TODO Сложно!!
	query := `
		SELECT DISTINCT p.id, p.name, p.code, p.description
		FROM permissions p
			LEFT JOIN user_permissions up ON p.id = up.permission_id and up.user_id = $1
			LEFT JOIN user_roles ur ON ur.user_id = $1
			LEFT JOIN role_permissions rp ON rp.role_id = ur.role_id and p.id = rp.permission_id
		WHERE ur.user_id = $1 or up.user_id = $1
		ORDER BY p.id`
	rows, err := r.db.GetConnection().Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
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

// AddUserPermissions добавляет права пользователю
func (r *userRepository) AddUserPermissions(ctx context.Context, userID int, permissionIDs []int) error {
	// Проверяем, существуют ли все указанные права
	query := `SELECT COUNT(*) FROM permissions WHERE id = ANY($1)`
	var count int
	err := r.db.GetConnection().QueryRow(ctx, query, permissionIDs).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check permissions: %w", err)
	}

	// Если количество найденных прав не совпадает с количеством запрошенных, значит некоторые права не существуют
	if count != len(permissionIDs) {
		return fmt.Errorf("some permissions do not exist")
	}

	// Добавляем права пользователю
	for _, permissionID := range permissionIDs {
		query = `INSERT INTO user_permissions (user_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := r.db.GetConnection().Exec(ctx, query, userID, permissionID)
		if err != nil {
			return fmt.Errorf("failed to add permission to user: %w", err)
		}
	}

	return nil
}

// DeleteUserPermissions удаляет права у пользователя
func (r *userRepository) DeleteUserPermissions(ctx context.Context, userID int, permissionIDs []int) error {
	// Удаляем права у пользователя
	query := `DELETE FROM user_permissions WHERE user_id = $1 AND permission_id = ANY($2)`
	_, err := r.db.GetConnection().Exec(ctx, query, userID, permissionIDs)
	if err != nil {
		return fmt.Errorf("failed to delete user permissions: %w", err)
	}

	return nil
}

// GetUserRoles получает роли пользователя
func (r *userRepository) GetUserRoles(ctx context.Context, userID int) ([]*models.Role, error) {
	query := `
		SELECT r.id, r.name, r.code, r.description
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.id
	`
	rows, err := r.db.GetConnection().Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
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

// AddUserRoles добавляет роли пользователю
func (r *userRepository) AddUserRoles(ctx context.Context, userID int, roleIDs []int) error {
	// Проверяем, существуют ли все указанные роли
	query := `SELECT COUNT(*) FROM roles WHERE id = ANY($1)`
	var count int
	err := r.db.GetConnection().QueryRow(ctx, query, roleIDs).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check roles: %w", err)
	}

	// Если количество найденных ролей не совпадает с количеством запрошенных, значит некоторые роли не существуют
	if count != len(roleIDs) {
		return fmt.Errorf("some roles do not exist")
	}

	// Добавляем роли пользователю
	for _, roleID := range roleIDs {
		query = `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := r.db.GetConnection().Exec(ctx, query, userID, roleID)
		if err != nil {
			return fmt.Errorf("failed to add role to user: %w", err)
		}
	}

	return nil
}

// DeleteUserRoles удаляет роли у пользователя
func (r *userRepository) DeleteUserRoles(ctx context.Context, userID int, roleIDs []int) error {
	// Удаляем роли у пользователя
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = ANY($2)`
	_, err := r.db.GetConnection().Exec(ctx, query, userID, roleIDs)
	if err != nil {
		return fmt.Errorf("failed to delete user roles: %w", err)
	}

	return nil
}

// SetUserTariff устанавливает тариф пользователю
func (r *userRepository) SetUserTariff(ctx context.Context, userID int, tariffID *int32) error {
	var tariffIDVal interface{}
	if tariffID != nil {
		tariffIDVal = *tariffID
	} else {
		tariffIDVal = nil
	}

	query := `UPDATE users SET tariff_id = $1 WHERE id = $2`
	_, err := r.db.GetConnection().Exec(ctx, query, tariffIDVal, userID)
	if err != nil {
		return fmt.Errorf("failed to set user tariff: %w", err)
	}

	return nil
}

// GetUserTariff получает тариф пользователя
func (r *userRepository) GetUserTariff(ctx context.Context, userID int) (*models.Tariff, error) {
	query := `SELECT t.id, t.name, t.description, t.price
	          FROM users u
	          LEFT JOIN tariffs t ON u.tariff_id = t.id
	          WHERE u.id = $1`

	row := r.db.GetConnection().QueryRow(ctx, query, userID)

	tariff := &models.Tariff{}
	err := row.Scan(&tariff.ID, &tariff.Name, &tariff.Description, &tariff.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user tariff: %w", err)
	}

	return tariff, nil
}

// InitDB инициализирует таблицы в БД
func (r *userRepository) InitDB() error {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash TEXT,
    organization_id integer,
    tariff_id integer
);
alter table users
    drop IF EXISTS organization;
alter table users
    add IF NOT EXISTS organization_id integer;
alter table users
    add IF NOT EXISTS password_hash TEXT;
alter table users
    add IF NOT EXISTS tariff_id integer;

-- Таблица для связи пользователей и прав
CREATE TABLE IF NOT EXISTS user_permissions (
	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, permission_id)
);

-- Таблица для связи пользователей и ролей
CREATE TABLE IF NOT EXISTS user_roles (
	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, role_id)
);
`
	_, err := r.db.GetConnection().Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Println("UserRepository initialized successfully")
	return nil
}
