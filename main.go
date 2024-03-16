package main

import (
	"log"

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

	log.Println("Listening on port 8000!")
	if err := app.Listen(":8000"); err != nil {
		log.Fatal("Failed to listen")
	}
}
