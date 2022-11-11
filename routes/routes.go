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
}

func UserRouter(app fiber.Router) {
	app.Put("/username", putUsername)
	app.Put("/password", putPassword)
	app.Put("/wallet", putWallet)
	app.Get("/wallets", getWallets)
	app.Get("/assigned_wallets", getAssignedWallets)
	app.Put("/session/wallet", putActiveWallet)
	app.Delete("/session", deleteSession)
	app.Delete("/session/:sessionid", deleteSessionById)
	app.Get("/sessions", getSessions)
	app.Delete("/sessions", deleteSessions)
}

func WalletRouter(app fiber.Router) {
	app.Get("/assets", auth.New(auth.Wallet), getAssets)
	app.Post("/transactions", auth.New(auth.Wallet), postTransaction)
	app.Get("/transactions", auth.New(auth.Wallet), getTransactions)
	app.Delete("/transaction/:transactionid", auth.New(auth.Wallet), deleteTransaction)
	app.Delete("/transactions", auth.New(auth.Wallet), deleteTransactions)
	app.Post("/user", auth.New(auth.User), postUserToWallet)
	app.Get("/users", auth.New(auth.User), getAssignedUsers)
	app.Delete("/user/:userid", auth.New(auth.User), deleteUserFromWallet)
	app.Post("/sessions", auth.New(auth.User), postWalletSession)
	app.Get("/sessions", auth.New(auth.User), getWalletSessions)
	app.Post("/sessions/token", refreshWalletSession)
	app.Delete("/sessions/:sessionid", auth.New(auth.User), deleteWalletSession)
	app.Delete("/sessions", auth.New(auth.User), deleteWalletSessions)
}

func WalletsRouter(app fiber.Router) {
	app.Post("/", auth.New(auth.User), postWallet)
	app.Post("/:walletid/assets", auth.New(auth.Admin), postAssetToWallet)
	app.Delete("/:walletid/assets/:assetid", auth.New(auth.Admin), deleteAssetFromWallet)
}

func AssetsRouter(app fiber.Router) {
	app.Post("/", postAsset)
	app.Put("/:id/value", putAssetValue)
	app.Put("/:id/name", putAssetName)
	app.Delete("/:id", deleteAsset)
}
