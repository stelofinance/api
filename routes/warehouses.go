package routes

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
)

func postWarehouse(c *fiber.Ctx) error {
	type requestBody struct {
		Name        string   `json:"name" validate:"required"`
		Coordinates [2]int64 `json:"coordinates"`
	}
	var body requestBody

	// Parse and validate body
	if c.BodyParser(&body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}
	if validate.Struct(body) != nil {
		return c.Status(400).SendString(constants.ErrorG000)
	}

	// Check they have permission to create the warehouse
	canCreateWarehouses, err := database.Q.GetUserCanCreateWarehouses(c.Context(), c.Locals("uid").(int64))
	if err != nil {
		log.Println("Error getting can create warehouses permission", err.Error())
		return c.Status(500).SendString(constants.ErrorS000)
	}

	if !canCreateWarehouses {
		return c.Status(403).SendString(constants.ErrorH000)
	}

	// Create the warehouse
	err = database.Q.InsertWarehouse(
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

	return c.Status(201).SendString("Warehouse created")
}
