package routes

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/tools"
)

type centrifugoJWT struct {
	Channels []string `json:"channels"`
	jwt.StandardClaims
}

func getCentrifugoToken(c *fiber.Ctx) error {
	widString := strconv.FormatInt(c.Locals("wid").(int64), 10)

	// Create the JWT
	claims := &centrifugoJWT{
		Channels: []string{"wallet:transactions#" + widString},
		StandardClaims: jwt.StandardClaims{
			Subject:   widString,
			ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString(tools.EnvVars.CentrifugoJWTKey)
	if err != nil {
		log.Printf("Error creating JWT: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).JSON(fiber.Map{
		"token": jwtString,
	})
}
