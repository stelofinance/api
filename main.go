package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/database"
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

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Stelo's Rebirth has begun!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
