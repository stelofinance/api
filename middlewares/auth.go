package middlewares

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stelofinance/api/tools"
)

// TODO: Optimize the heck out of this
func AuthSession(c *fiber.Ctx) error {
	cookie := c.Cookies("sjwt")
	if cookie == "" {
		return c.Status(400).JSON(fiber.Map{
			"code": "A0000",
		})
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// Get secret
		// TODO: Maybe get once at app start, and use cache?
		secret, err := tools.GetEnvVariable("JWT_SECRET")
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve secret")
		}
		return []byte(secret), nil
	})
	if !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(400).JSON(fiber.Map{
				"code": "A0002",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"code": "A0001",
		})
	}

	return c.Next()
}
