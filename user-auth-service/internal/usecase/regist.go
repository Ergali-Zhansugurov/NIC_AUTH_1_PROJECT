package usecase

import (
	"awesomeProject4/user-auth-service/internal/auth"
	"awesomeProject4/user-auth-service/internal/domains/interfaces"
	"awesomeProject4/user-auth-service/internal/domains/models"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type UserUseCase struct {
	UserRepo     interfaces.UserRepository
	TokenManager auth.TokenManager
	Cache        interfaces.Cache
	EmailSender  EmailSender // Интерфейс для отправки email
}

// EmailSender определяет интерфейс для отправки писем
type EmailSender interface {
	SendConfirmationEmail(email, code string) error
	SendRecoveryEmail(email, code string) error
}

func NewUserUseCase(userRepo interfaces.UserRepository, tokenManager auth.TokenManager, cache interfaces.Cache, emailSender EmailSender) *UserUseCase {
	return &UserUseCase{
		UserRepo:     userRepo,
		TokenManager: tokenManager,
		Cache:        cache,
		EmailSender:  emailSender,
	}
}

// RegisterUser регистрирует нового пользователя
func (uc *UserUseCase) RegisterUser(ctx context.Context, user *models.User) error {
	// Проверка существования пользователя
	err := uc.UserRepo.FindisUsernewByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	// Создание пользователя
	if err := uc.UserRepo.CreateUser(user); err != nil {
		return err
	}

	// Генерация кода подтверждения
	confirmCode, err := auth.GenerateConfirmationCode()
	if err != nil {
		return err
	}

	// Отправка подтверждающего письма
	if err := uc.EmailSender.SendConfirmationEmail(user.Email, confirmCode); err != nil {
		return err
	}

	// Сохранение кода подтверждения в кэш с TTL 15 минут
	cacheKey := fmt.Sprintf("confirm_%d", user.ID)
	logrus.Println("Cache Key: ", cacheKey)
	if err := uc.Cache.Set(ctx, cacheKey, confirmCode, 15*time.Hour); err != nil {
		return err
	}
	logrus.Println("cod sended to email and cache")
	return nil
}
