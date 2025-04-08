package usecase

import (
	"context"
	"fmt"
	"strings"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type CityUC struct {
	cityRepo repository.CityRepository
}

func NewCityUseCase(cityRepo repository.CityRepository) *CityUC {
	return &CityUC{
		cityRepo: cityRepo,
	}
}

func (uc *CityUC) Create(ctx context.Context, city *models.City) (int64, error) {
	if err := validateCity(city); err != nil {
		return 0, err
	}

	city.Name = strings.TrimSpace(city.Name)
	existingCity, err := uc.cityRepo.GetByName(ctx, city.Name)
	if err == nil && existingCity != nil {
		return 0, fmt.Errorf("город с названием '%s' уже существует", city.Name)
	}

	return uc.cityRepo.Create(ctx, city)
}

func (uc *CityUC) GetByID(ctx context.Context, id int64) (*models.City, error) {
	city, err := uc.cityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить город: %w", err)
	}
	return city, nil
}

func (uc *CityUC) Update(ctx context.Context, city *models.City) error {
	if err := validateCity(city); err != nil {
		return err
	}

	existingCity, err := uc.cityRepo.GetByID(ctx, city.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти город для обновления: %w", err)
	}

	city.Name = strings.TrimSpace(city.Name)
	if existingCity.Name != city.Name {
		dupCity, err := uc.cityRepo.GetByName(ctx, city.Name)
		if err == nil && dupCity != nil && dupCity.ID != city.ID {
			return fmt.Errorf("город с названием '%s' уже существует", city.Name)
		}
	}

	return uc.cityRepo.Update(ctx, city)
}

func (uc *CityUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.cityRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти город для удаления: %w", err)
	}

	return uc.cityRepo.Delete(ctx, id)
}

func (uc *CityUC) List(ctx context.Context) ([]*models.City, error) {
	return uc.cityRepo.List(ctx)
}

func validateCity(city *models.City) error {
	city.Name = strings.TrimSpace(city.Name)
	if city.Name == "" {
		return fmt.Errorf("название города не может быть пустым")
	}

	if len(city.Name) < 2 {
		return fmt.Errorf("название города должно быть не менее 2 символов")
	}

	return nil
}
