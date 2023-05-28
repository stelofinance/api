package routes

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
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
		log.Printf("Error hasing password: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Update the user's password
	err = database.Q.UpdateUserPassword(c.Context(), db.UpdateUserPasswordParams{
		Password: string(hashedPassword),
		ID:       c.Locals("uid").(int64),
	})

	if err != nil {
		log.Printf("Error updating password: {%v}", err.Error())
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
		log.Printf("Error creating db transaction: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	rows, err := qtx.UpdateUserWallet(c.Context(), db.UpdateUserWalletParams{
		WalletID: pgtype.Int8{Int64: body.WalletID, Valid: true},
		ID:       c.Locals("uid").(int64),
	})

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorU003)
	}

	if err != nil {
		log.Printf("Error updating user wallet: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	rows, err = qtx.CountWalletsByIdAndUserId(c.Context(), db.CountWalletsByIdAndUserIdParams{
		ID:     body.WalletID,
		UserID: c.Locals("uid").(int64),
	})

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorW000)
	}

	if err != nil {
		log.Printf("Error counting wallets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	return c.Status(200).SendString("Primary wallet updated")
}

func getWallets(c *fiber.Ctx) error {
	wallets, err := database.Q.GetWalletsByUserId(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Printf("Error retrieving wallets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Create the response object
	type walletsAPI struct {
		ID      int64  `json:"id"`
		Address string `json:"address"`
		UserID  int64  `json:"user_id"`
		Webhook string `json:"webhook"`
	}
	walletsResponse := []walletsAPI{}
	for _, wallet := range wallets {
		walletsResponse = append(walletsResponse, walletsAPI{
			ID:      wallet.ID,
			Address: wallet.Address,
			UserID:  wallet.UserID,
			Webhook: wallet.Webhook.String,
		})
	}

	return c.Status(200).JSON(walletsResponse)
}

func postWallet(c *fiber.Ctx) error {
	type requestBody struct {
		Address string `json:"address" validate:"min=3,max=24,alpha,lowercase"`
	}

	var body requestBody

	// Parse body then if they didn't set an address
	// gen one for them, otherwise validate theirs
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if body.Address == "" {
		body.Address = tools.RandString(6)
	} else if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Add the wallet into the DB
	err := database.Q.CreateWallet(c.Context(), db.CreateWalletParams{
		UserID:  c.Locals("uid").(int64),
		Address: body.Address,
	})
	if err != nil {
		log.Printf("Error creating wallet: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Wallet created",
		"address": body.Address,
	})
}

func getAssignedWallets(c *fiber.Ctx) error {
	wallets, err := database.Q.GetAssignedWalletsByUserId(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Printf("Error getting assigned wallets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if wallets == nil {
		wallets = []db.GetAssignedWalletsByUserIdRow{}
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
			return c.Status(404).SendString(constants.ErrorW000)
		}
	}

	// Update their session
	err := database.Q.UpdateUserSessionWallet(c.Context(), db.UpdateUserSessionWalletParams{
		ID:       c.Locals("sid").(int64),
		WalletID: body.WalletID,
	})
	if err != nil {
		log.Printf("Error updating session: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Active wallet updated")
}

func getSessions(c *fiber.Ctx) error {
	userSessions, err := database.Q.GetUserSessions(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Printf("Error getting sessions: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).JSON(userSessions)
}

func getSession(c *fiber.Ctx) error {
	result, err := database.Q.GetUserSessionInfo(c.Context(), db.GetUserSessionInfoParams{
		UserID:   c.Locals("uid").(int64),
		WalletID: c.Locals("wid").(int64),
	})

	if err != nil {
		log.Printf("Error getting session info: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).JSON(fiber.Map{
		"user_id":        c.Locals("uid").(int64),
		"username":       result.Username,
		"wallet_id":      c.Locals("wid").(int64),
		"wallet_address": result.WalletAddress,
	})
}

func deleteSession(c *fiber.Ctx) error {
	// Clear session from DB
	database.Q.DeleteSession(c.Context(), c.Locals("sid").(int64))

	// Clear the cookie
	domain := ".stelo.finance"
	if !tools.EnvVars.ProductionEnv {
		domain = "localhost"
	}

	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "stelo_token",
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
		log.Printf("Error deleting session: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	} else if count == 0 {
		return c.Status(404).SendString(constants.ErrorU004)
	}

	if int64(id) == c.Locals("sid").(int64) {
		// Clear the cookie
		domain := ".stelo.finance"
		if !tools.EnvVars.ProductionEnv {
			domain = "localhost"
		}
		c.Cookie(&fiber.Cookie{
			Domain:   domain,
			Name:     "stelo_token",
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
	domain := ".stelo.finance"
	if !tools.EnvVars.ProductionEnv {
		domain = "localhost"
	}
	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "stelo_token",
		Value:    "",
		Secure:   tools.EnvVars.ProductionEnv,
		HTTPOnly: true,
		SameSite: "strict",
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.Status(200).SendString("All sessions deleted")
}
