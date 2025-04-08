package usecase

import (
	"context"
	"fmt"

	"restaurant-management/internal/models"
	"restaurant-management/internal/repository"
)

type UserUC struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUC {
	return &UserUC{
		userRepo: userRepo,
	}
}

func (uc *UserUC) Create(ctx context.Context, user *models.User) (int64, error) {
	existingUser, err := uc.userRepo.GetByPhone(ctx, user.PhoneNumber)
	if err == nil && existingUser != nil {
		return 0, fmt.Errorf("пользователь с номером телефона %s уже существует", user.PhoneNumber)
	}

	if err := validateUser(user); err != nil {
		return 0, err
	}

	if user.Language == "" {
		user.Language = "ru"
	}

	if !user.IsActive {
		user.IsActive = true
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUC) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	return user, nil
}

func (uc *UserUC) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	user, err := uc.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	return user, nil
}

func (uc *UserUC) Update(ctx context.Context, user *models.User) error {
	existingUser, err := uc.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя для обновления: %w", err)
	}

	if existingUser.PhoneNumber != user.PhoneNumber {
		dupUser, err := uc.userRepo.GetByPhone(ctx, user.PhoneNumber)
		if err == nil && dupUser != nil && dupUser.ID != user.ID {
			return fmt.Errorf("пользователь с номером телефона %s уже существует", user.PhoneNumber)
		}
	}

	if err := validateUser(user); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUC) Delete(ctx context.Context, id int64) error {
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя для удаления: %w", err)
	}

	return uc.userRepo.Delete(ctx, id)
}

func (uc *UserUC) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return uc.userRepo.List(ctx, limit, offset)
}

func validateUser(user *models.User) error {
	if user.PhoneNumber == "" {
		return fmt.Errorf("номер телефона не может быть пустым")
	}

	if len(user.PhoneNumber) < 10 {
		return fmt.Errorf("некорректный номер телефона")
	}

	if user.Name == "" {
		return fmt.Errorf("имя не может быть пустым")
	}

	return nil
}
