package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type MenuTypeRepository struct {
	db *pgxpool.Pool
}

func NewMenuTypeRepository(db *pgxpool.Pool) *MenuTypeRepository {
	return &MenuTypeRepository{db: db}
}

func (r *MenuTypeRepository) Create(ctx context.Context, menuType *models.MenuType) (int64, error) {
	query := `
        INSERT INTO menu_types (name, img)
        VALUES ($1, $2)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query, menuType.Name, menuType.Img).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать тип меню: %w", err)
	}

	return id, nil
}

func (r *MenuTypeRepository) GetByID(ctx context.Context, id int64) (*models.MenuType, error) {
	query := `
        SELECT id, name, img
        FROM menu_types
        WHERE id = $1
    `
	var menuType models.MenuType
	err := r.db.QueryRow(ctx, query, id).Scan(
		&menuType.ID,
		&menuType.Name,
		&menuType.Img,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("тип меню с ID %d не найден", id)
		}
		return nil, fmt.Errorf("не удалось получить тип меню: %w", err)
	}

	return &menuType, nil
}

func (r *MenuTypeRepository) Update(ctx context.Context, menuType *models.MenuType) error {
	query := `
        UPDATE menu_types
        SET name = $1, img = $2
        WHERE id = $3
    `
	commandTag, err := r.db.Exec(ctx, query, menuType.Name, menuType.Img, menuType.ID)

	if err != nil {
		return fmt.Errorf("не удалось обновить тип меню: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("тип меню с ID %d не найден", menuType.ID)
	}

	return nil
}

func (r *MenuTypeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM menu_types WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить тип меню: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("тип меню с ID %d не найден", id)
	}

	return nil
}

func (r *MenuTypeRepository) List(ctx context.Context) ([]*models.MenuType, error) {
	query := `
        SELECT id, name, img
        FROM menu_types
        ORDER BY name
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список типов меню: %w", err)
	}
	defer rows.Close()

	var menuTypes []*models.MenuType
	for rows.Next() {
		var menuType models.MenuType
		if err := rows.Scan(
			&menuType.ID,
			&menuType.Name,
			&menuType.Img,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании типа меню: %w", err)
		}
		menuTypes = append(menuTypes, &menuType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по типам меню: %w", err)
	}

	return menuTypes, nil
}
