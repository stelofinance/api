package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgtype"
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
		log.Printf("Error getting assets: {%v}", err.Error())
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
		Recipient string           `json:"recipient" validate:"required"`
		Type      uint8            `json:"type" validate:"lte=2"`
		Memo      string           `json:"memo" validate:"max=64"`
		Assets    map[string]int64 `json:"assets" validate:"gt=0,dive,gt=0"`
	}
	var body requestBody

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if body.Type >= 3 {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	tx, err := database.DB.Begin(c.Context())
	defer tx.Rollback(c.Context())
	if err != nil {
		log.Printf("Error creating db transaction: {%v}", err.Error())
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
		log.Printf("Error getting assets id & name: {%v}", err.Error())
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
			log.Printf("Error subtracting quantity from asset: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			return c.Status(400).SendString(constants.ErrorI001)
		}
	}

	// Now put assets in other wallet, first get their wallet ID though
	// Create the wallet asset record if not already
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
		if resp.StatusCode != 200 {
			return c.Status(400).SendString(constants.ErrorW010)
		}
		resp.Body.Close()
	}

	for asset, quantity := range body.Assets {
		rows, err := qtx.AddWalletAssetQuantity(c.Context(), db.AddWalletAssetQuantityParams{
			Quantity: quantity,
			WalletID: recipientWalletID,
			AssetID:  assetIDs[asset],
		})

		if err != nil {
			log.Printf("Error adding asset quantity: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			err := qtx.CreateWalletAsset(c.Context(), db.CreateWalletAssetParams{
				Quantity: quantity,
				WalletID: recipientWalletID,
				AssetID:  assetIDs[asset],
			})

			if err != nil {
				log.Printf("Error creating wallet asset: {%v}", err.Error())
				return c.Status(500).SendString(constants.ErrorS000)
			}
		}
	}

	// Create transaction record
	transactionID, err := qtx.CreateTransaction(c.Context(), db.CreateTransactionParams{
		SendingWalletID:   c.Locals("wid").(int64),
		ReceivingWalletID: recipientWalletID,
		CreatedAt:         time.Now(),
		Memo: pgtype.Text{
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
		log.Printf("Error inserting transaction assets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	// Send to centrifugo cannel
	go func(recipient int64, sender int64, assets map[string]int64) {
		jsonBody, err := json.Marshal(fiber.Map{
			"method": "publish",
			"params": fiber.Map{
				"channel": "wallet:transactions#" + strconv.FormatInt(recipientWalletID, 10),
				"data": fiber.Map{
					"sender": sender,
					"assets": assets,
				},
			},
		})
		if err != nil {
			log.Printf("Error creating json body for centrifugo: {%v}", err.Error())
			return
		}

		req, err := http.NewRequest("POST", tools.EnvVars.CentrifugoAddr, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Printf("Error creating req to centrifugo: {%v}", err.Error())
			return
		}

		req.Header.Set("Authorization", "apikey "+tools.EnvVars.CentrifugoApiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request to centrifugo: {%v}", err.Error())
			return
		}
		resp.Body.Close()
	}(recipientWalletID, c.Locals("wid").(int64), body.Assets)

	return c.Status(201).SendString("Transaction created")
}

func getTransactions(c *fiber.Ctx) error {
	// TODO: Add pagination, so an offset is needed
	query := struct {
		Limit  int32 `query:"limit" validate:"min=0,max=100"`
		Offset int32 `query:"offset" validate:"min=0,max=1000"`
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
		Offset:          query.Offset,
	})
	if err != nil {
		log.Printf("Error getting transactions: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	var transactionIDs []int64
	for _, transaction := range transactions {
		transactionIDs = append(transactionIDs, transaction.ID)
	}

	transactionAssets, err := database.Q.GetTransactionAssetsByTransactionIds(c.Context(), transactionIDs)
	if err != nil {
		log.Printf("Error getting tx assets by tx ids: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	var assetIDs []int64
	for _, transactionAsset := range transactionAssets {
		assetIDs = append(assetIDs, transactionAsset.AssetID)
	}

	assets, err := database.Q.GetAssetsByIds(c.Context(), assetIDs)
	if err != nil {
		log.Printf("Error getting assets: {%v}", err.Error())
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
		log.Printf("Error deleting transaction: {%v}", err.Error())
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
		log.Printf("Error deleteing transactions: {%v}", err.Error())
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
		log.Printf("Error counting wallets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if count == 0 {
		return c.Status(400).SendString(constants.ErrorW006)
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
		return c.Status(400).SendString(constants.ErrorW006)
	}

	// Revoke the user's access to the wallet
	err = database.Q.DeleteWalletUser(c.Context(), db.DeleteWalletUserParams{
		WalletID: c.Locals("wid").(int64),
		UserID:   params.UserID,
	})
	if err != nil {
		log.Printf("Error deleting wallet user: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Delete any sessions that were using that wallet
	err = database.Q.DeleteUserSessionByUserIdAndWalletId(c.Context(), db.DeleteUserSessionByUserIdAndWalletIdParams{
		UserID:   params.UserID,
		WalletID: c.Locals("wid").(int64),
	})
	if err != nil {
		log.Printf("Error deleting user session: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("User removed from wallet")
}

func getAssignedUsers(c *fiber.Ctx) error {
	assignedUsers, err := database.Q.GetAssignedUsersByWalletId(c.Context(), c.Locals("wid").(int64))
	if err != nil {
		log.Printf("Error getting assinged users: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if assignedUsers == nil {
		assignedUsers = []db.GetAssignedUsersByWalletIdRow{}
	}

	return c.Status(200).JSON(assignedUsers)
}

func putWalletOwner(c *fiber.Ctx) error {
	body := struct {
		UserID int64 `json:"user_id" validate:"required"`
	}{}
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Switch their user_wallet entry for requesters
	// user_id, then change the wallet's user_id to new owner
	tx, err := database.DB.Begin(c.Context())
	defer tx.Rollback(c.Context())
	if err != nil {
		log.Printf("Error creating db transaction: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	rows, err := qtx.UpdateWalletUserUserID(c.Context(), db.UpdateWalletUserUserIDParams{
		UserID:   c.Locals("uid").(int64),
		WalletID: c.Locals("wid").(int64),
		UserID_2: body.UserID,
	})

	if rows == 0 {
		return c.Status(400).SendString(constants.ErrorW007)
	}
	if err != nil {
		log.Printf("Error updating wallet user: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	err = qtx.UpdateWalletUserID(c.Context(), db.UpdateWalletUserIDParams{
		UserID:   body.UserID,
		ID:       c.Locals("wid").(int64),
		UserID_2: c.Locals("uid").(int64),
	})
	if err != nil {
		log.Printf("Error updating wallet user id: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	tx.Commit(c.Context())

	return c.Status(200).SendString("New owner assigned")
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
		Name: pgtype.Text{
			String: body.Name,
			Valid:  body.Name != "",
		},
		UsedAt: time.Now(),
	})
	if err != nil {
		log.Printf("Error creating wallet session: {%v}", err.Error())
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
		log.Printf("Error creating JWT: {%v}", err.Error())
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
		log.Printf("Error updating wallet session used at: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorW005)
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
		log.Printf("Error creating JWT: {%v}", err.Error())
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
		log.Printf("Error getting wallet sessions: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	type walletSessionAPI struct {
		ID       int64     `json:"id"`
		WalletID int64     `json:"wallet_id"`
		Name     string    `json:"name"`
		UsedAt   time.Time `json:"used_at"`
	}

	walletSessionsAPI := []walletSessionAPI{}
	for _, walletSession := range walletSessions {
		walletSessionsAPI = append(walletSessionsAPI, walletSessionAPI{
			ID:       walletSession.ID,
			WalletID: walletSession.WalletID,
			Name:     walletSession.Name.String,
			UsedAt:   walletSession.UsedAt,
		})
	}

	return c.Status(200).JSON(walletSessionsAPI)
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
		log.Printf("Error deleting wallet session: {%v}", err.Error())
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
		log.Printf("Error deleting wallet sessions: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("All sessions deleted")
}

func deleteWallet(c *fiber.Ctx) error {
	tx, err := database.DB.Begin(c.Context())
	defer tx.Rollback(c.Context())
	if err != nil {
		log.Printf("Error creating db transaction: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	qtx := database.Q.WithTx(tx)

	// Check if requester owns the wallet
	count, err := database.Q.CountWalletsByIdAndUserId(c.Context(), db.CountWalletsByIdAndUserIdParams{
		ID:     c.Locals("wid").(int64),
		UserID: c.Locals("uid").(int64),
	})
	if err != nil {
		log.Printf("Error counting wallets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if count == 0 {
		return c.Status(400).SendString(constants.ErrorW006)
	}

	// Check if wallet is primary
	user, err := qtx.GetUserById(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Printf("Error getting user: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if !user.WalletID.Valid {
		log.Printf("Error, user created without wallet: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if user.WalletID.Int64 == c.Locals("wid").(int64) {
		return c.Status(400).SendString(constants.ErrorW008)
	}

	// Get all the assets currently in the wallet
	assets, err := qtx.GetWalletAssets(c.Context(), c.Locals("wid").(int64))
	if err != nil {
		log.Printf("Error getting wallet assets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Delete all the assets in the wallet
	_, err = qtx.DeleteWalletAssets(c.Context(), c.Locals("wid").(int64))
	if err != nil {
		log.Printf("Error deleting wallet assets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Add/create all the assets in the primary wallet
	for _, asset := range assets {
		rows, err := qtx.AddWalletAssetQuantity(c.Context(), db.AddWalletAssetQuantityParams{
			Quantity: asset.Quantity,
			WalletID: user.WalletID.Int64,
			AssetID:  asset.AssetID,
		})

		if err != nil {
			log.Printf("Error adding wallet asset quantity: {%v}", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		} else if rows == 0 {
			err := qtx.CreateWalletAsset(c.Context(), db.CreateWalletAssetParams{
				Quantity: asset.Quantity,
				WalletID: user.WalletID.Int64,
				AssetID:  asset.AssetID,
			})

			if err != nil {
				log.Printf("Error creating wallet asset: {%v}", err.Error())
				return c.Status(500).SendString(constants.ErrorS000)
			}
		}
	}

	// Create transaction record
	transactionID, err := qtx.CreateTransaction(c.Context(), db.CreateTransactionParams{
		SendingWalletID:   user.WalletID.Int64,
		ReceivingWalletID: user.WalletID.Int64,
		CreatedAt:         time.Now(),
		Memo: pgtype.Text{
			String: "funds from deleted wallet",
			Valid:  true,
		},
	})

	var transactionAssets []db.CreateTransactionAssetsParams
	for _, asset := range assets {
		transactionAssets = append(transactionAssets, db.CreateTransactionAssetsParams{
			TransactionID: transactionID,
			AssetID:       asset.AssetID,
			Quantity:      asset.Quantity,
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
		log.Printf("Error creating transaction assets: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// TODO: Add in websocket support

	// Update their session
	// UNSAFE: Type assertion could panic
	err = qtx.UpdateUserSessionWallet(c.Context(), db.UpdateUserSessionWalletParams{
		ID:       c.Locals("sid").(int64),
		WalletID: user.WalletID.Int64,
	})
	if err != nil {
		log.Printf("Error updating user session: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Delete the wallet
	_, err = qtx.DeleteWallet(c.Context(), c.Locals("wid").(int64))
	if err != nil {
		log.Printf("Error deleting wallet: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	// Create new JWT for their cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.UserJWT{
		UserID:    c.Locals("uid").(int64), // UNSAFE: Type assertion could panic
		SessionID: c.Locals("sid").(int64), // UNSAFE: Type assertion could panic
		WalletID:  user.WalletID.Int64,
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

	tx.Commit(c.Context())

	return c.Status(200).SendString("Wallet deleted")
}

func putWalletWebhook(c *fiber.Ctx) error {
	body := struct {
		Webhook string `json:"webhook" validate:"required,max=128,url"`
	}{}
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Check if wallet is primary
	user, err := database.Q.GetUserById(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Printf("Error getting user: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if !user.WalletID.Valid {
		log.Printf("Error, user created without wallet: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}
	if user.WalletID.Int64 == c.Locals("wid").(int64) {
		return c.Status(400).SendString(constants.ErrorW009)
	}

	// Update webhook
	err = database.Q.UpdateWalletWebhook(c.Context(), db.UpdateWalletWebhookParams{
		ID: c.Locals("wid").(int64),
		Webhook: pgtype.Text{
			String: body.Webhook,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("Error updating wallet webhook: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Wallet webhook updated")
}

func deleteWalletWebhook(c *fiber.Ctx) error {
	err := database.Q.DeleteWalletWebhook(c.Context(), c.Locals("wid").(int64))
	if err != nil {
		log.Printf("Error deleting wallet webhook: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Webhook removed from wallet")
}