package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"web-service/config"
)

func Me(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return fiber.ErrUnauthorized
	}

	// ===== FETCH USER FROM DB =====
	var (
		username string
		email string
		last_login string
		last_os string
		last_ip string
		role_name string
	)

	err := config.DB.QueryRow(`
		SELECT uname, email, last_login, last_os, last_ip, role_name
		FROM tv_user
		WHERE userid = $1
	`, userID).Scan(&username, &email, &last_login, &last_os, &last_ip, &role_name)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.ErrUnauthorized
		}
		return fiber.ErrInternalServerError
	}

	// ===== DEVICE INFO =====
	device := fiber.Map{
		"user_agent": c.Get("User-Agent"),
		"ip":         c.IP(),
		"platform":   c.Get("Sec-CH-UA-Platform"),
	}

	return c.JSON(fiber.Map{
		"id": userID,
		"username": username,
		"email": email,
		"last_login": last_login,
		"last_os": last_os,
		"last_ip": last_ip,
		"role": role_name,
		"device": device,
	})
}