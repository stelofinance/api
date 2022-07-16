package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Stelo's Rebirth has begun!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
