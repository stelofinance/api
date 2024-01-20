package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/pusher"
)

func postWarehouse(c *fiber.Ctx) error {
	var body struct {
		Name        string   `json:"name" validate:"required,min=2,max=32"`
		Coordinates [2]int64 `json:"coordinates"`
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

	// Check they have permission to create the warehouse
	canCreateWarehouses, err := qtx.GetUserCanCreateWarehouses(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Println("Error getting can create warehouses permission", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if !canCreateWarehouses {
		return c.Status(403).SendString(constants.ErrorH000)
	}

	// Create the warehouse
	warehouseId, err := qtx.InsertWarehouse(
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

	// Add owner to warehouse's workers
	err = qtx.InsertWarehouseWorker(c.Context(), db.InsertWarehouseWorkerParams{
		WarehouseID: warehouseId,
		UserID:      c.Locals("uid").(int64),
	})
	if err != nil {
		log.Println("Error adding owner to new warehouse's workers", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

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

	// Make sure requester is owner of warehouse
	warehouseId, err := strconv.Atoi(c.Params("warehouseid"))
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}
	userId, err := qtx.GetWarehouseUserIdLock(c.Context(), int64(warehouseId))
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
		rows, err = qtx.AddWarehouseCollateral(c.Context(), db.AddWarehouseCollateralParams{
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
		rows, err = qtx.SubtractWarehouseCollateral(c.Context(), db.SubtractWarehouseCollateralParams{
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

func putWarehouseOwner(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username" validate:"required"`
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

	// Make sure requester is owner of warehouse
	warehouseId, err := strconv.ParseInt(c.Params("warehouseid"), 10, 64)
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}
	userId, err := qtx.GetWarehouseUserIdLock(c.Context(), int64(warehouseId))
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

	// Make sure new owner is a worker
	isWorker, err := qtx.ExistsWarehouseWorker(c.Context(), db.ExistsWarehouseWorkerParams{
		WarehouseID: warehouseId,
		Username:    body.Username,
	})
	if err != nil {
		log.Println("Error checking is new owner is a worker", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if !isWorker {
		return c.Status(400).SendString(constants.ErrorH004)
	}

	// Update user who owns warehouse
	err = qtx.UpdateWarehouseUserIdByUsername(c.Context(), db.UpdateWarehouseUserIdByUsernameParams{
		ID:       warehouseId,
		Username: body.Username,
	})
	if err != nil {
		log.Printf("Error updating warehouse user id: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	return c.Status(200).SendString("Warehouse owner updated")
}

func postWarehouseWorker(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username" validate:"required"`
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

	// Make sure requester is owner of warehouse
	warehouseId, err := strconv.ParseInt(c.Params("warehouseid"), 10, 64)
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}
	userId, err := qtx.GetWarehouseUserIdLock(c.Context(), int64(warehouseId))
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

	// Assign user to warehouse
	err = qtx.InsertWarehouseWorkerByUsername(c.Context(), db.InsertWarehouseWorkerByUsernameParams{
		WarehouseID: warehouseId,
		Username:    body.Username,
	})
	if err != nil {
		log.Printf("Error updating warehouse user id: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	return c.Status(201).SendString("Warehouse worker added to warehouse")
}

func deleteWarehouseWorker(c *fiber.Ctx) error {
	workerId, err := strconv.ParseInt(c.Params("workerid"), 10, 64)
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	// Delete worker
	err = database.Q.DeleteWarehouseWorker(c.Context(), db.DeleteWarehouseWorkerParams{
		ID:     workerId,
		UserID: c.Locals("uid").(int64),
	})
	if err != nil {
		log.Println("Error deleting warehouse worker", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Worker removed from warehouse")
}

func getWarehouseWorkers(c *fiber.Ctx) error {
	warehouseId, err := strconv.ParseInt(c.Params("warehouseid"), 10, 64)
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	workers, err := database.Q.GetWarehouseWorkers(c.Context(), warehouseId)
	if err != nil {
		log.Println("Error getting warehouse workers", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).JSON(workers)
}

func postWarehouseAssets(c *fiber.Ctx) error {
	var body struct {
		Recipient string           `json:"recipient" validate:"required"`
		Type      uint8            `json:"type" validate:"lte=2"`
		Memo      string           `json:"memo" validate:"max=64"`
		Assets    map[string]int64 `json:"assets" validate:"gt=0,dive,keys,ne=stelo,endkeys,gt=0"`
	}

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	warehouseId, err := strconv.ParseInt(c.Params("warehouseid"), 10, 64)
	if err != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	// Default the memo if there is none
	if body.Memo == "" {
		body.Memo = "warehouse deposit"
	}

	tx, err := database.DB.Begin(c.Context())
	defer tx.Rollback(c.Context())
	if err != nil {
		log.Printf("Error creating db transaction: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	// Get asset ids
	var assetNames []string
	for asset := range body.Assets {
		if asset == "stelo" {
			return c.Status(500).SendString(constants.ErrorS000)
		}
		assetNames = append(assetNames, asset)
	}
	assets, err := qtx.GetAssetsByNames(c.Context(), assetNames)
	if err != nil {
		log.Printf("Error getting assets id & name: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if len(assets) != len(body.Assets) {
		return c.Status(404).SendString(constants.ErrorI000)
	}
	assetsMap := make(map[string]db.Asset)
	for _, asset := range assets {
		assetsMap[asset.Name] = asset
	}

	// Make sure warehouse has enough collateral
	var collateralNeeded int64 = 0
	for asset, quantity := range body.Assets {
		collateralNeeded += assetsMap[asset].Value * quantity
	}
	warehouseResult, err := qtx.GetWarehouseCollateralLiabilityAndRatioLock(c.Context(), warehouseId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(404).SendString(constants.ErrorH001)
		}
		log.Println("Unable to fetch warehouse", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	collateralRatio, err := warehouseResult.CollateralRatio.Float64Value()
	log.Println("Collateral ratio", collateralRatio)
	if err != nil {
		log.Println("Error getting collateral ratio", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	availableCollateral := float64(warehouseResult.Collateral) - float64(warehouseResult.Liability)*collateralRatio.Float64
	log.Println("Available collateral", availableCollateral)
	log.Println("Collateral actually needed", float64(collateralNeeded)*collateralRatio.Float64)
	if availableCollateral < float64(collateralNeeded)*collateralRatio.Float64 {
		return c.Status(400).SendString(constants.ErrorH003)
	}

	// Get recipient's walletId
	var recipientWalletID int64
	webhook := pgtype.Text{
		String: "",
		Valid:  false,
	}
	// Username type
	if body.Type == 0 {
		walletID, err := qtx.GetWalletByUsername(c.Context(), body.Recipient)
		if err != nil {
			return c.Status(404).SendString(constants.ErrorU003)
		}
		if !walletID.Valid {
			log.Printf("user created without wallet")
			return c.Status(500).SendString(constants.ErrorS000)
		}
		recipientWalletID = walletID.Int64
		// Address type
	} else if body.Type == 1 {
		wallet, err := qtx.GetWalletIdAndWebhookByAddress(c.Context(), body.Recipient)
		if err != nil {
			return c.Status(404).SendString(constants.ErrorW000)
		}
		recipientWalletID = wallet.ID
		webhook = wallet.Webhook
		// Wallet Id type
	} else {
		wallet_id, err := strconv.ParseInt(body.Recipient, 10, 0)
		if err != nil {
			return c.Status(400).SendString(constants.ErrorG000)
		}
		wallet_webhook, err := qtx.GetWalletWebhook(c.Context(), wallet_id)
		recipientWalletID = wallet_id
		webhook = wallet_webhook
	}

	// If there is a webhook hit it first
	if webhook.Valid {
		// Create the body
		postBody, err := json.Marshal(fiber.Map{
			"wallet_id": c.Locals("wid").(int64),
			"memo":      body.Memo,
			"assets":    body.Assets,
		})
		if err != nil {
			log.Printf("Error creating json body: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		}
		responseBody := bytes.NewBuffer(postBody)

		// Hit the webhook
		resp, err := http.Post(webhook.String, "application/json", responseBody)
		if err != nil {
			log.Printf("Error posting to webhook: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		}
		if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
			return c.Status(400).SendString(constants.ErrorW010)
		}
		resp.Body.Close()
	}

	// Deposit assets into recipient's account
	for asset, quantity := range body.Assets {
		rows, err := qtx.AddWalletAssetQuantity(c.Context(), db.AddWalletAssetQuantityParams{
			Quantity: quantity,
			WalletID: recipientWalletID,
			AssetID:  assetsMap[asset].ID,
		})

		if err != nil {
			log.Printf("Error adding asset quantity: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			err := qtx.CreateWalletAsset(c.Context(), db.CreateWalletAssetParams{
				Quantity: quantity,
				WalletID: recipientWalletID,
				AssetID:  assetsMap[asset].ID,
			})

			if err != nil {
				log.Printf("Error creating wallet asset: {%v}", err.Error())
				return c.Status(500).SendString(constants.ErrorS000)
			}
		}
	}

	// Create transaction record
	transactionID, err := qtx.CreateTransaction(c.Context(), db.CreateTransactionParams{
		SendingWalletID:   recipientWalletID,
		ReceivingWalletID: recipientWalletID,
		CreatedAt:         time.Now(),
		Memo: pgtype.Text{
			String: body.Memo,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("Error creating transaction record: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	var transactionAssets []db.CreateTransactionAssetsParams
	for asset, quantity := range body.Assets {
		transactionAssets = append(transactionAssets, db.CreateTransactionAssetsParams{
			TransactionID: transactionID,
			AssetID:       assetsMap[asset].ID,
			Quantity:      quantity,
		})
	}

	txAssetsResult := qtx.CreateTransactionAssets(c.Context(), transactionAssets)

	var insertErrorOccured bool
	txAssetsResult.Exec(func(i int, err error) {
		if err != nil {
			insertErrorOccured = true
		}
	})
	if insertErrorOccured {
		log.Printf("Error inserting transaction assets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Deposit assets into warehouse
	for asset, quantity := range body.Assets {
		rows, err := qtx.AddWarehouseAssetQuantity(c.Context(), db.AddWarehouseAssetQuantityParams{
			Quantity:    quantity,
			WarehouseID: warehouseId,
			AssetID:     assetsMap[asset].ID,
		})

		if err != nil {
			log.Printf("Error adding warehouse asset quantity: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			err := qtx.CreateWarehouseAsset(c.Context(), db.CreateWarehouseAssetParams{
				Quantity:    quantity,
				WarehouseID: warehouseId,
				AssetID:     assetsMap[asset].ID,
			})

			if err != nil {
				log.Printf("Error creating warehouse asset: {%v}", err.Error())
				return c.Status(500).SendString(constants.ErrorS000)
			}
		}
	}

	// Adjust warehouse liability
	err = qtx.AddWarehouseLiabiliy(c.Context(), db.AddWarehouseLiabiliyParams{
		ID:        warehouseId,
		Liability: collateralNeeded,
	})
	if err != nil {
		log.Println("Unable to add warehouse liability", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	// Send to pusher cannel
	go func(recipient int64, sender int64, memo string, assets map[string]int64) {
		var data fiber.Map

		if memo != "" {
			data = fiber.Map{
				"sender": sender,
				"memo":   memo,
				"assets": assets,
			}
		} else {
			data = fiber.Map{
				"sender": sender,
				"assets": assets,
			}
		}

		err := pusher.PusherClient.Trigger("private-wallet@"+fmt.Sprint(recipient), "transaction:incoming", data)
		if err != nil {
			log.Printf("Error posting transaction to Pusher: {%v}", err.Error())
		}
	}(recipientWalletID, recipientWalletID, body.Memo, body.Assets)

	return c.Status(201).SendString("Assets deposited, transaction created")
}
