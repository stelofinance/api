package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/middlewares"
)

var validate = validator.New()

func UsersRouter(app fiber.Router) {
	app.Post("/", auth.New(auth.Guest), postUser)
	app.Post("/:username/sessions", auth.New(auth.Guest), postSession)
	app.Put("/:username/can_create_warehouses", auth.New(auth.Admin), putCanCreateWarehouses)
}

func PusherRouter(app fiber.Router) {
	app.Post("/auth", auth.New(auth.Wallet), postPusherAuth)
}

func UserRouter(app fiber.Router) {
	app.Put("/username", putUsername)
	app.Put("/password", putPassword)
	app.Put("/wallet", putWallet)
	app.Get("/wallets", getWallets)
	app.Post("/wallets", postWallet)
	app.Get("/assigned_wallets", getAssignedWallets)
	app.Put("/session/wallet", putActiveWallet)
	app.Get("/session", getSession)
	app.Delete("/session", deleteSession)
	app.Delete("/sessions/:sessionid", deleteSessionById)
	app.Get("/sessions", getSessions)
	app.Delete("/sessions", deleteSessions)
}

func WalletRouter(app fiber.Router) {
	app.Get("/assets", auth.New(auth.Wallet), getAssets)
	app.Post("/transactions", auth.New(auth.Wallet), postTransaction)
	app.Get("/transactions", auth.New(auth.Wallet), getTransactions)
	app.Post("/users", auth.New(auth.User), postUserToWallet)
	app.Get("/users", auth.New(auth.User), getAssignedUsers)
	app.Delete("/users/:userid", auth.New(auth.User), deleteUserFromWallet)
	app.Post("/sessions", auth.New(auth.User), postWalletSession)
	app.Get("/sessions", auth.New(auth.User), getWalletSessions)
	app.Delete("/sessions/:sessionid", auth.New(auth.User), deleteWalletSession)
	app.Delete("/sessions", auth.New(auth.User), deleteWalletSessions)
	app.Put("/owner", auth.New(auth.User), putWalletOwner)
	app.Delete("/", auth.New(auth.User), deleteWallet)
	app.Put("/webhook", auth.New(auth.User), putWalletWebhook)
	app.Delete("/webhook", auth.New(auth.User), deleteWalletWebhook)
}

func WalletsRouter(app fiber.Router) {
	app.Post("/:walletid/assets", postAssetToWallet)
	app.Delete("/:walletid/assets/:assetid", deleteAssetFromWallet)
}

func AssetsRouter(app fiber.Router) {
	app.Post("/", postAsset)
	app.Put("/:id/value", putAssetValue)
	app.Put("/:id/name", putAssetName)
	app.Delete("/:id", deleteAsset)
}

func WarehousesRouter(app fiber.Router) {
	app.Post("/", postWarehouse)
	app.Put("/:warehouseid/collateral", auth.NewWarehouse(auth.Owner), putCollateral)
	app.Get("/:warehouseid/collateral", auth.NewWarehouse(auth.Owner), getWarehouseCollateral)
	app.Put("/:warehouseid/owner", auth.NewWarehouse(auth.Owner), putWarehouseOwner)
	app.Post("/:warehouseid/workers", auth.NewWarehouse(auth.Owner), postWarehouseWorker)
	app.Delete("/:warehouseid/workers/:workerid", auth.NewWarehouse(auth.Owner), deleteWarehouseWorker)
	app.Get("/:warehouseid/workers", auth.NewWarehouse(auth.Worker), getWarehouseWorkers)
	app.Post("/:warehouseid/assets", auth.NewWarehouse(auth.Worker), postWarehouseAssets)
	app.Delete("/:warehouseid/assets", auth.NewWarehouse(auth.Worker), deleteWarehouseAssets)
	app.Get("/:warehouseid/assets", auth.NewWarehouse(auth.Worker), getWarehouseAssets)
}
