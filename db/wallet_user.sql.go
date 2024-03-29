// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: wallet_user.sql

package db

import (
	"context"
)

const countAssignedWallet = `-- name: CountAssignedWallet :one
SELECT count(*) FROM wallet_user WHERE user_id = $1 AND wallet_id = $2
`

type CountAssignedWalletParams struct {
	UserID   int64 `json:"user_id"`
	WalletID int64 `json:"wallet_id"`
}

func (q *Queries) CountAssignedWallet(ctx context.Context, arg CountAssignedWalletParams) (int64, error) {
	row := q.db.QueryRow(ctx, countAssignedWallet, arg.UserID, arg.WalletID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createWalletUser = `-- name: CreateWalletUser :exec
INSERT INTO wallet_user (wallet_id, user_id) VALUES ($1, $2)
`

type CreateWalletUserParams struct {
	WalletID int64 `json:"wallet_id"`
	UserID   int64 `json:"user_id"`
}

func (q *Queries) CreateWalletUser(ctx context.Context, arg CreateWalletUserParams) error {
	_, err := q.db.Exec(ctx, createWalletUser, arg.WalletID, arg.UserID)
	return err
}

const deleteWalletUser = `-- name: DeleteWalletUser :exec
DELETE FROM wallet_user WHERE wallet_id = $1 AND user_id = $2
`

type DeleteWalletUserParams struct {
	WalletID int64 `json:"wallet_id"`
	UserID   int64 `json:"user_id"`
}

func (q *Queries) DeleteWalletUser(ctx context.Context, arg DeleteWalletUserParams) error {
	_, err := q.db.Exec(ctx, deleteWalletUser, arg.WalletID, arg.UserID)
	return err
}

const updateWalletUserUserID = `-- name: UpdateWalletUserUserID :execrows
UPDATE wallet_user SET user_id = $1 WHERE wallet_id = $2 AND user_id = $3
`

type UpdateWalletUserUserIDParams struct {
	UserID   int64 `json:"user_id"`
	WalletID int64 `json:"wallet_id"`
	UserID_2 int64 `json:"user_id_2"`
}

func (q *Queries) UpdateWalletUserUserID(ctx context.Context, arg UpdateWalletUserUserIDParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateWalletUserUserID, arg.UserID, arg.WalletID, arg.UserID_2)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
