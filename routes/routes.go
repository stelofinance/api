package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/middlewares"
)

func UsersRouter(app fiber.Router) {
	app.Post("/", postUser)
	app.Post("/:username/sessions", postSession)
	app.Post("/test", middlewares.AuthSession, testCookie)
}
