// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const getAssignedUsersByWalletId = `-- name: GetAssignedUsersByWalletId :many
SELECT "user".id, "user".username 
FROM "user"
INNER JOIN wallet_user 
    ON "user".id = wallet_user.user_id 
        AND wallet_user.wallet_id = $1
`

type GetAssignedUsersByWalletIdRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (q *Queries) GetAssignedUsersByWalletId(ctx context.Context, walletID int64) ([]GetAssignedUsersByWalletIdRow, error) {
	rows, err := q.db.Query(ctx, getAssignedUsersByWalletId, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAssignedUsersByWalletIdRow
	for rows.Next() {
		var i GetAssignedUsersByWalletIdRow
		if err := rows.Scan(&i.ID, &i.Username); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT id, username, password, created_at, wallet_id FROM "user" WHERE username = $1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.WalletID,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, password, created_at, wallet_id FROM "user" WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.WalletID,
	)
	return i, err
}

const getUserIdByUsername = `-- name: GetUserIdByUsername :one
SELECT id FROM "user" WHERE username = $1
`

func (q *Queries) GetUserIdByUsername(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRow(ctx, getUserIdByUsername, username)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getWalletByUsername = `-- name: GetWalletByUsername :one
SELECT wallet_id FROM "user" WHERE username = $1
`

func (q *Queries) GetWalletByUsername(ctx context.Context, username string) (sql.NullInt64, error) {
	row := q.db.QueryRow(ctx, getWalletByUsername, username)
	var wallet_id sql.NullInt64
	err := row.Scan(&wallet_id)
	return wallet_id, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO "user" (username, password, created_at) VALUES ($1, $2, $3) RETURNING id
`

type InsertUserParams struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertUser, arg.Username, arg.Password, arg.CreatedAt)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE "user" SET password = $1 WHERE id = $2
`

type UpdateUserPasswordParams struct {
	Password string `json:"password"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Password, arg.ID)
	return err
}

const updateUserUsername = `-- name: UpdateUserUsername :exec
UPDATE "user" SET username = $1 WHERE id = $2
`

type UpdateUserUsernameParams struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateUserUsername(ctx context.Context, arg UpdateUserUsernameParams) error {
	_, err := q.db.Exec(ctx, updateUserUsername, arg.Username, arg.ID)
	return err
}

const updateUserWallet = `-- name: UpdateUserWallet :execrows
UPDATE "user" SET wallet_id = $1 WHERE id = $2
`

type UpdateUserWalletParams struct {
	WalletID sql.NullInt64 `json:"wallet_id"`
	ID       int64         `json:"id"`
}

func (q *Queries) UpdateUserWallet(ctx context.Context, arg UpdateUserWalletParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateUserWallet, arg.WalletID, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
