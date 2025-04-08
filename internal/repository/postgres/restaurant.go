package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type RestaurantRepository struct {
	db *pgxpool.Pool
}

func NewRestaurantRepository(db *pgxpool.Pool) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) Create(ctx context.Context, restaurant *models.Restaurant) (int64, error) {
	query := `
        INSERT INTO restaurants (name, city_id, address_ru, address_kz, is_active, _2gis_map)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query,
		restaurant.Name,
		restaurant.CityID,
		restaurant.AddressRU,
		restaurant.AddressKZ,
		restaurant.IsActive,
		restaurant.Map2GIS,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать ресторан: %w", err)
	}

	return id, nil
}

func (r *RestaurantRepository) GetByID(ctx context.Context, id int64) (*models.Restaurant, error) {
	query := `
        SELECT id, name, city_id, address_ru, address_kz, is_active, _2gis_map
        FROM restaurants
        WHERE id = $1
    `
	var restaurant models.Restaurant
	err := r.db.QueryRow(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.CityID,
		&restaurant.AddressRU,
		&restaurant.AddressKZ,
		&restaurant.IsActive,
		&restaurant.Map2GIS,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ресторан с ID %d не найден", id)
		}
		return nil, fmt.Errorf("не удалось получить ресторан: %w", err)
	}

	return &restaurant, nil
}

func (r *RestaurantRepository) GetByCity(ctx context.Context, cityID int64) ([]*models.Restaurant, error) {
	query := `
        SELECT id, name, city_id, address_ru, address_kz, is_active, _2gis_map
        FROM restaurants
        WHERE city_id = $1
        ORDER BY name
    `
	rows, err := r.db.Query(ctx, query, cityID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить рестораны по городу: %w", err)
	}
	defer rows.Close()

	var restaurants []*models.Restaurant
	for rows.Next() {
		var restaurant models.Restaurant
		if err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.CityID,
			&restaurant.AddressRU,
			&restaurant.AddressKZ,
			&restaurant.IsActive,
			&restaurant.Map2GIS,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании ресторана: %w", err)
		}
		restaurants = append(restaurants, &restaurant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по ресторанам: %w", err)
	}

	return restaurants, nil
}

// Update обновляет данные ресторана
func (r *RestaurantRepository) Update(ctx context.Context, restaurant *models.Restaurant) error {
	query := `
        UPDATE restaurants
        SET name = $1, city_id = $2, address_ru = $3, address_kz = $4, is_active = $5, _2gis_map = $6
        WHERE id = $7
    `
	commandTag, err := r.db.Exec(ctx, query,
		restaurant.Name,
		restaurant.CityID,
		restaurant.AddressRU,
		restaurant.AddressKZ,
		restaurant.IsActive,
		restaurant.Map2GIS,
		restaurant.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить ресторан: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("ресторан с ID %d не найден", restaurant.ID)
	}

	return nil
}

func (r *RestaurantRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM restaurants WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить ресторан: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("ресторан с ID %d не найден", id)
	}

	return nil
}

func (r *RestaurantRepository) List(ctx context.Context, active bool) ([]*models.Restaurant, error) {
	var query string
	var args []interface{}

	if active {
		query = `
            SELECT id, name, city_id, address_ru, address_kz, is_active, _2gis_map
            FROM restaurants
            WHERE is_active = true
            ORDER BY name
        `
	} else {
		query = `
            SELECT id, name, city_id, address_ru, address_kz, is_active, _2gis_map
            FROM restaurants
            ORDER BY name
        `
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список ресторанов: %w", err)
	}
	defer rows.Close()

	var restaurants []*models.Restaurant
	for rows.Next() {
		var restaurant models.Restaurant
		if err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.CityID,
			&restaurant.AddressRU,
			&restaurant.AddressKZ,
			&restaurant.IsActive,
			&restaurant.Map2GIS,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании ресторана: %w", err)
		}
		restaurants = append(restaurants, &restaurant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по ресторанам: %w", err)
	}

	return restaurants, nil
}
