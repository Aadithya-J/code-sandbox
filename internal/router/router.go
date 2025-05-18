package router

import (
	"github.com/Aadithya-J/code-sandbox/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")

	auth.Post("/register", handler.UserRegister)
	auth.Post("/login", handler.UserLogin)
}
