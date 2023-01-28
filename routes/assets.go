package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
)

func postAsset(c *fiber.Ctx) error {
	type requestBody struct {
		Name  string `json:"name" validate:"required"`
		Value int64  `json:"value" validate:"required"`
	}
	var body requestBody

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	err := database.Q.CreateAsset(c.Context(), db.CreateAssetParams{
		Name:  body.Name,
		Value: body.Value,
	})

	if err != nil {
		log.Printf("Error creating asset: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(201).SendString("Asset created")
}

func putAssetValue(c *fiber.Ctx) error {
	// Get id param
	params := struct {
		ID int64 `params:"id"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	type requestBody struct {
		Value int64 `json:"value"`
	}
	var body requestBody

	// Parse body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	rows, err := database.Q.UpdateAssetValue(c.Context(), db.UpdateAssetValueParams{
		Value: body.Value,
		ID:    params.ID,
	})

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorI000)
	}
	if err != nil {
		log.Printf("Error updating asset value: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Asset value updated")
}

func putAssetName(c *fiber.Ctx) error {
	// Get id param
	params := struct {
		ID int64 `params:"id"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	type requestBody struct {
		Name string `json:"name" validate:"required"`
	}
	var body requestBody

	// Parse body and validate it
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	rows, err := database.Q.UpdateAssetName(c.Context(), db.UpdateAssetNameParams{
		Name: body.Name,
		ID:   params.ID,
	})

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorI000)
	}
	if err != nil {
		log.Printf("Error updating asset name: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Asset name updated")
}

func deleteAsset(c *fiber.Ctx) error {
	// Get id param
	params := struct {
		ID int64 `params:"id"`
	}{}
	if c.ParamsParser(&params) != nil {
		return c.Status(400).SendString(constants.ErrorG001)
	}

	rows, err := database.Q.DeleteAsset(c.Context(), params.ID)

	if rows == 0 {
		return c.Status(404).SendString(constants.ErrorI000)
	}
	if err != nil {
		log.Printf("Error deleting asset: {%v}", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	return c.Status(200).SendString("Asset successfully removed")
}