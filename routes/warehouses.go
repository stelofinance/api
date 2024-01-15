package routes

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
)

func postWarehouse(c *fiber.Ctx) error {
	var body struct {
		Name        string   `json:"name" validate:"required"`
		Coordinates [2]int64 `json:"coordinates"`
	}

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Check they have permission to create the warehouse
	canCreateWarehouses, err := database.Q.GetUserCanCreateWarehouses(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Println("Error getting can create warehouses permission", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if !canCreateWarehouses {
		return c.Status(403).SendString(constants.ErrorH000)
	}

	// Create the warehouse
	err = database.Q.InsertWarehouse(
		c.Context(),
		db.InsertWarehouseParams{
			Name:     body.Name,
			UserID:   c.Locals("uid").(int64),
			Location: fmt.Sprintf("POINT(%d %d)", body.Coordinates[0], body.Coordinates[1]),
		})
	if err != nil {
		log.Println("Error creating warehouse", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(201).SendString("Warehouse created")
}

func putCollateral(c *fiber.Ctx) error {
	var body struct {
		Amount int64 `json:"amount" validate:"ne=0"`
	}

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	tx, err := database.DB.Begin(c.Context())
	defer tx.Rollback(c.Context())
	if err != nil {
		log.Printf("Error creating db transaction: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	// Make sure requeser is owner of warehouse
	warehouseId, err := strconv.Atoi(c.Params("warehouseid"))
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}
	userId, err := qtx.GetWarehouseUserId(c.Context(), int64(warehouseId))
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).SendString(constants.ErrorH001)
		}
		log.Println("Failed to fetch user id from warehouse", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if userId != c.Locals("uid").(int64) {
		return c.Status(403).SendString(constants.ErrorH002)
	}

	// Find Stelo asset id
	steloId, err := qtx.GetAssetIdByName(c.Context(), "stelo")
	if err != nil {
		log.Println("Error getting stelo asset id", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Adding collateral
	if body.Amount > 0 {
		// Remove stelo from wallet
		rows, err := qtx.SubtractWalletAssetQuantity(c.Context(), db.SubtractWalletAssetQuantityParams{
			Quantity: body.Amount,
			WalletID: c.Locals("wid").(int64),
			AssetID:  steloId,
		})

		if err != nil {
			log.Printf("Error subtracting quantity from stelo: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			return c.Status(400).SendString(constants.ErrorI001)
		}

		// Create transaction record
		transactionId, err := qtx.CreateTransaction(c.Context(), db.CreateTransactionParams{
			SendingWalletID:   c.Locals("wid").(int64),
			ReceivingWalletID: c.Locals("wid").(int64),
			CreatedAt:         time.Now(),
			Memo: pgtype.Text{
				String: "Warehouse collateral transfer",
				Valid:  true,
			},
		})
		if err != nil {
			log.Println("Unable to create transaction")
			return c.Status(500).SendString(constants.ErrorS000)
		}
		err = qtx.CreateTransactionAsset(c.Context(), db.CreateTransactionAssetParams{
			TransactionID: transactionId,
			AssetID:       steloId,
			Quantity:      body.Amount,
		})
		if err != nil {
			log.Println("Unable to create transaction asset")
			return c.Status(500).SendString(constants.ErrorS000)
		}

		// Add collateral to warehouse
		rows, err = qtx.AddWarehouseCollateralQuantity(c.Context(), db.AddWarehouseCollateralQuantityParams{
			ID:         int64(warehouseId),
			Collateral: body.Amount,
		})

		if err != nil {
			log.Println("Unable to adjust warehouse collateral")
			return c.Status(500).SendString(constants.ErrorS000)
		}
		if rows == 0 {
			return c.Status(404).SendString(constants.ErrorH001)
		}

	} else {
		amount := -body.Amount
		// Add stelo to wallet
		rows, err := qtx.AddWalletAssetQuantity(c.Context(), db.AddWalletAssetQuantityParams{
			Quantity: amount,
			WalletID: c.Locals("wid").(int64),
			AssetID:  steloId,
		})

		if err != nil {
			log.Printf("Error adding asset quantity: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			err := qtx.CreateWalletAsset(c.Context(), db.CreateWalletAssetParams{
				Quantity: amount,
				WalletID: c.Locals("wid").(int64),
				AssetID:  steloId,
			})

			if err != nil {
				log.Printf("Error creating wallet asset: {%v}", err.Error())
				return c.Status(500).SendString(constants.ErrorS000)
			}
		}

		// Create transaction record
		transactionId, err := qtx.CreateTransaction(c.Context(), db.CreateTransactionParams{
			SendingWalletID:   c.Locals("wid").(int64),
			ReceivingWalletID: c.Locals("wid").(int64),
			CreatedAt:         time.Now(),
			Memo: pgtype.Text{
				String: "Warehouse collateral transfer",
				Valid:  true,
			},
		})
		if err != nil {
			log.Println("Unable to create transaction")
			return c.Status(500).SendString(constants.ErrorS000)
		}
		err = qtx.CreateTransactionAsset(c.Context(), db.CreateTransactionAssetParams{
			TransactionID: transactionId,
			AssetID:       steloId,
			Quantity:      amount,
		})
		if err != nil {
			log.Println("Unable to create transaction asset")
			return c.Status(500).SendString(constants.ErrorS000)
		}

		// Remove collateral from warehouse
		rows, err = qtx.SubtractWarehouseCollateralQuantity(c.Context(), db.SubtractWarehouseCollateralQuantityParams{
			ID:         int64(warehouseId),
			Collateral: amount,
		})

		if err != nil {
			log.Println("Unable to adjust warehouse collateral")
			return c.Status(500).SendString(constants.ErrorS000)
		}
		if rows == 0 {
			return c.Status(400).SendString(constants.ErrorH003)
		}
	}

	tx.Commit(c.Context())

	return c.Status(200).SendString("Collateral adjusted")
}
