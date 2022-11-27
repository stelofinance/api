package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/tools"
)

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
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Wallet created",
		"address": body.Address,
	})
}

func postAssetToWallet(c *fiber.Ctx) error {
	// Get walletid param
	params := struct {
		WalletID int64 `params:"walletid"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	body := struct {
		ID       int64 `json:"id" validate:"required"`
		Quantity int64 `json:"quantity" validate:"required"`
	}{}

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	err := database.Q.CreateWalletAsset(c.Context(), db.CreateWalletAssetParams{
		WalletID: params.WalletID,
		AssetID:  body.ID,
		Quantity: body.Quantity,
	})

	if err != nil {
		return c.Status(404).SendString(constants.ErrorW002)
	}

	return c.Status(200).SendString("Asset added wallet")
}

func deleteAssetFromWallet(c *fiber.Ctx) error {
	// Get walletid and assetid param
	params := struct {
		WalletID int64 `params:"walletid"`
		AssetID  int64 `params:"assetid"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	rows, _ := database.Q.DeleteWalletAsset(c.Context(), db.DeleteWalletAssetParams{
		WalletID: params.WalletID,
		AssetID:  params.AssetID,
	})

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorI000)
	}

	return c.Status(200).SendString("Asset removed from wallet")
}