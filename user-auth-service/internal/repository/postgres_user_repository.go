package repository

import (
	"awesomeProject4/user-auth-service/internal/db/postgres"
	"awesomeProject4/user-auth-service/internal/domains/models"
	"context"
	"fmt"
)

type PostgresUserRepository struct {
	db *postgres.PostgresDB
}

func NewPostgresUserRepository(postgresDB *postgres.PostgresDB) *PostgresUserRepository {
	return &PostgresUserRepository{db: postgresDB}
}

func (repo *PostgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Валидация пользователя
	if err := user.Validate(); err != nil {
		return err
	}
	// Хэширование пароля
	if err := user.HashPassword(); err != nil {
		return err
	}
	query := `
		INSERT INTO users (username, email, password,  status, created_at)
		VALUES ($1, $2, $3, $4,  NOW()) RETURNING id`
	err := repo.db.DB.QueryRowxContext(ctx, query, user.Username, user.Email, user.Password, user.Status).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return nil
}

// GetUserByID получает пользователя по его ID
func (repo *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE id = $1`

	var user models.User
	err := repo.db.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// IsUserExists проверяет, существует ли пользователь с указанным email
func (repo *PostgresUserRepository) IsUserExists(ctx context.Context, user models.User) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)`

	var exists bool
	err := repo.db.DB.GetContext(ctx, &exists, query, user.Email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *PostgresUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users WHERE email = $1`
	var user *models.User
	err := repo.db.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *PostgresUserRepository) FindUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users WHERE id = $1`
	var user *models.User
	err := repo.db.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repo *PostgresUserRepository) UpdateUserStatus(ctx context.Context, userID int, status models.Status) error {
	query := `
		UPDATE users SET status = $1 WHERE id = $2`
	result, err := repo.db.DB.ExecContext(ctx, query, status, userID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса пользователя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", userID)
	}
	return nil
}

// UpdateUserPassword обновляет пароль пользователя
func (repo *PostgresUserRepository) UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error {
	query := `
		UPDATE users SET password = $1 WHERE id = $2`
	result, err := repo.db.DB.ExecContext(ctx, query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пароля пользователя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", userID)
	}
	return nil
}
