package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JWTProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte("ebba34271e82839008c40f03c287ffcc"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}
