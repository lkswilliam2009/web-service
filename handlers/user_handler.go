package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"encoding/json"

	"web-service/config"
	"web-service/models"
)

func Me(c *fiber.Ctx) error {
	// ===== JWT =====
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.ErrUnauthorized
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return fiber.ErrUnauthorized
	}

	role, _ := claims["role"].(string)
	branch, _ := claims["branch"].(string)

	// parse groups dari JWT
	var groups []models.Group
	if raw, ok := claims["groups"]; ok {
		bytes, _ := json.Marshal(raw)
		_ = json.Unmarshal(bytes, &groups)
	}

	// ===== FETCH USER BASIC INFO =====
	var (
		username string
		email string
		role_name string
		branch_name string
		lastLogin string
		lastOS string
		lastIP string
	)

	err := config.DB.QueryRow(`
		SELECT uname, email, role_name, branch_name, last_login, last_os, last_ip
		FROM tv_user
		WHERE userid = $1
	`, userID).Scan(
		&username,
		&email,
		&role_name,
		&branch_name,
		&lastLogin,
		&lastOS,
		&lastIP,
	)

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

	// ===== RESPONSE =====
	return c.JSON(fiber.Map{
		"id": userID,
		"username": username,
		"email":email,
		"last_login": lastLogin,
		"last_os": lastOS,
		"last_ip": lastIP,
		"role": role,
		"role_name": role_name,
		"branch": branch,
		"branch_name": branch_name,
		"groups": groups,
		"device": device,
	})
}