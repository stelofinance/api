// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: wallet.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countWalletsByIdAndUserId = `-- name: CountWalletsByIdAndUserId :one
SELECT count(*) FROM wallet WHERE id = $1 AND user_id = $2
`

type CountWalletsByIdAndUserIdParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) CountWalletsByIdAndUserId(ctx context.Context, arg CountWalletsByIdAndUserIdParams) (int64, error) {
	row := q.db.QueryRow(ctx, countWalletsByIdAndUserId, arg.ID, arg.UserID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createWallet = `-- name: CreateWallet :exec
INSERT INTO wallet (address, user_id) VALUES ($1, $2)
`

type CreateWalletParams struct {
	Address string `json:"address"`
	UserID  int64  `json:"user_id"`
}

func (q *Queries) CreateWallet(ctx context.Context, arg CreateWalletParams) error {
	_, err := q.db.Exec(ctx, createWallet, arg.Address, arg.UserID)
	return err
}

const deleteWallet = `-- name: DeleteWallet :execrows
DELETE FROM wallet WHERE id = $1
`

func (q *Queries) DeleteWallet(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.Exec(ctx, deleteWallet, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteWalletWebhook = `-- name: DeleteWalletWebhook :exec
UPDATE wallet set webhook = NULL WHERE id = $1
`

func (q *Queries) DeleteWalletWebhook(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteWalletWebhook, id)
	return err
}

const getAssignedWalletsByUserId = `-- name: GetAssignedWalletsByUserId :many
SELECT wallet.id, wallet.address, wallet.user_id 
FROM wallet 
INNER JOIN wallet_user 
    ON wallet.id = wallet_user.wallet_id 
        AND wallet_user.user_id = $1
`

type GetAssignedWalletsByUserIdRow struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
	UserID  int64  `json:"user_id"`
}

func (q *Queries) GetAssignedWalletsByUserId(ctx context.Context, userID int64) ([]GetAssignedWalletsByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getAssignedWalletsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAssignedWalletsByUserIdRow
	for rows.Next() {
		var i GetAssignedWalletsByUserIdRow
		if err := rows.Scan(&i.ID, &i.Address, &i.UserID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWalletIdAndWebhookByAddress = `-- name: GetWalletIdAndWebhookByAddress :one
SELECT id, webhook FROM wallet WHERE address = $1
`

type GetWalletIdAndWebhookByAddressRow struct {
	ID      int64       `json:"id"`
	Webhook pgtype.Text `json:"webhook"`
}

func (q *Queries) GetWalletIdAndWebhookByAddress(ctx context.Context, address string) (GetWalletIdAndWebhookByAddressRow, error) {
	row := q.db.QueryRow(ctx, getWalletIdAndWebhookByAddress, address)
	var i GetWalletIdAndWebhookByAddressRow
	err := row.Scan(&i.ID, &i.Webhook)
	return i, err
}

const getWalletWebhook = `-- name: GetWalletWebhook :one
SELECT webhook FROM wallet WHERE id = $1
`

func (q *Queries) GetWalletWebhook(ctx context.Context, id int64) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, getWalletWebhook, id)
	var webhook pgtype.Text
	err := row.Scan(&webhook)
	return webhook, err
}

const getWalletsByUserId = `-- name: GetWalletsByUserId :many
SELECT id, address, user_id, webhook FROM wallet WHERE user_id = $1
`

func (q *Queries) GetWalletsByUserId(ctx context.Context, userID int64) ([]Wallet, error) {
	rows, err := q.db.Query(ctx, getWalletsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wallet
	for rows.Next() {
		var i Wallet
		if err := rows.Scan(
			&i.ID,
			&i.Address,
			&i.UserID,
			&i.Webhook,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertWallet = `-- name: InsertWallet :one
INSERT INTO wallet (address, user_id) VALUES ($1, $2) RETURNING id
`

type InsertWalletParams struct {
	Address string `json:"address"`
	UserID  int64  `json:"user_id"`
}

func (q *Queries) InsertWallet(ctx context.Context, arg InsertWalletParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertWallet, arg.Address, arg.UserID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateWalletUserID = `-- name: UpdateWalletUserID :exec
UPDATE wallet SET user_id = $1 WHERE id = $2 AND user_id = $3
`

type UpdateWalletUserIDParams struct {
	UserID   int64 `json:"user_id"`
	ID       int64 `json:"id"`
	UserID_2 int64 `json:"user_id_2"`
}

func (q *Queries) UpdateWalletUserID(ctx context.Context, arg UpdateWalletUserIDParams) error {
	_, err := q.db.Exec(ctx, updateWalletUserID, arg.UserID, arg.ID, arg.UserID_2)
	return err
}

const updateWalletWebhook = `-- name: UpdateWalletWebhook :exec
UPDATE wallet SET webhook = $1 WHERE id = $2
`

type UpdateWalletWebhookParams struct {
	Webhook pgtype.Text `json:"webhook"`
	ID      int64       `json:"id"`
}

func (q *Queries) UpdateWalletWebhook(ctx context.Context, arg UpdateWalletWebhookParams) error {
	_, err := q.db.Exec(ctx, updateWalletWebhook, arg.Webhook, arg.ID)
	return err
}
