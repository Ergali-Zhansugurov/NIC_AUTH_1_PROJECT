package main

import (
	"awesomeProject4/user-auth-service/internal/auth"
	"awesomeProject4/user-auth-service/internal/config"
	"awesomeProject4/user-auth-service/internal/db/postgres"
	"awesomeProject4/user-auth-service/internal/repository"
	"awesomeProject4/user-auth-service/internal/server/http"
	"awesomeProject4/user-auth-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	// Подключение к PostgreSQL
	postgresDB, err := postgres.ConnectPostgres(cfg)
	if err != nil {
		logrus.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer postgresDB.DB.Close()

	// Подключение к Redis
	redisClient, err := repository.ConnectRedis(cfg)
	if err != nil {
		logrus.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisClient.Client.Close()

	// Инициализация репозитория пользователей
	userRepo := repository.NewPostgresUserRepository(postgresDB)

	// Инициализация менеджера токенов JWT
	tokenManager, err := auth.NewManager(cfg.JWTSecret)
	if err != nil {
		logrus.Fatalf("Ошибка создания менеджера токенов: %v", err)
	}

	// Инициализация отправителя email
	emailSender := repository.NewSMTPEmailSender(cfg)

	// Инициализация бизнес-логики
	userUseCase := usecase.NewUserUseCase(userRepo, tokenManager, redisClient, emailSender)

	// Инициализация HTTP-обработчиков
	userHandler := http.NewUserHandler(userUseCase)

	// Настройка маршрутов
	router := gin.Default()

	// Маршруты аутентификации
	router.POST("/register", userHandler.Register)
	router.POST("/confirm", userHandler.ConfirmEmail)
	router.POST("/login", userHandler.Login)
	router.POST("/forgot-password", userHandler.ForgotPassword)
	router.POST("/reset-password", userHandler.ResetPassword)

	// Запуск сервера
	if err := router.Run(cfg.ServerPort); err != nil {
		logrus.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
