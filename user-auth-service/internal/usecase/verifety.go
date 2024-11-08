package usecase

import (
	"awesomeProject4/user-auth-service/internal/domains/models"
	"context"
	"errors"
	"fmt"
)

func (uc *UserUseCase) ConfirmEmail(ctx context.Context, userID int, code string) error {
	cacheKey := fmt.Sprintf("confirm_%d", userID)
	storedCode, err := uc.Cache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}

	if storedCode != code {
		return errors.New("неверный код подтверждения")
	}

	// Обновление статуса пользователя
	if err := uc.UserRepo.UpdateUserStatus(ctx, userID, models.StatusConfirmed); err != nil {
		return err
	}

	// Удаление кода подтверждения из кэша
	if err := uc.Cache.Delete(ctx, cacheKey); err != nil {
		return err
	}

	return nil
}
