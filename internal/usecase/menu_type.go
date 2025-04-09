package usecase

import (
	"context"
	"fmt"
	"strings"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type MenuTypeUC struct {
	menuTypeRepo repository.MenuTypeRepository
}

func NewMenuTypeUseCase(menuTypeRepo repository.MenuTypeRepository) *MenuTypeUC {
	return &MenuTypeUC{
		menuTypeRepo: menuTypeRepo,
	}
}

func (uc *MenuTypeUC) Create(ctx context.Context, menuType *models.MenuType) (int64, error) {
	if err := validateMenuType(menuType); err != nil {
		return 0, err
	}

	menuTypes, err := uc.menuTypeRepo.List(ctx)
	if err == nil {
		for _, mt := range menuTypes {
			if strings.EqualFold(mt.Name, menuType.Name) {
				return 0, fmt.Errorf("тип меню с названием '%s' уже существует", menuType.Name)
			}
		}
	}

	return uc.menuTypeRepo.Create(ctx, menuType)
}

func (uc *MenuTypeUC) GetByID(ctx context.Context, id int64) (*models.MenuType, error) {
	menuType, err := uc.menuTypeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить тип меню: %w", err)
	}
	return menuType, nil
}

func (uc *MenuTypeUC) Update(ctx context.Context, menuType *models.MenuType) error {
	if err := validateMenuType(menuType); err != nil {
		return err
	}

	existingMenuType, err := uc.menuTypeRepo.GetByID(ctx, menuType.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти тип меню для обновления: %w", err)
	}

	if existingMenuType.Name != menuType.Name {
		menuTypes, err := uc.menuTypeRepo.List(ctx)
		if err == nil {
			for _, mt := range menuTypes {
				if strings.EqualFold(mt.Name, menuType.Name) && mt.ID != menuType.ID {
					return fmt.Errorf("тип меню с названием '%s' уже существует", menuType.Name)
				}
			}
		}
	}

	return uc.menuTypeRepo.Update(ctx, menuType)
}

func (uc *MenuTypeUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.menuTypeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти тип меню для удаления: %w", err)
	}

	return uc.menuTypeRepo.Delete(ctx, id)
}

func (uc *MenuTypeUC) List(ctx context.Context) ([]*models.MenuType, error) {
	return uc.menuTypeRepo.List(ctx)
}

func validateMenuType(menuType *models.MenuType) error {
	menuType.Name = strings.TrimSpace(menuType.Name)
	if menuType.Name == "" {
		return fmt.Errorf("название типа меню не может быть пустым")
	}

	if len(menuType.Name) < 2 {
		return fmt.Errorf("название типа меню должно быть не менее 2 символов")
	}

	return nil
}
