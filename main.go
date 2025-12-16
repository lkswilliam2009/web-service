package main

import (
	"github.com/gofiber/fiber/v2"

	"web-service/config"
	"web-service/routes"
)

func main() {
	config.ConnectDB()

	app := fiber.New()
	routes.Setup(app)

	app.Listen(":3000")
}
