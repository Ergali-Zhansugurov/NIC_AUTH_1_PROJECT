package auth

import (
	"github.com/gofiber/fiber/v2"

	"os"
)

func CheckAuthToken() fiber.Handler {
	return New(Config{
		SigningKey: []byte(os.Getenv("jwt_secret")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		},
	})
}
