package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type CityRepository struct {
	db *pgxpool.Pool
}

func NewCityRepository(db *pgxpool.Pool) *CityRepository {
	return &CityRepository{db: db}
}

func (r *CityRepository) Create(ctx context.Context, city *models.City) (int64, error) {
	query := `
        INSERT INTO cities (name)
        VALUES ($1)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query, city.Name).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать город: %w", err)
	}

	return id, nil
}

func (r *CityRepository) GetByID(ctx context.Context, id int64) (*models.City, error) {
	query := `
        SELECT id, name
        FROM cities
        WHERE id = $1
    `
	var city models.City
	err := r.db.QueryRow(ctx, query, id).Scan(
		&city.ID,
		&city.Name,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("город с ID %d не найден", id)
		}
		return nil, fmt.Errorf("не удалось получить город: %w", err)
	}

	return &city, nil
}

func (r *CityRepository) GetByName(ctx context.Context, name string) (*models.City, error) {
	query := `
        SELECT id, name
        FROM cities
        WHERE name = $1
    `
	var city models.City
	err := r.db.QueryRow(ctx, query, name).Scan(
		&city.ID,
		&city.Name,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("город с названием %s не найден", name)
		}
		return nil, fmt.Errorf("не удалось получить город: %w", err)
	}

	return &city, nil
}

func (r *CityRepository) Update(ctx context.Context, city *models.City) error {
	query := `
        UPDATE cities
        SET name = $1
        WHERE id = $2
    `
	commandTag, err := r.db.Exec(ctx, query, city.Name, city.ID)

	if err != nil {
		return fmt.Errorf("не удалось обновить город: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("город с ID %d не найден", city.ID)
	}

	return nil
}

func (r *CityRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM cities WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить город: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("город с ID %d не найден", id)
	}

	return nil
}

func (r *CityRepository) List(ctx context.Context) ([]*models.City, error) {
	query := `
        SELECT id, name
        FROM cities
        ORDER BY name
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список городов: %w", err)
	}
	defer rows.Close()

	var cities []*models.City
	for rows.Next() {
		var city models.City
		if err := rows.Scan(
			&city.ID,
			&city.Name,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании города: %w", err)
		}
		cities = append(cities, &city)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по городам: %w", err)
	}

	return cities, nil
}
