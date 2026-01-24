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

		// MASTER DATA EMPLOYEE ENDPOINT
		api.Get("/employee-data/master/list", handlers.EmployeeData)
		api.Get("/employee-data/master/detail", handlers.EmployeeDataDetail)
		api.Post("/employee-data/master/create", handlers.EmployeeInsert)
		api.Post("/employee-data/master/edit", handlers.EmployeeUpdate)
		api.Post("/employee-data/master/remove", handlers.EmployeeSoftDelete)

		// GLOBAL FILE UPLOAD DOWNLOAD ENDPOINT
		api.Post("/file-managers/upload", handlers.UploadDokumen)
		api.Get("/file-managers/document/:pegawai_id/:doc_type", handlers.GetDokumen)
}