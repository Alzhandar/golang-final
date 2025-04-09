package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type MenuRepository struct {
	db *pgxpool.Pool
}

func NewMenuRepository(db *pgxpool.Pool) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) Create(ctx context.Context, menu *models.Menu) (int64, error) {
	query := `
        INSERT INTO menus (restaurant_id, name_ru, name_kz, img)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query,
		menu.RestaurantID,
		menu.NameRU,
		menu.NameKZ,
		menu.Img,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать меню: %w", err)
	}

	return id, nil
}

func (r *MenuRepository) GetByID(ctx context.Context, id int64) (*models.Menu, error) {
	query := `
        SELECT id, restaurant_id, name_ru, name_kz, img
        FROM menus
        WHERE id = $1
    `
	var menu models.Menu
	err := r.db.QueryRow(ctx, query, id).Scan(
		&menu.ID,
		&menu.RestaurantID,
		&menu.NameRU,
		&menu.NameKZ,
		&menu.Img,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("меню с ID %d не найдено", id)
		}
		return nil, fmt.Errorf("не удалось получить меню: %w", err)
	}

	return &menu, nil
}

func (r *MenuRepository) GetByRestaurant(ctx context.Context, restaurantID int64) ([]*models.Menu, error) {
	query := `
        SELECT id, restaurant_id, name_ru, name_kz, img
        FROM menus
        WHERE restaurant_id = $1
        ORDER BY name_ru
    `
	rows, err := r.db.Query(ctx, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить меню ресторана: %w", err)
	}
	defer rows.Close()

	var menus []*models.Menu
	for rows.Next() {
		var menu models.Menu
		if err := rows.Scan(
			&menu.ID,
			&menu.RestaurantID,
			&menu.NameRU,
			&menu.NameKZ,
			&menu.Img,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании меню: %w", err)
		}
		menus = append(menus, &menu)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по меню: %w", err)
	}

	return menus, nil
}

func (r *MenuRepository) Update(ctx context.Context, menu *models.Menu) error {
	query := `
        UPDATE menus
        SET restaurant_id = $1, name_ru = $2, name_kz = $3, img = $4
        WHERE id = $5
    `
	commandTag, err := r.db.Exec(ctx, query,
		menu.RestaurantID,
		menu.NameRU,
		menu.NameKZ,
		menu.Img,
		menu.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить меню: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("меню с ID %d не найдено", menu.ID)
	}

	return nil
}

func (r *MenuRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM menus WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить меню: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("меню с ID %d не найдено", id)
	}

	return nil
}
