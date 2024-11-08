package cache

import (
	"awesomeProject4/user-auth-service/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

// ConnectRedis создает новое подключение к Redis
func ConnectRedis(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Проверим подключение
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{Client: client}, nil
}
func (r *RedisClient) SetAuthToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	return r.Client.Set(ctx, userID, token, ttl).Err()
}
func (r *RedisClient) GetAuthToken(ctx context.Context, userID string) (string, error) {
	return r.Client.Get(ctx, userID).Result()
}
func (r *RedisClient) DeleteAuthToken(ctx context.Context, userID string) error {
	return r.Client.Del(ctx, userID).Err()
}
