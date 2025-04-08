package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type SectionRepository struct {
	db *pgxpool.Pool
}

func NewSectionRepository(db *pgxpool.Pool) *SectionRepository {
	return &SectionRepository{db: db}
}

func (r *SectionRepository) Create(ctx context.Context, section *models.Section) (int64, error) {
	query := `
        INSERT INTO sections (restaurant_id, name)
        VALUES ($1, $2)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query, section.RestaurantID, section.Name).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать секцию: %w", err)
	}

	return id, nil
}

func (r *SectionRepository) GetByID(ctx context.Context, id int64) (*models.Section, error) {
	query := `
        SELECT id, restaurant_id, name
        FROM sections
        WHERE id = $1
    `
	var section models.Section
	err := r.db.QueryRow(ctx, query, id).Scan(
		&section.ID,
		&section.RestaurantID,
		&section.Name,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("секция с ID %d не найдена", id)
		}
		return nil, fmt.Errorf("не удалось получить секцию: %w", err)
	}

	return &section, nil
}

func (r *SectionRepository) GetByRestaurant(ctx context.Context, restaurantID int64) ([]*models.Section, error) {
	query := `
        SELECT id, restaurant_id, name
        FROM sections
        WHERE restaurant_id = $1
        ORDER BY name
    `
	rows, err := r.db.Query(ctx, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить секции ресторана: %w", err)
	}
	defer rows.Close()

	var sections []*models.Section
	for rows.Next() {
		var section models.Section
		if err := rows.Scan(
			&section.ID,
			&section.RestaurantID,
			&section.Name,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании секции: %w", err)
		}
		sections = append(sections, &section)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по секциям: %w", err)
	}

	return sections, nil
}

func (r *SectionRepository) Update(ctx context.Context, section *models.Section) error {
	query := `
        UPDATE sections
        SET restaurant_id = $1, name = $2
        WHERE id = $3
    `
	commandTag, err := r.db.Exec(ctx, query, section.RestaurantID, section.Name, section.ID)

	if err != nil {
		return fmt.Errorf("не удалось обновить секцию: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("секция с ID %d не найдена", section.ID)
	}

	return nil
}

func (r *SectionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sections WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить секцию: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("секция с ID %d не найдена", id)
	}

	return nil
}
