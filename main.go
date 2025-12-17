package main

import (
	"github.com/gofiber/fiber/v2"

	"web-service/config"
	"web-service/routes"
	"web-service/middleware"
)

func main() {
	config.ConnectDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	routes.Setup(app)

	app.Listen(":3000")
}
