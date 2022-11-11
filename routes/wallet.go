package routes

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/middlewares"
	"github.com/stelofinance/api/tools"
)

func getAssets(c *fiber.Ctx) error {
	// Get their wallet's assets
	walletAssetsResult, err := database.Q.GetWalletAssets(c.Context(), c.Locals("wid").(int64))

	// Create array of the asset IDs
	var assetIDs []int64
	for _, walletAsset := range walletAssetsResult {
		assetIDs = append(assetIDs, walletAsset.AssetID)
	}

	// Get all the assets (preloading basically)
	assets, err := database.Q.GetAssetsByIds(c.Context(), assetIDs)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Create a map of the assets using their id
	assetMap := make(map[int64]db.Asset)
	for _, asset := range assets {
		assetMap[asset.ID] = asset
	}

	// Create the response object
	type walletAssetAPI struct {
		Asset    db.Asset `json:"asset"`
		Quantity int64    `json:"quantity"`
	}
	walletAssetsAPI := []walletAssetAPI{}
	for _, walletAsset := range walletAssetsResult {
		walletAssetsAPI = append(walletAssetsAPI, walletAssetAPI{
			Asset:    assetMap[walletAsset.AssetID],
			Quantity: walletAsset.Quantity,
		})
	}

	return c.Status(200).JSON(walletAssetsAPI)
}

