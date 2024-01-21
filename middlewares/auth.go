package auth

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/stelofinance/api/constants"
	"github.com/stelofinance/api/database"
	"github.com/stelofinance/api/db"
	"github.com/stelofinance/api/tools"
)

type Role struct {
	slug string
}

var (
	Unknown = Role{""}
	Guest   = Role{"stlg"}
	User    = Role{"stlu"}
	Wallet  = Role{"stlw"}
	Admin   = Role{"stla"}
)

func (r Role) String() string {
	return r.slug
}

func New(role Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if role == Guest {
			if c.GetReqHeaders()["Authorization"] == "" &&
				c.Cookies("stelo_token") == "" {
				return c.Next()
			}
			return c.Status(400).SendString(constants.ErrorA002)
		}

		var fullToken string

		header := c.GetReqHeaders()["Authorization"]
		cookie := c.Cookies("stelo_token")
		if header == "" && cookie == "" {
			return c.Status(401).SendString(constants.ErrorA000)
		} else if header != "" {
			fullToken = header
		} else if cookie != "" {
			fullToken = cookie
		} else {
			log.Fatal("Logic error in middleware, cookie & header empty")
			return c.Status(500).SendString(constants.ErrorS000)
		}

		tokenArray := strings.Split(fullToken, "_")

		if len(tokenArray) != 2 {
			return c.Status(400).SendString(constants.ErrorA001)
		}

		key := tokenArray[1]

		var tokenRole Role
		switch tokenArray[0] {
		case User.slug:
			tokenRole = User
		case Wallet.slug:
			tokenRole = Wallet
		case Admin.slug:
			tokenRole = Admin
		default:
			return c.Status(400).SendString(constants.ErrorA001)
		}

		if role == Admin {
			if tokenRole != Admin || key != tools.EnvVars.AdminKey {
				return c.Status(403).SendString(constants.ErrorA003)
			}
			return c.Next()
		} else if role == User {
			if tokenRole != User {
				return c.Status(400).SendString(constants.ErrorA001)
			}

			row, err := database.Q.GetUserSession(c.Context(), key)
			if err != nil {
				return c.Status(400).SendString(constants.ErrorA001)
			}

			c.Locals("uid", row.UserID)
			c.Locals("sid", row.ID)
			c.Locals("wid", row.WalletID)
			return c.Next()
		} else if role == Wallet {

			var walletId int64

			if tokenRole == User {
				userRow, err := database.Q.GetUserSession(c.Context(), key)
				if err != nil {
					return c.Status(400).SendString(constants.ErrorA001)
				}
				walletId = userRow.WalletID

			} else if tokenRole == Wallet {
				walletIdReq, err := database.Q.GetWalletSession(c.Context(), key)
				if err != nil {
					return c.Status(400).SendString(constants.ErrorA001)
				}
				walletId = walletIdReq

			} else {
				return c.Status(400).SendString(constants.ErrorA001)
			}

			// Removed, note left for this context:
			// Currently doesn't need this value, and if it does
			// then it breaks the user sessions be able to be used for
			// wallet requests
			// c.Locals("sid", row.ID)

			c.Locals("wid", walletId)
			return c.Next()
		}

		log.Fatal("Auth no Role matched")
		return c.Status(500).SendString(constants.ErrorS000)
	}
}

type WarehouseRole int

const (
	Worker WarehouseRole = iota
	Owner
)

// `warehouseid` must be in the path for this auth middleware
func NewWarehouse(role WarehouseRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		warehouseId, err := strconv.ParseInt(c.Params("warehouseid"), 10, 64)
		if err != nil {
			return c.Status(400).SendString(constants.ErrorG001)
		}

		if role == Worker {
			isWorker, err := database.Q.ExistsWarehouseWorker(c.Context(), db.ExistsWarehouseWorkerParams{
				WarehouseID: warehouseId,
				UserID:      c.Locals("uid").(int64),
			})
			if err != nil {
				log.Println("Error checking is new owner is a worker", err.Error())
				return c.Status(500).SendString(constants.ErrorS000)
			}
			if !isWorker {
				return c.Status(400).SendString(constants.ErrorH004)
			}

			return c.Next()
		}
		// Else it's owner
		userId, err := database.Q.GetWarehouseUserId(c.Context(), int64(warehouseId))
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(404).SendString(constants.ErrorH001)
			}
			log.Println("Failed to fetch user id from warehouse", err.Error())
			return c.Status(500).SendString(constants.ErrorS000)
		}

		if userId != c.Locals("uid").(int64) {
			return c.Status(403).SendString(constants.ErrorH002)
		}

		return c.Next()
	}
}
