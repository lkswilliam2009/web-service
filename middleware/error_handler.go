package middleware

import (
	"github.com/gofiber/fiber/v2"
	appErr "web-service/errors"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	if e, ok := err.(*appErr.AppError); ok {
		return c.Status(e.HTTPCode).JSON(fiber.Map{
			"code":    e.Code,     // ← kirim error code ke client
			"message": e.Message,  // ← message user
		})
	}

	if fe, ok := err.(*fiber.Error); ok {
		return c.Status(fe.Code).JSON(fiber.Map{
			"message": fe.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":    "SYS_001",
		"message": "Internal server error",
	})
}

