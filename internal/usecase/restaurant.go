package usecase

import (
	"context"
	"fmt"
	"strings"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type RestaurantUC struct {
	restaurantRepo repository.RestaurantRepository
	cityRepo       repository.CityRepository
}

func NewRestaurantUseCase(restaurantRepo repository.RestaurantRepository, cityRepo repository.CityRepository) *RestaurantUC {
	return &RestaurantUC{
		restaurantRepo: restaurantRepo,
		cityRepo:       cityRepo,
	}
}

func (uc *RestaurantUC) Create(ctx context.Context, restaurant *models.Restaurant) (int64, error) {
	if err := validateRestaurant(restaurant); err != nil {
		return 0, err
	}

	_, err := uc.cityRepo.GetByID(ctx, restaurant.CityID)
	if err != nil {
		return 0, fmt.Errorf("указанный город не существует: %w", err)
	}

	return uc.restaurantRepo.Create(ctx, restaurant)
}

func (uc *RestaurantUC) GetByID(ctx context.Context, id int64) (*models.Restaurant, error) {
	restaurant, err := uc.restaurantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить ресторан: %w", err)
	}
	return restaurant, nil
}

func (uc *RestaurantUC) GetByCity(ctx context.Context, cityID int64) ([]*models.Restaurant, error) {
	_, err := uc.cityRepo.GetByID(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("указанный город не существует: %w", err)
	}

	return uc.restaurantRepo.GetByCity(ctx, cityID)
}

func (uc *RestaurantUC) Update(ctx context.Context, restaurant *models.Restaurant) error {
	if err := validateRestaurant(restaurant); err != nil {
		return err
	}

	_, err := uc.restaurantRepo.GetByID(ctx, restaurant.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти ресторан для обновления: %w", err)
	}

	_, err = uc.cityRepo.GetByID(ctx, restaurant.CityID)
	if err != nil {
		return fmt.Errorf("указанный город не существует: %w", err)
	}

	return uc.restaurantRepo.Update(ctx, restaurant)
}

func (uc *RestaurantUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.restaurantRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти ресторан для удаления: %w", err)
	}

	return uc.restaurantRepo.Delete(ctx, id)
}

func (uc *RestaurantUC) List(ctx context.Context, active bool) ([]*models.Restaurant, error) {
	return uc.restaurantRepo.List(ctx, active)
}

func validateRestaurant(restaurant *models.Restaurant) error {
	restaurant.Name = strings.TrimSpace(restaurant.Name)
	if restaurant.Name == "" {
		return fmt.Errorf("название ресторана не может быть пустым")
	}

	if len(restaurant.Name) < 3 {
		return fmt.Errorf("название ресторана должно быть не менее 3 символов")
	}

	restaurant.AddressRU = strings.TrimSpace(restaurant.AddressRU)
	if restaurant.AddressRU == "" {
		return fmt.Errorf("адрес ресторана (RU) не может быть пустым")
	}

	if restaurant.CityID <= 0 {
		return fmt.Errorf("необходимо указать корректный ID города")
	}

	return nil
}
