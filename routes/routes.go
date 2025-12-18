package routes

import (
	"github.com/gofiber/fiber/v2"

	"web-service/handlers"
	"web-service/middleware"
)

func Setup(app *fiber.App) {
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Post("/refresh", handlers.Refresh)

	api := app.Group("/api", middleware.JWTProtected())
		api.Get("/who-is", handlers.Me)
		api.Post("/logout", handlers.Logout)
}
