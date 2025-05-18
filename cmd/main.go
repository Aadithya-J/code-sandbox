package main

import (
	"github.com/Aadithya-J/code-sandbox/internal/db"
	"github.com/Aadithya-J/code-sandbox/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	db.Init()
	router.SetupRoutes(app)

	app.Listen(":8080")
}
