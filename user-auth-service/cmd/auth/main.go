package main

import (
	"awesomeProject4/user-auth-service/internal/config"
	"awesomeProject4/user-auth-service/internal/db/cache"
	"context"
	"log"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	// Подключение к Redis

	redisClient, err := cache.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
	defer redisClient.Client.Close()

	// Пример: сохраняем токен в Redis после логина пользователя
	token := "example-jwt-token" // предполагаемый JWT токен
	userID := "12345"            // ID пользователя
	ttl := time.Hour * 24        // срок действия сессии - 24 часа

	err = redisClient.SetAuthToken(context.Background(), userID, token, ttl)
	if err != nil {
		log.Fatalf("Failed to set auth token: %v", err)
	}

	// Проверка токена из Redis
	storedToken, err := redisClient.GetAuthToken(context.Background(), userID)
	if err != nil {
		log.Fatalf("Failed to get auth token: %v", err)
	}
	log.Printf("Stored token for user %s: %s", userID, storedToken)
}
