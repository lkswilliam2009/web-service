package handlers

import (
	"time"
	_"fmt"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"

	"web-service/config"
	"web-service/utils"
)

func Register(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Password hashing failed",
		})
	}

	if body.Username == "" || body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username, Email fields, Password are required",
		})
	}

	_, err = config.DB.Exec(
		"INSERT INTO users(uname,email,password) VALUES($1,$2,$3)",
		body.Username, body.Email, hash,
	)

	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_uname_key":
				return c.Status(409).JSON(fiber.Map{
					"error": "Username already exists",
				})
			case "users_email_key":
				return c.Status(409).JSON(fiber.Map{
					"error": "Email already exists",
				})
			}
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Register success",
		"data": fiber.Map{
			"username": body.Username,
			"email":    body.Email,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if body.Username == "" && body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username or Email fields and Password are required",
		})
	}

	var (
		id    string
		hash  string
		uname string
		roleid sql.NullString
		err error
	)
	if body.Email != "" {
		err = config.DB.QueryRow(
			"SELECT userid,password,uname,roleid FROM users WHERE email=$1",
			body.Email,
		).Scan(&id, &hash, &uname, &roleid)
	} else {
		err = config.DB.QueryRow(
			"SELECT userid,password,uname,roleid FROM users WHERE uname=$1",
			body.Username,
		).Scan(&id, &hash, &uname, &roleid)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid credentials, user not found!",
			})
		}

		// error DB lain
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
			"details": err.Error(),
		})
	}

	if !roleid.Valid {
		return c.Status(403).JSON(fiber.Map{
			"error": "User has no role assigned, contact your Administrator to assign!",
		})
	}

	if err := utils.CheckPassword(hash, body.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	access, err := utils.GenerateAccessToken(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Access token generation failed",
		})
	}

	refresh, err := utils.GenerateRefreshToken(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Refresh token generation failed",
		})
	}

	_, err = config.DB.Exec(
		"UPDATE users SET refresh_token=$1, last_login=$2 WHERE userid=$3",
		refresh,time.Now().UTC(), id,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Token save failed",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
		"data": fiber.Map{
			"username": uname,
		},
	})
}

func Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Refresh token required",
		})
	}

	// Parse & validate refresh token
	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return config.RefreshSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Cek refresh token di DB
	var savedToken string
	err = config.DB.QueryRow(
		"SELECT refresh_token FROM users WHERE userid=$1",
		userID,
	).Scan(&savedToken)

	if err != nil || savedToken != body.RefreshToken {
		return c.Status(401).JSON(fiber.Map{
			"error": "Refresh token revoked",
		})
	}

	access, err := utils.GenerateAccessToken(userID)
	newRefresh, _ := utils.GenerateRefreshToken(userID)

	config.DB.Exec(
		"UPDATE users SET refresh_token=$1 WHERE userid=$2",
		newRefresh, userID,
	)

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": newRefresh,
	})
}

func Logout(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid user id",
		})
	}

	_, err := config.DB.Exec(
		"UPDATE users SET refresh_token=NULL WHERE userid=$1",
		userID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Logout failed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logout success",
	})
}


