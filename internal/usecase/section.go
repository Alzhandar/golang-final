package usecase

import (
	"context"
	"fmt"
	"strings"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type SectionUC struct {
	sectionRepo    repository.SectionRepository
	restaurantRepo repository.RestaurantRepository
}

func NewSectionUseCase(sectionRepo repository.SectionRepository, restaurantRepo repository.RestaurantRepository) *SectionUC {
	return &SectionUC{
		sectionRepo:    sectionRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (uc *SectionUC) Create(ctx context.Context, section *models.Section) (int64, error) {
	if err := validateSection(section); err != nil {
		return 0, err
	}

	_, err := uc.restaurantRepo.GetByID(ctx, section.RestaurantID)
	if err != nil {
		return 0, fmt.Errorf("ресторан не найден: %w", err)
	}

	sections, err := uc.sectionRepo.GetByRestaurant(ctx, section.RestaurantID)
	if err == nil {
		for _, s := range sections {
			if strings.EqualFold(s.Name, section.Name) {
				return 0, fmt.Errorf("секция с названием '%s' уже существует в этом ресторане", section.Name)
			}
		}
	}

	return uc.sectionRepo.Create(ctx, section)
}

func (uc *SectionUC) GetByID(ctx context.Context, id int64) (*models.Section, error) {
	section, err := uc.sectionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить секцию: %w", err)
	}
	return section, nil
}

func (uc *SectionUC) GetByRestaurant(ctx context.Context, restaurantID int64) ([]*models.Section, error) {
	_, err := uc.restaurantRepo.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("указанный ресторан не существует: %w", err)
	}

	return uc.sectionRepo.GetByRestaurant(ctx, restaurantID)
}

func (uc *SectionUC) Update(ctx context.Context, section *models.Section) error {
	if err := validateSection(section); err != nil {
		return err
	}

	existingSection, err := uc.sectionRepo.GetByID(ctx, section.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти секцию для обновления: %w", err)
	}

	_, err = uc.restaurantRepo.GetByID(ctx, section.RestaurantID)
	if err != nil {
		return fmt.Errorf("указанный ресторан не существует: %w", err)
	}

	if existingSection.RestaurantID == section.RestaurantID && existingSection.Name != section.Name ||
		existingSection.RestaurantID != section.RestaurantID {
		sections, err := uc.sectionRepo.GetByRestaurant(ctx, section.RestaurantID)
		if err == nil {
			for _, s := range sections {
				if strings.EqualFold(s.Name, section.Name) && s.ID != section.ID {
					return fmt.Errorf("секция с названием '%s' уже существует в этом ресторане", section.Name)
				}
			}
		}
	}

	return uc.sectionRepo.Update(ctx, section)
}

func (uc *SectionUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.sectionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти секцию для удаления: %w", err)
	}

	return uc.sectionRepo.Delete(ctx, id)
}

func validateSection(section *models.Section) error {
	section.Name = strings.TrimSpace(section.Name)

	if section.Name == "" {
		return fmt.Errorf("название секции не может быть пустым")
	}

	if section.RestaurantID <= 0 {
		return fmt.Errorf("некорректный ID ресторана")
	}

	return nil
}
