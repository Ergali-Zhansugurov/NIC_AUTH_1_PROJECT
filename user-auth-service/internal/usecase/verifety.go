package usecase

import (
	"awesomeProject4/user-auth-service/internal/domains/models"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

func (uc *UserUseCase) ConfirmEmail(ctx context.Context, userID int, code string) error {
	logrus.Println("start ConfirmEmail in usecase")
	cacheKey := fmt.Sprintf("confirm_%d", userID)
	logrus.Println("cod outing frome cache")
	logrus.Println("Cache Key: ", cacheKey)
	storedCode, err := uc.Cache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}
	logrus.Println("geted from radis ConfirmEmail in usecase ")
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
