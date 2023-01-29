package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/middlewares"
	"github.com/stelofinance/api/routes"
	"github.com/stelofinance/api/tools"
)

func main() {
	// Load in env
	if tools.LoadEnv() != nil {
		log.Fatal("Failed to load environment")
	}

	// Connect to db
	if database.ConnectDb() != nil {
		log.Fatal("Failed to connect to database")
	}

	app := fiber.New()

	// Log request
	app.Use(logger.New(logger.Config{
		Format: "${status} - ${method} ${path}\n",
	}))

	// Setup routes
	routes.CentrifugoRouter(app.Group("/centrifugo"))
	routes.UsersRouter(app.Group("/users"))
	routes.UserRouter(app.Group("/user", auth.New(auth.User)))
	routes.WalletRouter(app.Group("/wallet"))
	routes.WalletsRouter(app.Group("/wallets", auth.New(auth.Admin)))
	routes.AssetsRouter(app.Group("/assets", auth.New(auth.Admin)))

	// Run app
	log.Fatal(app.Listen(":8080"))
}
