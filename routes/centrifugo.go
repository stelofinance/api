package routes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func postConnection(c *fiber.Ctx) error {
	widString := strconv.FormatInt(c.Locals("wid").(int64), 10)

	var channels [1]string
	channels[0] = "wallet:transactions#" + widString
	return c.Status(200).JSON(fiber.Map{
		"result": fiber.Map{
			"user":     widString,
			"channels": channels,
		},
	})
}
