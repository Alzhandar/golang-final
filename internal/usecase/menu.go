package usecase

import (
	"context"
	"fmt"
	"strings"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type MenuUC struct {
	menuRepo       repository.MenuRepository
	restaurantRepo repository.RestaurantRepository
}

func NewMenuUseCase(menuRepo repository.MenuRepository, restaurantRepo repository.RestaurantRepository) *MenuUC {
	return &MenuUC{
		menuRepo:       menuRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (uc *MenuUC) Create(ctx context.Context, menu *models.Menu) (int64, error) {
	if err := validateMenu(menu); err != nil {
		return 0, err
	}

	_, err := uc.restaurantRepo.GetByID(ctx, menu.RestaurantID)
	if err != nil {
		return 0, fmt.Errorf("указанный ресторан не существует: %w", err)
	}

	return uc.menuRepo.Create(ctx, menu)
}

func (uc *MenuUC) GetByID(ctx context.Context, id int64) (*models.Menu, error) {
	menu, err := uc.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить меню: %w", err)
	}
	return menu, nil
}

func (uc *MenuUC) GetByRestaurant(ctx context.Context, restaurantID int64) ([]*models.Menu, error) {
	_, err := uc.restaurantRepo.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("указанный ресторан не существует: %w", err)
	}

	return uc.menuRepo.GetByRestaurant(ctx, restaurantID)
}

func (uc *MenuUC) Update(ctx context.Context, menu *models.Menu) error {
	if err := validateMenu(menu); err != nil {
		return err
	}

	_, err := uc.menuRepo.GetByID(ctx, menu.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти меню для обновления: %w", err)
	}

	_, err = uc.restaurantRepo.GetByID(ctx, menu.RestaurantID)
	if err != nil {
		return fmt.Errorf("указанный ресторан не существует: %w", err)
	}

	return uc.menuRepo.Update(ctx, menu)
}

func (uc *MenuUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.menuRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти меню для удаления: %w", err)
	}

	return uc.menuRepo.Delete(ctx, id)
}

func validateMenu(menu *models.Menu) error {
	menu.NameRU = strings.TrimSpace(menu.NameRU)
	if menu.NameRU == "" {
		return fmt.Errorf("название меню на русском не может быть пустым")
	}

	menu.NameKZ = strings.TrimSpace(menu.NameKZ)

	if menu.RestaurantID <= 0 {
		return fmt.Errorf("необходимо указать корректный ID ресторана")
	}

	return nil
}
