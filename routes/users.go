package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func postUser(c *fiber.Ctx) error {
	type requestBody struct {
		Username string `json:"username" validate:"required,min=2,max=32"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	var body requestBody

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).JSON(fiber.Map{
			"code": "G0000",
		})
	}
	if validate.Struct(body) != nil {
		return c.Status(400).JSON(fiber.Map{
			"code": "G0000",
		})
	}

	return c.Status(200).JSON(body)
}
