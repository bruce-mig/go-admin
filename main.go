package main

import (
	"github.com/bruce-mig/go-admin/db"
	"github.com/bruce-mig/go-admin/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:8000",
	}))

	routes.SetupRoutes(app)

	println("Listening on port 8080!")
	app.Listen(":8000")
}
