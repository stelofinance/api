// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"database/sql"
	"time"
)

type Asset struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type Transaction struct {
	ID                int64          `json:"id"`
	SendingWalletID   int64          `json:"sending_wallet_id"`
	ReceivingWalletID int64          `json:"receiving_wallet_id"`
	CreatedAt         time.Time      `json:"created_at"`
	Memo              sql.NullString `json:"memo"`
}

type TransactionAsset struct {
	ID            int64 `json:"id"`
	TransactionID int64 `json:"transaction_id"`
	AssetID       int64 `json:"asset_id"`
	Quantity      int64 `json:"quantity"`
}

type User struct {
	ID        int64         `json:"id"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	CreatedAt time.Time     `json:"created_at"`
	WalletID  sql.NullInt64 `json:"wallet_id"`
}

type UserSession struct {
	ID       int64     `json:"id"`
	UsedAt   time.Time `json:"used_at"`
	UserID   int64     `json:"user_id"`
	WalletID int64     `json:"wallet_id"`
}

type Wallet struct {
	ID      int64          `json:"id"`
	Address string         `json:"address"`
	UserID  int64          `json:"user_id"`
	Webhook sql.NullString `json:"webhook"`
}

type WalletAsset struct {
	ID       int64 `json:"id"`
	WalletID int64 `json:"wallet_id"`
	AssetID  int64 `json:"asset_id"`
	Quantity int64 `json:"quantity"`
}

type WalletSession struct {
	ID       int64          `json:"id"`
	WalletID int64          `json:"wallet_id"`
	Name     sql.NullString `json:"name"`
	UsedAt   time.Time      `json:"used_at"`
}

type WalletUser struct {
	ID       int64 `json:"id"`
	WalletID int64 `json:"wallet_id"`
	UserID   int64 `json:"user_id"`
}
