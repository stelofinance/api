package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/routes"
	"github.com/stelofinance/api/tools"
)

func main() {
	// Load in env
	if tools.LoadEnv() != nil {
		log.Fatal("Failed to load environment")
	}

	// Connect to db
	if database.ConnectDb() != nil {
		log.Fatal("Failed to connect to database")
	}

	// Run auto migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run auto migrations\n", err.Error())
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Stelo's Rebirth has begun!",
		})
	})

	// Setup routes
	routes.UsersRouter(app.Group("/users"))

	// Run app
	log.Fatal(app.Listen(":3000"))
}
