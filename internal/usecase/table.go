package usecase

import (
	"context"
	"fmt"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type TableUC struct {
	tableRepo   repository.TableRepository
	sectionRepo repository.SectionRepository
}

func NewTableUseCase(tableRepo repository.TableRepository, sectionRepo repository.SectionRepository) *TableUC {
	return &TableUC{
		tableRepo:   tableRepo,
		sectionRepo: sectionRepo,
	}
}

func (uc *TableUC) Create(ctx context.Context, table *models.Table) (int64, error) {
	_, err := uc.sectionRepo.GetByID(ctx, table.SectionID)
	if err != nil {
		return 0, fmt.Errorf("указанная секция не существует: %w", err)
	}

	if err := validateTable(table); err != nil {
		return 0, err
	}

	tables, err := uc.tableRepo.GetBySection(ctx, table.SectionID)
	if err == nil {
		for _, t := range tables {
			if t.NumberOfTable == table.NumberOfTable {
				return 0, fmt.Errorf("столик с номером %d уже существует в этой секции", table.NumberOfTable)
			}
		}
	}

	id, err := uc.tableRepo.Create(ctx, table)
	if err != nil {
		return 0, err
	}

	if table.QR == "" {
		qr, err := uc.tableRepo.GenerateQR(ctx, id)
		if err != nil {
			return id, nil
		}
		table.QR = qr

		table.ID = id
		_ = uc.tableRepo.Update(ctx, table)
	}

	return id, nil
}

func (uc *TableUC) GetByID(ctx context.Context, id int64) (*models.Table, error) {
	table, err := uc.tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить столик: %w", err)
	}
	return table, nil
}

func (uc *TableUC) GetBySection(ctx context.Context, sectionID int64) ([]*models.Table, error) {
	_, err := uc.sectionRepo.GetByID(ctx, sectionID)
	if err != nil {
		return nil, fmt.Errorf("указанная секция не существует: %w", err)
	}

	return uc.tableRepo.GetBySection(ctx, sectionID)
}

func (uc *TableUC) Update(ctx context.Context, table *models.Table) error {
	existingTable, err := uc.tableRepo.GetByID(ctx, table.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти столик для обновления: %w", err)
	}

	_, err = uc.sectionRepo.GetByID(ctx, table.SectionID)
	if err != nil {
		return fmt.Errorf("указанная секция не существует: %w", err)
	}

	if err := validateTable(table); err != nil {
		return err
	}

	if existingTable.SectionID == table.SectionID && existingTable.NumberOfTable != table.NumberOfTable ||
		existingTable.SectionID != table.SectionID {
		tables, err := uc.tableRepo.GetBySection(ctx, table.SectionID)
		if err == nil {
			for _, t := range tables {
				if t.NumberOfTable == table.NumberOfTable && t.ID != table.ID {
					return fmt.Errorf("столик с номером %d уже существует в этой секции", table.NumberOfTable)
				}
			}
		}
	}

	return uc.tableRepo.Update(ctx, table)
}

func (uc *TableUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.tableRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти столик для удаления: %w", err)
	}

	return uc.tableRepo.Delete(ctx, id)
}

func (uc *TableUC) GenerateQR(ctx context.Context, tableID int64) (string, error) {
	_, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return "", fmt.Errorf("не удалось найти столик для генерации QR-кода: %w", err)
	}

	return uc.tableRepo.GenerateQR(ctx, tableID)
}

func validateTable(table *models.Table) error {
	if table.NumberOfTable <= 0 {
		return fmt.Errorf("номер столика должен быть положительным числом")
	}

	if table.SectionID <= 0 {
		return fmt.Errorf("необходимо указать корректный ID секции")
	}

	return nil
}
