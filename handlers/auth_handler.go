package handlers

import (
	"time"
	"fmt"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	_"github.com/lib/pq"
	"github.com/mssola/useragent"

	"web-service/config"
	"web-service/utils"

	appErr "web-service/errors"
)

func Register(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return appErr.InvalidJSON(err)
	}

	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		return appErr.Internal(err)
	}

	if body.Username == "" {
		return appErr.BadRequest("Username required")
	}

	if body.Email == "" {
		return appErr.BadRequest("Email required")
	}

	if body.Password == "" {
		return appErr.BadRequest("Password required")
	}

	_, err = config.DB.Exec(
		"INSERT INTO users(uname,email,password) VALUES($1,$2,$3)",
		body.Username, body.Email, hash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return appErr.Unauthorized("Invalid credentials")
		}
		return appErr.FromDB(err)
	}

	return c.JSON(fiber.Map{
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
		return appErr.InvalidJSON(err)
	}

	if body.Username == "" && body.Email == "" {
		return appErr.BadRequest("Username or Email is required")
	}

	if body.Password == "" {
		return appErr.BadRequest("Password is required")
	}

	var (
		id     string
		hash   string
		uname  string
		roleID sql.NullString
		Last_login sql.NullString
		Last_os sql.NullString
		Last_browser sql.NullString
		Last_ip sql.NullString
		Role_name sql.NullString
		Role_description sql.NullString
	)

	var err error
	if body.Email != "" {
		err = config.DB.QueryRow(
			`SELECT userid, password, uname, roleid, last_login, last_os, last_browser, last_ip, role_name, role_description
			 FROM tv_user
			 WHERE email = $1`,
			body.Email,
		).Scan(&id, &hash, &uname, &roleID, &Last_login, &Last_os, &Last_browser, &Last_ip, &Role_name, &Role_description)
	} else {
		err = config.DB.QueryRow(
			`SELECT userid, password, uname, roleid, last_login, last_os, last_browser, last_ip, role_name, role_description
			 FROM tv_user
			 WHERE uname = $1`,
			body.Username,
		).Scan(&id, &hash, &uname, &roleID, &Last_login, &Last_os, &Last_browser, &Last_ip, &Role_name, &Role_description)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return appErr.Unauthorized("Invalid credentials")
		}
		return appErr.FromDB(err)
	}

	if !roleID.Valid {
		return appErr.Forbidden(
			"User has no role assigned, contact your administrator",
		)
	}

	if err := utils.CheckPassword(hash, body.Password); err != nil {
		return appErr.Unauthorized("Invalid credentials")
	}

	access, err := utils.GenerateAccessToken(id)
	if err != nil {
		return appErr.Internal(err)
	}

	refresh, err := utils.GenerateRefreshToken(id)
	if err != nil {
		return appErr.Internal(err)
	}

	ip := c.IP()
	uaString := c.Get("User-Agent")
	ua := useragent.New(uaString)

	os := ua.OS()
	browser, _ := ua.Browser()

	_, err = config.DB.Exec(
		`UPDATE users
		 SET refresh_token = $1,
		     last_login = $2,
		     last_os=$3,
		     last_browser=$4,
		     last_ip=$5
		 WHERE userid = $6`,
		refresh,
		time.Now().UTC(),
		os,
		browser,
		ip,
		id,
	)

	if err != nil {
		return appErr.FromDB(err)
	}

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
		"data": fiber.Map{
			"username": uname,
			"last_login": Last_login,
			"last_os": Last_os,
			"last_browser": Last_browser,
			"last_ip": Last_ip,
			"role_name": Role_name,
			"role_description": Role_description,
		},
	})
}

func Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil {
		return appErr.InvalidJSON(err)
	}

	if body.RefreshToken == "" {
		return appErr.BadRequest("Refresh token is required")
	}

	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {

		// pastikan algoritma benar
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return config.RefreshSecret, nil
	})

	if err != nil || !token.Valid {
		return appErr.Unauthorized("Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return appErr.Unauthorized("Invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return appErr.Unauthorized("Invalid user id in token")
	}

	var savedToken string
	err = config.DB.QueryRow(
		"SELECT refresh_token FROM users WHERE userid=$1",
		userID,
	).Scan(&savedToken)

	if err != nil {
		if err == sql.ErrNoRows {
			return appErr.Unauthorized("User not found")
		}
		return appErr.FromDB(err)
	}

	if savedToken != body.RefreshToken {
		return appErr.Unauthorized("Refresh token revoked")
	}

	access, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return appErr.Internal(err)
	}

	newRefresh, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return appErr.Internal(err)
	}

	_, err = config.DB.Exec(
		"UPDATE users SET refresh_token=$1 WHERE userid=$2",
		newRefresh,
		userID,
	)

	if err != nil {
		return appErr.FromDB(err)
	}

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": newRefresh,
	})
}

func Logout(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return appErr.Unauthorized("Invalid or missing token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return appErr.Unauthorized("Invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return appErr.Unauthorized("Invalid user id")
	}

	_, err := config.DB.Exec(
		"UPDATE users SET refresh_token=NULL WHERE userid=$1",
		userID,
	)

	if err != nil {
		return appErr.FromDB(err)
	}

	return c.JSON(fiber.Map{
		"message": "Logout success",
	})
}