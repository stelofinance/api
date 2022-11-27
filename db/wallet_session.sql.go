// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: wallet_session.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createWalletSession = `-- name: CreateWalletSession :one
INSERT INTO wallet_session (wallet_id, name, used_at) VALUES ($1, $2, $3) RETURNING id
`

type CreateWalletSessionParams struct {
	WalletID int64          `json:"wallet_id"`
	Name     sql.NullString `json:"name"`
	UsedAt   time.Time      `json:"used_at"`
}

func (q *Queries) CreateWalletSession(ctx context.Context, arg CreateWalletSessionParams) (int64, error) {
	row := q.db.QueryRow(ctx, createWalletSession, arg.WalletID, arg.Name, arg.UsedAt)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteWalletSession = `-- name: DeleteWalletSession :execrows
DELETE FROM wallet_session WHERE id = $1
`

func (q *Queries) DeleteWalletSession(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.Exec(ctx, deleteWalletSession, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteWalletSessionsByWalletId = `-- name: DeleteWalletSessionsByWalletId :exec
DELETE FROM wallet_session WHERE wallet_id = $1
`

func (q *Queries) DeleteWalletSessionsByWalletId(ctx context.Context, walletID int64) error {
	_, err := q.db.Exec(ctx, deleteWalletSessionsByWalletId, walletID)
	return err
}

const getWalletSessionsByWalletId = `-- name: GetWalletSessionsByWalletId :many
SELECT id, wallet_id, name, used_at FROM wallet_session WHERE wallet_id = $1
`

func (q *Queries) GetWalletSessionsByWalletId(ctx context.Context, walletID int64) ([]WalletSession, error) {
	rows, err := q.db.Query(ctx, getWalletSessionsByWalletId, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WalletSession
	for rows.Next() {
		var i WalletSession
		if err := rows.Scan(
			&i.ID,
			&i.WalletID,
			&i.Name,
			&i.UsedAt,
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

const updateWalletSessionUsedAt = `-- name: UpdateWalletSessionUsedAt :execrows
UPDATE wallet_session SET used_at = $1 WHERE id = $2
`

type UpdateWalletSessionUsedAtParams struct {
	UsedAt time.Time `json:"used_at"`
	ID     int64     `json:"id"`
}

func (q *Queries) UpdateWalletSessionUsedAt(ctx context.Context, arg UpdateWalletSessionUsedAtParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateWalletSessionUsedAt, arg.UsedAt, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