func postTransaction(c *fiber.Ctx) error {
	type requestBody struct {
		Recipient  string           `json:"recipient" validate:"required"`
		IsUsername bool             `json:"is_username"`
		Memo       string           `json:"memo" validate:"max=64"`
		Assets     map[string]int64 `json:"assets" validate:"gt=0,dive,gt=0"`
	}
	var body requestBody

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
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	// TODO: Optimize this whole thing maybe? Perhaps channels for each asset
	// so big transactions don't take 200+ ms
	// Create a string int64 map for asset names to ids
	var assetNames []string
	for asset := range body.Assets {
		assetNames = append(assetNames, asset)
	}
	assets, err := qtx.GetAssetsIdNameByNames(c.Context(), assetNames)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	assetIDs := make(map[string]int64)
	for _, asset := range assets {
		assetIDs[asset.Name] = asset.ID
	}

	// Remove assets from requesters wallet
	for asset, quantity := range body.Assets {
		if assetIDs[asset] == 0 {
			return c.Status(400).SendString(constants.ErrorI001)
		}
		rows, err := qtx.SubtractWalletAssetQuantity(c.Context(), db.SubtractWalletAssetQuantityParams{
			Quantity: quantity,
			WalletID: c.Locals("wid").(int64),
			AssetID:  assetIDs[asset],
		})

		if err != nil {
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			return c.Status(400).SendString(constants.ErrorI001)
		}
	}

	// Now put assets in other wallet, first get their wallet ID though
	// Create the wallet asset record if not already
	var recipientWalletID int64
	if body.IsUsername {
		walletID, err := qtx.GetWalletByUsername(c.Context(), body.Recipient)
		if err != nil || !walletID.Valid {
			return c.Status(500).SendString(constants.ErrorU003)
		}
		recipientWalletID = walletID.Int64
	} else {
		walletID, err := qtx.GetWalletIdByAddress(c.Context(), body.Recipient)
		if err != nil {
			return c.Status(500).SendString(constants.ErrorW000)
		}
		recipientWalletID = walletID
	}

	for asset, quantity := range body.Assets {
		rows, err := qtx.AddWalletAssetQuantity(c.Context(), db.AddWalletAssetQuantityParams{
			Quantity: quantity,
			WalletID: recipientWalletID,
			AssetID:  assetIDs[asset],
		})

		if err != nil {
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			err := qtx.CreateWalletAsset(c.Context(), db.CreateWalletAssetParams{
				Quantity: quantity,
				WalletID: recipientWalletID,
				AssetID:  assetIDs[asset],
			})

			if err != nil {
				return c.Status(500).SendString(constants.ErrorS000)
			}
		}
	}

	// Create transaction record
	transactionID, err := qtx.CreateTransaction(c.Context(), db.CreateTransactionParams{
		SendingWalletID:   c.Locals("wid").(int64),
		ReceivingWalletID: recipientWalletID,
		CreatedAt:         time.Now(),
		Memo: sql.NullString{
			String: body.Memo,
			Valid:  body.Memo != "",
		},
	})

	var transactionAssets []db.CreateTransactionAssetsParams
	for asset, quantity := range body.Assets {
		transactionAssets = append(transactionAssets, db.CreateTransactionAssetsParams{
			TransactionID: transactionID,
			AssetID:       assetIDs[asset],
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
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	// TODO: Add in websocket support

	return c.Status(201).SendString("Transaction created")
}

func getTransactions(c *fiber.Ctx) error {
	// TODO: Add pagination, so an offset is needed
	query := struct {
		Limit int32 `query:"limit" validate:"min=0,max=100"`
	}{}
	// Parse and validate params
	if c.QueryParser(&query) != nil {
		return c.Status(400).SendString(constants.ErrorG002)
	}
	if validate.Struct(query) != nil {
		return c.Status(400).SendString(constants.ErrorG002)
	}
	// Set default if not already
	if query.Limit == 0 {
		query.Limit = 10
	}

	// Retrieve their transactions, edger load the assets
	transactions, err := database.Q.GetTransactions(c.Context(), db.GetTransactionsParams{
		SendingWalletID: c.Locals("wid").(int64),
		Limit:           query.Limit,
	})
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	var transactionIDs []int64
	for _, transaction := range transactions {
		transactionIDs = append(transactionIDs, transaction.ID)
	}

	transactionAssets, err := database.Q.GetTransactionAssetsByTransactionIds(c.Context(), transactionIDs)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	var assetIDs []int64
	for _, transactionAsset := range transactionAssets {
		assetIDs = append(assetIDs, transactionAsset.AssetID)
	}

	assets, err := database.Q.GetAssetsByIds(c.Context(), assetIDs)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Create a map of the assets using their id
	assetMap := make(map[int64]db.Asset)
	for _, asset := range assets {
		assetMap[asset.ID] = asset
	}

	type transactionAssetAPI struct {
		Quantity int64    `json:"quantity"`
		Asset    db.Asset `json:"asset"`
	}

	type transactionAPI struct {
		ID                int64                 `json:"id"`
		SendingWalletID   int64                 `json:"sending_wallet_id"`
		ReceivingWalletID int64                 `json:"receiving_wallet_id"`
		CreatedAt         time.Time             `json:"created_at"`
		Memo              string                `json:"memo"`
		Assets            []transactionAssetAPI `json:"assets"`
	}

	var transactionsAPI []transactionAPI
	for _, transaction := range transactions {
		var assets []transactionAssetAPI

		// TODO: Optimize this somehow
		for _, transactionAsset := range transactionAssets {
			if transactionAsset.TransactionID == transaction.ID {
				assets = append(assets, transactionAssetAPI{
					Quantity: transactionAsset.Quantity,
					Asset:    assetMap[transactionAsset.AssetID],
				})
			}
		}

		transactionsAPI = append(transactionsAPI, transactionAPI{
			ID:                transaction.ID,
			SendingWalletID:   transaction.SendingWalletID,
			ReceivingWalletID: transaction.ReceivingWalletID,
			CreatedAt:         transaction.CreatedAt,
			Memo:              transaction.Memo.String,
			Assets:            assets,
		})
	}

	if transactionsAPI == nil {
		transactionsAPI = []transactionAPI{}
	}

	return c.Status(200).JSON(transactionsAPI)
}

func deleteTransaction(c *fiber.Ctx) error {
	// Get walletid param
	params := struct {
		TransactionID int64 `params:"transactionid"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	rows, err := database.Q.DeleteTransactionById(c.Context(), db.DeleteTransactionByIdParams{
		SendingWalletID:   c.Locals("wid").(int64),
		ReceivingWalletID: c.Locals("wid").(int64),
		ID:                params.TransactionID,
	})
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorW004)
	}

	return c.Status(200).SendString("Transaction deleted")
}

func deleteTransactions(c *fiber.Ctx) error {
	// Get walletid param
	body := struct {
		TransactionIDs []int64 `json:"transaction_ids" validate:"required"`
	}{}
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	rows, err := database.Q.DeleteTransactionsById(c.Context(), db.DeleteTransactionsByIdParams{
		SendingWalletID:   c.Locals("wid").(int64),
		ReceivingWalletID: c.Locals("wid").(int64),
		Column3:           body.TransactionIDs,
	})

	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorW004)
	}

	return c.Status(200).SendString("Transactions deleted")
}

func postUserToWallet(c *fiber.Ctx) error {
	body := struct {
		Username string `json:"username" validate:"required,min=2,max=32,alphanum"`
	}{}

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Check if requester owns the wallet
	count, err := database.Q.CountWalletsByIdAndUserId(c.Context(), db.CountWalletsByIdAndUserIdParams{
		ID:     c.Locals("wid").(int64),
		UserID: c.Locals("uid").(int64),
	})
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if count == 0 {
		return c.Status(400).SendString(constants.ErrorU002)
	}

	// Get the username's user id
	userID, err := database.Q.GetUserIdByUsername(c.Context(), body.Username)
	if err != nil {
		return c.Status(404).SendString(constants.ErrorU003)
	}

	// Can't assign themself to their own wallet
	if userID == c.Locals("uid").(int64) {
		return c.Status(400).SendString(constants.ErrorW003)
	}

	// Assign the user to the wallet
	err = database.Q.CreateWalletUser(c.Context(), db.CreateWalletUserParams{
		WalletID: c.Locals("wid").(int64),
		UserID:   userID,
	})
	// TODO: Better error here, incase it's not a duplicate key
	if err != nil {
		return c.Status(400).SendString(constants.ErrorW003)
	}

	return c.Status(200).SendString("User assigned to wallet")
}

func deleteUserFromWallet(c *fiber.Ctx) error {
	// Get walletid param
	params := struct {
		UserID int64 `params:"userid"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	// Check if requester owns the wallet
	count, err := database.Q.CountWalletsByIdAndUserId(c.Context(), db.CountWalletsByIdAndUserIdParams{
		ID:     c.Locals("wid").(int64),
		UserID: c.Locals("uid").(int64),
	})
	if count == 0 {
		return c.Status(400).SendString(constants.ErrorU002)
	}

	// Revoke the user's access to the wallet
	err = database.Q.DeleteWalletUser(c.Context(), db.DeleteWalletUserParams{
		WalletID: c.Locals("wid").(int64),
		UserID:   params.UserID,
	})
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Delete any sessions that were using that wallet
	err = database.Q.DeleteUserSessionByUserIdAndWalletId(c.Context(), db.DeleteUserSessionByUserIdAndWalletIdParams{
		UserID:   params.UserID,
		WalletID: c.Locals("wid").(int64),
	})
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("User removed from wallet")
}

func getAssignedUsers(c *fiber.Ctx) error {
	assignedUsers, err := database.Q.GetAssignedUsersByWalletId(c.Context(), c.Locals("wid").(int64))
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if assignedUsers == nil {
		assignedUsers = []db.GetAssignedUsersByWalletIdRow{}
	}

	return c.Status(200).JSON(assignedUsers)
}

func postWalletSession(c *fiber.Ctx) error {
	body := struct {
		Name string `json:"name"`
	}{}

	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Add wallet session into db
	walletSessionID, err := database.Q.CreateWalletSession(c.Context(), db.CreateWalletSessionParams{
		WalletID: c.Locals("wid").(int64),
		Name: sql.NullString{
			String: body.Name,
			Valid:  body.Name != "",
		},
		UsedAt: time.Now(),
	})
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Generate the JWT
	claims := &auth.WalletJWT{
		SessionID: walletSessionID,
		WalletID:  c.Locals("wid").(int64),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString(tools.EnvVars.JwtSecret)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(201).JSON(fiber.Map{
		"token":   jwtString,
		"message": "Session created, token attached",
	})
}

func refreshWalletSession(c *fiber.Ctx) error {
	body := struct {
		Token string `json:"token"`
	}{}

	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	claims := &auth.WalletJWT{}
	token, err := jwt.ParseWithClaims(body.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return tools.EnvVars.JwtSecret, nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		// TODO: Better error than this
		return c.Status(400).SendString(constants.ErrorG000)
	}

	rows, err := database.Q.UpdateWalletSessionUsedAt(c.Context(), db.UpdateWalletSessionUsedAtParams{
		UsedAt: time.Now(),
		ID:     claims.SessionID,
	})

	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if rows == 0 {
		return c.Status(400).SendString(constants.ErrorW005)
	}

	// Generate the JWT
	claims = &auth.WalletJWT{
		SessionID: claims.SessionID,
		WalletID:  claims.WalletID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString(tools.EnvVars.JwtSecret)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(201).JSON(fiber.Map{
		"token":   jwtString,
		"message": "Session refreshed, token attached",
	})
}

func getWalletSessions(c *fiber.Ctx) error {
	walletSessions, err := database.Q.GetWalletSessionsByWalletId(c.Context(), c.Locals("wid").(int64))
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if walletSessions == nil {
		walletSessions = []db.WalletSession{}
	}

	return c.Status(200).JSON(walletSessions)
}

func deleteWalletSession(c *fiber.Ctx) error {
	// Get walletid param
	params := struct {
		SessionID int64 `params:"sessionid"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	rows, err := database.Q.DeleteWalletSession(c.Context(), params.SessionID)
	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorW005)
	}

	return c.Status(200).SendString("Session deleted")
}

func deleteWalletSessions(c *fiber.Ctx) error {
	err := database.Q.DeleteWalletSessionsByWalletId(c.Context(), c.Locals("wid").(int64))

	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("All sessions deleted")
}