package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/models"
	"github.com/stelofinance/api/tools"
	"golang.org/x/crypto/bcrypt"
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code": "S0000",
		})
	}

	// Create the user model and insert into db
	user := models.User{
		Username: body.Username,
		Password: string(hashedPassword),
		Wallet: models.Wallet{
			Address: tools.RandString(6),
		},
	}
	result := database.Db.Create(&user)

	// If error on insertion assume username was taken
	// TODO: handle this cleaner, check error, maybe address already taken?
	if result.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"code": "U0000",
		})
	}

	return c.Status(201).SendString("User created")
}
