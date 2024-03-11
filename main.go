package main

import (
	"github.com/bruce-mig/go-admin/db"
	"github.com/bruce-mig/go-admin/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db.Connect()

	app := fiber.New()

	routes.SetupRoutes(app)

	app.Listen(":8000")
	println("Listening on port 8080!")
}
