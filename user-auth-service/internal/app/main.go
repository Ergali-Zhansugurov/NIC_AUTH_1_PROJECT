package app

import (
	"awesomeProject4/user-auth-service/internal/auth"
	"awesomeProject4/user-auth-service/internal/config"
	"awesomeProject4/user-auth-service/internal/db/postgres"
	"awesomeProject4/user-auth-service/internal/repository"
	"awesomeProject4/user-auth-service/internal/server/http"
	"awesomeProject4/user-auth-service/internal/usecase"
	"fmt"
	"github.com/sirupsen/logrus"
)

func Run() {
	cfg := config.LoadConfig()

	// Подключение к PostgreSQL
	postgresDB, err := postgres.ConnectPostgres(cfg)
	if err != nil {
		logrus.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer postgresDB.DB.Close()

	fmt.Println("postgres connected")

	// Подключение к Redis
	redisClient, err := repository.ConnectRedis(cfg)
	if err != nil {
		logrus.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisClient.Client.Close()

	fmt.Println("redis connected")

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
	srv := http.NewServer(cfg)

	// Маршруты аутентификации
	srv.Router.POST("/register", userHandler.Register)
	srv.Router.POST("/confirm", userHandler.ConfirmEmail)
	srv.Router.POST("/login", userHandler.Login)
	srv.Router.POST("/forgot-password", userHandler.ForgotPassword)
	srv.Router.POST("/reset-password", userHandler.ResetPassword)

	// Запуск сервера
	if err := srv.Run(); err != nil {
		logrus.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
