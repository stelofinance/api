package routes

import (
	"log"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/tools"
	"golang.org/x/crypto/bcrypt"
)

func postUser(c *fiber.Ctx) error {
	type requestBody struct {
		Username string `json:"username" validate:"required,min=2,max=32,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	var body requestBody

	// Parse and validate body
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

	// Add user and their wallet into database
	tx, err := database.DB.Begin(c.Context())
	defer tx.Rollback(c.Context())
	if err != nil {
		log.Printf("Error creating transaction: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	userID, err := qtx.InsertUser(c.Context(), db.InsertUserParams{
		Username:  body.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	})
	// TODO: Handle this cleaner, may not be username was taken
	if err != nil {
		return c.Status(400).SendString(constants.ErrorU000)
	}

	walletID, err := qtx.InsertWallet(c.Context(), db.InsertWalletParams{
		Address: tools.RandString(6),
		UserID:  userID,
	})
	// TODO: Handle this cleaner, the address may have been taken
	if err != nil {
		log.Printf("Error inserting wallet to db (addr taken?): {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	_, err = qtx.UpdateUserWallet(c.Context(), db.UpdateUserWalletParams{
		WalletID: pgtype.Int8{Int64: walletID, Valid: true},
		ID:       userID,
	})
	if err != nil {
		log.Printf("Error updating user wallet: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	return c.Status(201).SendString("User created")
}

func postSession(c *fiber.Ctx) error {
	type requestBody struct {
		Password    string `json:"password" validate:"required,max=32"`
		ReturnToken bool   `json:"return_token"`
	}

	var body requestBody

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Query for user
	user, err := database.Q.GetUser(c.Context(), c.Params("username"))

	// Check if user wasn't found
	if err != nil {
		return c.Status(404).SendString(constants.ErrorU003)
	}

	// Return error if password doesn't match hash
	isPasswordValid := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if isPasswordValid != nil {
		return c.Status(400).SendString(constants.ErrorU001)
	}

	key := uniuri.NewLen(27)

	// Add session into DB
	err = database.Q.InsertUserSession(c.Context(), db.InsertUserSessionParams{
		Key:      key,
		UserID:   user.ID,
		WalletID: user.WalletID.Int64,
	})
	if err != nil {
		log.Printf("Error inserting session: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Return in body if requested
	if body.ReturnToken {
		return c.Status(201).JSON(fiber.Map{
			"token": "stlu_" + key,
		})
	}

	// Create and set the cookie
	cookie := fiber.Cookie{
		Name:     "stelo_token",
		Expires:  time.Now().Add(time.Hour * 24 * 31),
		Value:    "stlu_" + key,
		Secure:   tools.EnvVars.ProductionEnv,
		HTTPOnly: true,
		SameSite: "strict",
	}
	c.Cookie(&cookie)

	return c.Status(201).SendString("Session created")
}
