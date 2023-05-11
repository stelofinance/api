package routes

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/pusher"
)

func postPusherAuth(c *fiber.Ctx) error {
	type requestBody struct {
		SocketId    string `form:"socket_id" validate:"required"`
		ChannelName string `form:"channel_name" validate:"required"`
	}
	var body requestBody

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	channelParts := strings.Split(body.ChannelName, "@")

	if len(channelParts) != 2 {
		return c.Status(400).SendString(constants.ErrorP000)
	}

	if channelParts[0] != "private-wallet" {
		return c.Status(400).SendString(constants.ErrorP001)
	}

	if channelParts[1] != fmt.Sprint(c.Locals("wid")) {
		return c.Status(403).SendString(constants.ErrorW006)
	}

	response, err := pusher.PusherClient.AuthorizePrivateChannel(c.Body())

	if err != nil {
		return c.Status(500).SendString(constants.ErrorS000)
	}

	c.Set("Content-Type", "application/json")
	return c.Status(200).Send(response)
}
