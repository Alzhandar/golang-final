package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	query := `
        INSERT INTO users (phone_number, name, last_name, language, is_active)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query, user.PhoneNumber, user.Name, user.LastName,
		user.Language, user.IsActive).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать пользователя: %w", err)
	}

	return id, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
        SELECT id, phone_number, name, last_name, language, is_active
        FROM users
        WHERE id = $1
    `
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Name,
		&user.LastName,
		&user.Language,
		&user.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с ID %d не найден", id)
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
        SELECT id, phone_number, name, last_name, language, is_active
        FROM users
        WHERE phone_number = $1
    `
	var user models.User
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Name,
		&user.LastName,
		&user.Language,
		&user.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с номером телефона %s не найден", phone)
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET phone_number = $1, name = $2, last_name = $3, language = $4, is_active = $5
        WHERE id = $6
    `
	commandTag, err := r.db.Exec(ctx, query,
		user.PhoneNumber,
		user.Name,
		user.LastName,
		user.Language,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить пользователя: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", user.ID)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить пользователя: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", id)
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
        SELECT id, phone_number, name, last_name, language, is_active
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список пользователей: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.PhoneNumber,
			&user.Name,
			&user.LastName,
			&user.Language,
			&user.IsActive,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании пользователя: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по пользователям: %w", err)
	}

	return users, nil
}
