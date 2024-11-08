package usecase

import (
	"awesomeProject4/user-auth-service/internal/auth"
	"context"
	"errors"
	"fmt"
	"time"
)

// ForgotPassword инициирует процесс восстановления пароля
func (uc *UserUseCase) ForgotPassword(ctx context.Context, email string) error {
	user, err := uc.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Генерация кода восстановления
	recoveryCode, err := auth.GenerateConfirmationCode()
	if err != nil {
		return err
	}

	// Отправка письма восстановления
	if err := uc.EmailSender.SendRecoveryEmail(user.Email, recoveryCode); err != nil {
		return err
	}

	// Сохранение кода восстановления в кэш с TTL 10 минут
	cacheKey := fmt.Sprintf("recovery_%d", user.ID)
	if err := uc.Cache.Set(ctx, cacheKey, recoveryCode, 10*time.Minute); err != nil {
		return err
	}

	return nil
}

// ResetPassword сбрасывает пароль пользователя
func (uc *UserUseCase) ResetPassword(ctx context.Context, userID int, code, newPassword string) error {
	cacheKey := fmt.Sprintf("recovery_%d", userID)
	storedCode, err := uc.Cache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}

	if storedCode != code {
		return errors.New("неверный код восстановления")
	}

	// Хэширование нового пароля
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Обновление пароля в базе данных
	if err := uc.UserRepo.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return err
	}

	// Удаление кода восстановления из кэша
	if err := uc.Cache.Delete(ctx, cacheKey); err != nil {
		return err
	}

	return nil
}
