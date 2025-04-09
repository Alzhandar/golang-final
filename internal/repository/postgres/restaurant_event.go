package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type RestaurantEventRepository struct {
	db *pgxpool.Pool
}

func NewRestaurantEventRepository(db *pgxpool.Pool) *RestaurantEventRepository {
	return &RestaurantEventRepository{db: db}
}

func (r *RestaurantEventRepository) Create(ctx context.Context, event *models.RestaurantEvent) (int64, error) {
	query := `
        INSERT INTO restaurant_events (name, eventtype, "desc", price, img)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query,
		event.Name,
		event.EventType,
		event.Description,
		event.Price,
		event.Img,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать событие ресторана: %w", err)
	}

	return id, nil
}

func (r *RestaurantEventRepository) GetByID(ctx context.Context, id int64) (*models.RestaurantEvent, error) {
	query := `
        SELECT id, name, eventtype, "desc", price, img
        FROM restaurant_events
        WHERE id = $1
    `
	var event models.RestaurantEvent
	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Name,
		&event.EventType,
		&event.Description,
		&event.Price,
		&event.Img,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("событие ресторана с ID %d не найдено", id)
		}
		return nil, fmt.Errorf("не удалось получить событие ресторана: %w", err)
	}

	return &event, nil
}

func (r *RestaurantEventRepository) GetByType(ctx context.Context, eventType models.EventType) ([]*models.RestaurantEvent, error) {
	query := `
        SELECT id, name, eventtype, "desc", price, img
        FROM restaurant_events
        WHERE eventtype = $1
        ORDER BY name
    `
	rows, err := r.db.Query(ctx, query, eventType)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить события ресторана по типу: %w", err)
	}
	defer rows.Close()

	var events []*models.RestaurantEvent
	for rows.Next() {
		var event models.RestaurantEvent
		if err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.EventType,
			&event.Description,
			&event.Price,
			&event.Img,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании события ресторана: %w", err)
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по событиям ресторана: %w", err)
	}

	return events, nil
}

func (r *RestaurantEventRepository) Update(ctx context.Context, event *models.RestaurantEvent) error {
	query := `
        UPDATE restaurant_events
        SET name = $1, eventtype = $2, "desc" = $3, price = $4, img = $5
        WHERE id = $6
    `
	commandTag, err := r.db.Exec(ctx, query,
		event.Name,
		event.EventType,
		event.Description,
		event.Price,
		event.Img,
		event.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить событие ресторана: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("событие ресторана с ID %d не найдено", event.ID)
	}

	return nil
}

func (r *RestaurantEventRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM restaurant_events WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить событие ресторана: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("событие ресторана с ID %d не найдено", id)
	}

	return nil
}

func (r *RestaurantEventRepository) List(ctx context.Context) ([]*models.RestaurantEvent, error) {
	query := `
        SELECT id, name, eventtype, "desc", price, img
        FROM restaurant_events
        ORDER BY name
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список событий ресторана: %w", err)
	}
	defer rows.Close()

	var events []*models.RestaurantEvent
	for rows.Next() {
		var event models.RestaurantEvent
		if err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.EventType,
			&event.Description,
			&event.Price,
			&event.Img,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании события ресторана: %w", err)
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по событиям ресторана: %w", err)
	}

	return events, nil
}
