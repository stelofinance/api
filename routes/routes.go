package routes

import "github.com/gofiber/fiber/v2"

func UsersRouter(app fiber.Router) {
	app.Post("/", postUser)
	app.Post("/:username/sessions", postSession)
}
