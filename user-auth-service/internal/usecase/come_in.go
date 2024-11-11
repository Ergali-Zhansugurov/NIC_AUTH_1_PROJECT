package usecase

import (
	"awesomeProject4/user-auth-service/internal/domains/models"
	"context"
	"errors"
	"fmt"
	"time"
)

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	// Проверка пароля
	if err := user.CheckPassword(password); err != nil {
		return "", errors.New("неверные учетные данные")
	}
	// Проверка статуса подтверждения
	if user.Status != models.StatusConfirmed {
		return "", errors.New("email пользователя не подтвержден")
	}

	// Генерация JWT-токена
	token, err := uc.TokenManager.NewJWT(fmt.Sprintf("%d", user.ID), 24*time.Hour)
	if err != nil {
		return "", err
	}

	// Сохранение токена в кэше с TTL 24 часа
	cacheKey := fmt.Sprintf("token_%d", user.ID)
	if err := uc.Cache.Set(ctx, cacheKey, token, 24*time.Hour); err != nil {
		return "", err
	}

	return token, nil
}
