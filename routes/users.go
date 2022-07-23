package routes

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/models"
	"github.com/stelofinance/api/tools"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func postUser(c *fiber.Ctx) error {
	type requestBody struct {
		Username string `json:"username" validate:"required,min=2,max=32,alphanum"`
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

func postSession(c *fiber.Ctx) error {
	type requestBody struct {
		Password string `json:"password" validate:"required,max=32"`
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

	// Query for user
	var user models.User
	result := database.Db.Select("id", "password").Where("username = ?", c.Params("username")).First(&user)

	// Check if user wasn't found
	if result.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"code": "U0001",
		})
	}

	// Return error if password doesn't match hash
	isPasswordValid := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if isPasswordValid != nil {
		return c.Status(400).JSON(fiber.Map{
			"code": "U0001",
		})
	}

	// Create the JWT and set as cookie
	// Note: using 'exp' as the field name
	// makes the jwt module auto check it's
	// expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	})
	jwtSecret, err := tools.GetEnvVariable("JWT_SECRET")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code": "S0000",
		})
	}
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code": "S0000",
		})
	}

	// Set the cookie
	isProdEnv := false
	prodEnv, err := tools.GetEnvVariable("PRODUCTION_ENV")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code": "S0000",
		})
	} else if prodEnv == "true" {
		isProdEnv = true
	}
	cookie := fiber.Cookie{
		Name:     "sjwt",
		Value:    tokenString,
		Secure:   isProdEnv,
		HTTPOnly: true,
		SameSite: "Strict",
		Expires:  time.Now().Add(time.Minute * 30),
	}
	c.Cookie(&cookie)

	return c.Status(201).SendString("Session created")
}

func testCookie(c *fiber.Ctx) error {
	return c.Status(200).SendString("You're logged in and valid!")
}
