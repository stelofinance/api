package auth

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/tools"
)

type Role struct {
	slug string
}

var (
	Unknown = Role{""}
	Guest   = Role{"guest"}
	User    = Role{"user"}
	Wallet  = Role{"wallet"}
	Admin   = Role{"admin"}
)

func (r Role) String() string {
	return r.slug
}

type AuthHeader struct {
	Role  Role
	Token string
}

type UserJWT struct {
	UserID    int64 `json:"uid"`
	SessionID int64 `json:"sid"`
	WalletID  int64 `json:"wid"`
	jwt.StandardClaims
}

type WalletJWT struct {
	SessionID int64 `json:"sid"`
	WalletID  int64 `json:"wid"`
	jwt.StandardClaims
}

func getUserCookie(c *fiber.Ctx) (string, error) {
	userCookie := c.Cookies("ujwt")
	if userCookie == "" {
		return "", errors.New("No cookie")
	}
	return userCookie, nil
}

func getAuthHeader(c *fiber.Ctx) (AuthHeader, error) {
	header := c.GetReqHeaders()["Authorization"]
	headerArray := strings.Split(header, " ")
	if len(headerArray) != 2 {
		return AuthHeader{Role: Unknown, Token: ""},
			errors.New("Invalid header")
	}

	var role Role
	switch headerArray[0] {
	case Guest.slug:
		role = Guest
	case User.slug:
		role = User
	case Wallet.slug:
		role = Wallet
	case Admin.slug:
		role = Admin
	default:
		return AuthHeader{Role: Unknown, Token: ""},
			errors.New("Unknown role: " + headerArray[0])
	}

	return AuthHeader{Role: role, Token: headerArray[1]}, nil
}

// TODO: this is a mess and you know it
func New(r Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if r == Guest {
			if c.GetReqHeaders()["Authorization"] == "" &&
				c.Cookies("ujwt") == "" {
				return c.Next()
			}
			return c.Status(400).SendString(constants.ErrorA002)
		} else if r == Admin {
			header, err := getAuthHeader(c)
			if err != nil || header.Role != Admin || header.Token != tools.EnvVars.AdminKey {
				return c.Status(403).SendString(constants.ErrorA003)
			}
			return c.Next()
		} else if r == User {
			return AuthUser(c)
		} else if r == Wallet {
			if c.Cookies("ujwt") != "" {
				return AuthUser(c)
			}

			header, err := getAuthHeader(c)
			if err != nil {
				return c.Status(400).SendString(constants.ErrorA001)
			}

			if header.Role != Wallet {
				return c.Status(400).SendString(constants.ErrorA001)
			}

			claims := &WalletJWT{}
			token, _ := jwt.ParseWithClaims(header.Token, claims, func(token *jwt.Token) (interface{}, error) {
				return tools.EnvVars.JwtSecret, nil
			})

			if token.Valid {
				// Set the locals
				c.Locals("sid", claims.SessionID)
				c.Locals("wid", claims.WalletID)
				return c.Next()
			}
			return c.Status(400).SendString(constants.ErrorA001)
		}
		log.Fatal("Auth no Role matched")
		return nil
	}
}

// TODO: Please for the love of all things good clean up this
func AuthUser(c *fiber.Ctx) error {
	tokenString, err := getUserCookie(c)
	if err != nil {
		return c.Status(401).SendString(constants.ErrorA000)
	}

	claims := &UserJWT{}
	token, errJWT := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tools.EnvVars.JwtSecret, nil
	})

	if token.Valid {
		// Set the locals
		c.Locals("uid", claims.UserID)
		c.Locals("sid", claims.SessionID)
		c.Locals("wid", claims.WalletID)
		return c.Next()
	} else if errors.Is(errJWT, jwt.ErrTokenExpired) {
		// Update the session's used_at time
		// If there is an error or the session isn't
		// updated for whatever reason then return
		// invalid session
		rows, err := database.Q.UpdateUserSessionUsedAt(c.Context(), db.UpdateUserSessionUsedAtParams{
			UsedAt: time.Now(),
			ID:     claims.SessionID,
		})
		if err == nil && rows == 1 {
			// Create the JWT and set as cookie
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserJWT{
				UserID:    claims.UserID,
				SessionID: claims.SessionID,
				WalletID:  claims.WalletID,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
				},
			})
			jwtString, err := token.SignedString(tools.EnvVars.JwtSecret)
			if err != nil {
				log.Printf("Error creating JWT: {%v}", err.Error())
				return c.Status(500).SendString(constants.ErrorS000)
			}

			// Set the cookie
			cookie := fiber.Cookie{
				Name:     "ujwt",
				Value:    jwtString,
				Secure:   tools.EnvVars.ProductionEnv,
				HTTPOnly: true,
				SameSite: "Strict",
			}
			c.Cookie(&cookie)

			//  set the locals
			c.Locals("uid", claims.UserID)
			c.Locals("sid", claims.SessionID)
			c.Locals("wid", claims.WalletID)

			return c.Next()
		} else {
			return c.Status(400).SendString(constants.ErrorA001)
		}
	} else {
		// Clear the jwt cookie
		c.Cookie(&fiber.Cookie{
			Name:     "ujwt",
			Value:    "",
			Secure:   tools.EnvVars.ProductionEnv,
			HTTPOnly: true,
			SameSite: "strict",
			Expires:  time.Now().Add(-time.Hour),
		})

		return c.Status(400).SendString(constants.ErrorA001)
	}
}
