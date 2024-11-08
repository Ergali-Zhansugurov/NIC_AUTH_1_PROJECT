package auth

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

// Config для базовой настройки JWT-аутентификации
type Config struct {
	SigningKey     []byte             // Секретный ключ для подписи JWT
	SuccessHandler fiber.Handler      // Обработчик успешной авторизации
	ErrorHandler   fiber.ErrorHandler // Обработчик ошибки авторизации
}

// New создает middleware для авторизации на основе JWT
func New(config Config) fiber.Handler {
	if config.SuccessHandler == nil {
		config.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
		}
	}

	// Функция для проверки и извлечения JWT-токена
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return config.ErrorHandler(c, errors.New("Missing or malformed JWT"))
		}

		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.SigningKey, nil
		})

		if err != nil || !token.Valid {
			return config.ErrorHandler(c, err)
		}

		// Добавляем информацию из токена в контекст
		c.Locals("user", token)
		return config.SuccessHandler(c)
	}
}
