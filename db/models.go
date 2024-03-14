// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type TransferStatus string

const (
	TransferStatusOpen     TransferStatus = "open"
	TransferStatusDeclined TransferStatus = "declined"
	TransferStatusApproved TransferStatus = "approved"
	TransferStatusCleared  TransferStatus = "cleared"
)

func (e *TransferStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransferStatus(s)
	case string:
		*e = TransferStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TransferStatus: %T", src)
	}
	return nil
}

type NullTransferStatus struct {
	TransferStatus TransferStatus `json:"transfer_status"`
	Valid          bool           `json:"valid"` // Valid is true if TransferStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTransferStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TransferStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TransferStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTransferStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TransferStatus), nil
}

type Asset struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type Transaction struct {
	ID                int64       `json:"id"`
	SendingWalletID   int64       `json:"sending_wallet_id"`
	ReceivingWalletID int64       `json:"receiving_wallet_id"`
	CreatedAt         time.Time   `json:"created_at"`
	Memo              pgtype.Text `json:"memo"`
}

type TransactionAsset struct {
	ID            int64 `json:"id"`
	TransactionID int64 `json:"transaction_id"`
	AssetID       int64 `json:"asset_id"`
	Quantity      int64 `json:"quantity"`
}

type Transfer struct {
	ID                   int64          `json:"id"`
	CreatedAt            time.Time      `json:"created_at"`
	Status               TransferStatus `json:"status"`
	SendingWarehouseID   int64          `json:"sending_warehouse_id"`
	ReceivingWarehouseID int64          `json:"receiving_warehouse_id"`
}

type TransferAsset struct {
	ID         int64 `json:"id"`
	TransferID int64 `json:"transfer_id"`
	AssetID    int64 `json:"asset_id"`
	Quantity   int64 `json:"quantity"`
}

type User struct {
	ID                  int64       `json:"id"`
	Username            string      `json:"username"`
	Password            string      `json:"password"`
	CreatedAt           time.Time   `json:"created_at"`
	WalletID            pgtype.Int8 `json:"wallet_id"`
	CanCreateWarehouses bool        `json:"can_create_warehouses"`
}

type UserSession struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	WalletID int64  `json:"wallet_id"`
	Key      string `json:"key"`
}

type Wallet struct {
	ID      int64       `json:"id"`
	Address string      `json:"address"`
	UserID  int64       `json:"user_id"`
	Webhook pgtype.Text `json:"webhook"`
}

type WalletAsset struct {
	ID       int64 `json:"id"`
	WalletID int64 `json:"wallet_id"`
	AssetID  int64 `json:"asset_id"`
	Quantity int64 `json:"quantity"`
}

type WalletSession struct {
	ID       int64       `json:"id"`
	WalletID int64       `json:"wallet_id"`
	Name     pgtype.Text `json:"name"`
	Key      string      `json:"key"`
}

type WalletUser struct {
	ID       int64 `json:"id"`
	WalletID int64 `json:"wallet_id"`
	UserID   int64 `json:"user_id"`
}

type Warehouse struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	UserID          int64          `json:"user_id"`
	Location        interface{}    `json:"location"`
	Liability       int64          `json:"liability"`
	Collateral      int64          `json:"collateral"`
	CollateralRatio pgtype.Numeric `json:"collateral_ratio"`
}

type WarehouseAsset struct {
	ID          int64 `json:"id"`
	WarehouseID int64 `json:"warehouse_id"`
	AssetID     int64 `json:"asset_id"`
	Quantity    int64 `json:"quantity"`
}

type WarehouseWorker struct {
	ID          int64 `json:"id"`
	WarehouseID int64 `json:"warehouse_id"`
	UserID      int64 `json:"user_id"`
}
