package usecase

import (
	"context"
	"fmt"
	"strings"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type RestaurantEventUC struct {
	eventRepo repository.RestaurantEventRepository
}

func NewRestaurantEventUseCase(eventRepo repository.RestaurantEventRepository) *RestaurantEventUC {
	return &RestaurantEventUC{
		eventRepo: eventRepo,
	}
}

func (uc *RestaurantEventUC) Create(ctx context.Context, event *models.RestaurantEvent) (int64, error) {
	if err := validateRestaurantEvent(event); err != nil {
		return 0, err
	}

	return uc.eventRepo.Create(ctx, event)
}

func (uc *RestaurantEventUC) GetByID(ctx context.Context, id int64) (*models.RestaurantEvent, error) {
	event, err := uc.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить событие ресторана: %w", err)
	}
	return event, nil
}

func (uc *RestaurantEventUC) GetByType(ctx context.Context, eventType models.EventType) ([]*models.RestaurantEvent, error) {
	switch eventType {
	case models.EventTypeWedding, models.EventTypeBirthday, models.EventTypeCorporate:
	default:
		return nil, fmt.Errorf("неизвестный тип события: %s", eventType)
	}

	return uc.eventRepo.GetByType(ctx, eventType)
}

func (uc *RestaurantEventUC) Update(ctx context.Context, event *models.RestaurantEvent) error {
	if err := validateRestaurantEvent(event); err != nil {
		return err
	}

	_, err := uc.eventRepo.GetByID(ctx, event.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти событие ресторана для обновления: %w", err)
	}

	return uc.eventRepo.Update(ctx, event)
}

func (uc *RestaurantEventUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.eventRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти событие ресторана для удаления: %w", err)
	}

	return uc.eventRepo.Delete(ctx, id)
}

func (uc *RestaurantEventUC) List(ctx context.Context) ([]*models.RestaurantEvent, error) {
	return uc.eventRepo.List(ctx)
}

func validateRestaurantEvent(event *models.RestaurantEvent) error {
	event.Name = strings.TrimSpace(event.Name)
	if event.Name == "" {
		return fmt.Errorf("название события не может быть пустым")
	}

	switch event.EventType {
	case models.EventTypeWedding, models.EventTypeBirthday, models.EventTypeCorporate:
	default:
		return fmt.Errorf("неизвестный тип события: %s", event.EventType)
	}

	if event.Price < 0 {
		return fmt.Errorf("цена не может быть отрицательной")
	}

	return nil
}
