package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
	auth "github.com/stelofinance/api/middlewares"
	"github.com/stelofinance/api/tools"
	"golang.org/x/crypto/bcrypt"
)

func putUsername(c *fiber.Ctx) error {
	// Parse and validate body
	type requestBody struct {
		Username string `json:"username" validate:"required,min=2,max=32,alphanum"`
	}
	var body requestBody
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	err := database.Q.UpdateUserUsername(c.Context(), db.UpdateUserUsernameParams{
		Username: body.Username,
		ID:       c.Locals("uid").(int64),
	})
	if err != nil {
		return c.Status(400).SendString(constants.ErrorU000)
	}

	return c.Status(200).SendString("Username updated")
}

func putPassword(c *fiber.Ctx) error {
	// Parse and validate body
	type requestBody struct {
		Password string `json:"password" validate:"required,min=8,max=32"`
	}
	var body requestBody
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Update the user's password
	err = database.Q.UpdateUserPassword(c.Context(), db.UpdateUserPasswordParams{
		Password: string(hashedPassword),
		ID:       c.Locals("uid").(int64),
	})

	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Password updated")
}

func putWallet(c *fiber.Ctx) error {
	// Parse and validate body
	type requestBody struct {
		WalletID int64 `json:"wallet_id" validate:"required"`
	}
	var body requestBody
	if c.BodyParser(&body) != nil || validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Change the primary wallet ID, if they didn't own it
	// then revert the transaction
	tx, err := database.DB.Begin(c.Context())
	defer tx.Rollback(c.Context())
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	rows, err := qtx.UpdateUserWallet(c.Context(), db.UpdateUserWalletParams{
		WalletID: sql.NullInt64{Int64: body.WalletID, Valid: true},
		ID:       c.Locals("uid").(int64),
	})

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorU003)
	}

	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	rows, err = qtx.CountWalletsByIdAndUserId(c.Context(), db.CountWalletsByIdAndUserIdParams{
		ID:     body.WalletID,
		UserID: c.Locals("uid").(int64),
	})

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorU002)
	}

	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	return c.Status(200).SendString("Primary wallet updated")
}

func getWallets(c *fiber.Ctx) error {
	wallets, err := database.Q.GetWalletsByUserId(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).JSON(wallets)
}

func getAssignedWallets(c *fiber.Ctx) error {
	wallets, err := database.Q.GetAssignedWalletsByUserId(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if wallets == nil {
		wallets = []db.Wallet{}
	}

	return c.Status(200).JSON(wallets)
}

func putActiveWallet(c *fiber.Ctx) error {
	body := struct {
		WalletID int64 `json:"wallet_id"`
	}{}

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Check if they own the wallet
	// UNSAFE: Type assertion could panic
	rows, _ := database.Q.CountWalletsByIdAndUserId(c.Context(), db.CountWalletsByIdAndUserIdParams{
		ID:     body.WalletID,
		UserID: c.Locals("uid").(int64),
	})
	if rows == 0 {
		// UNSAFE: Type assertion could panic
		rows, _ := database.Q.CountAssignedWallet(c.Context(), db.CountAssignedWalletParams{
			UserID:   c.Locals("uid").(int64),
			WalletID: body.WalletID,
		})
		if rows == 0 {
			return c.Status(404).SendString(constants.ErrorU002)
		}
	}

	// Create the JWT and set as cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.UserJWT{
		UserID:    c.Locals("uid").(int64), // UNSAFE: Type assertion could panic
		SessionID: c.Locals("sid").(int64), // UNSAFE: Type assertion could panic
		WalletID:  body.WalletID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	})
	jwtString, err := token.SignedString(tools.EnvVars.JwtSecret)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Update their session
	// UNSAFE: Type assertion could panic
	err = database.Q.UpdateUserSessionWallet(c.Context(), db.UpdateUserSessionWalletParams{
		ID:       c.Locals("sid").(int64),
		WalletID: body.WalletID,
	})
	if err != nil {
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

	return c.Status(200).SendString("Active wallet updated")
}

func getSessions(c *fiber.Ctx) error {
	userSessions, err := database.Q.GetUserSessions(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).JSON(userSessions)
}

func deleteSession(c *fiber.Ctx) error {
	// Clear session from DB
	database.Q.DeleteSession(c.Context(), c.Locals("sid").(int64))

	// Clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "ujwt",
		Value:    "",
		Secure:   tools.EnvVars.ProductionEnv,
		HTTPOnly: true,
		SameSite: "strict",
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.Status(200).SendString("Session deleted")
}

func deleteSessionById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("sessionid")
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	count, err := database.Q.DeleteSession(c.Context(), int64(id))
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	} else if count == 0 {
		return c.Status(404).SendString(constants.ErrorU004)
	}

	if int64(id) == c.Locals("sid").(int64) {
		// Clear the cookie
		c.Cookie(&fiber.Cookie{
			Name:     "ujwt",
			Value:    "",
			Secure:   tools.EnvVars.ProductionEnv,
			HTTPOnly: true,
			SameSite: "strict",
			Expires:  time.Now().Add(-time.Hour),
		})
	}

	return c.Status(200).SendString("Session deleted")
}

func deleteSessions(c *fiber.Ctx) error {
	// Clear all sessions from DB related to the user
	database.Q.DeleteSessionsByUserId(c.Context(), c.Locals("uid").(int64))

	// Clear their session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "ujwt",
		Value:    "",
		Secure:   tools.EnvVars.ProductionEnv,
		HTTPOnly: true,
		SameSite: "strict",
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.Status(200).SendString("All sessions deleted")
}
