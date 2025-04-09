package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"restaurant-management/internal/models"
)

type TableRepository struct {
	db *pgxpool.Pool
}

func NewTableRepository(db *pgxpool.Pool) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Create(ctx context.Context, table *models.Table) (int64, error) {
	query := `
        INSERT INTO tables (number_of_table, section_id, qr)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int64
	err := r.db.QueryRow(ctx, query, table.NumberOfTable, table.SectionID, table.QR).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать столик: %w", err)
	}

	return id, nil
}

func (r *TableRepository) GetByID(ctx context.Context, id int64) (*models.Table, error) {
	query := `
        SELECT id, number_of_table, section_id, qr
        FROM tables
        WHERE id = $1
    `
	var table models.Table
	err := r.db.QueryRow(ctx, query, id).Scan(
		&table.ID,
		&table.NumberOfTable,
		&table.SectionID,
		&table.QR,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("столик с ID %d не найден", id)
		}
		return nil, fmt.Errorf("не удалось получить столик: %w", err)
	}

	return &table, nil
}

func (r *TableRepository) GetBySection(ctx context.Context, sectionID int64) ([]*models.Table, error) {
	query := `
        SELECT id, number_of_table, section_id, qr
        FROM tables
        WHERE section_id = $1
        ORDER BY number_of_table
    `
	rows, err := r.db.Query(ctx, query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить столики секции: %w", err)
	}
	defer rows.Close()

	var tables []*models.Table
	for rows.Next() {
		var table models.Table
		if err := rows.Scan(
			&table.ID,
			&table.NumberOfTable,
			&table.SectionID,
			&table.QR,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании столика: %w", err)
		}
		tables = append(tables, &table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по столикам: %w", err)
	}

	return tables, nil
}

func (r *TableRepository) Update(ctx context.Context, table *models.Table) error {
	query := `
        UPDATE tables
        SET number_of_table = $1, section_id = $2, qr = $3
        WHERE id = $4
    `
	commandTag, err := r.db.Exec(ctx, query, table.NumberOfTable, table.SectionID, table.QR, table.ID)

	if err != nil {
		return fmt.Errorf("не удалось обновить столик: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("столик с ID %d не найден", table.ID)
	}

	return nil
}

func (r *TableRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tables WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("не удалось удалить столик: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("столик с ID %d не найден", id)
	}

	return nil
}

func (r *TableRepository) GenerateQR(ctx context.Context, tableID int64) (string, error) {
	qrCode := fmt.Sprintf("table-%d", tableID)

	query := `UPDATE tables SET qr = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, qrCode, tableID)
	if err != nil {
		return "", fmt.Errorf("не удалось обновить QR-код столика: %w", err)
	}

	return qrCode, nil
}
