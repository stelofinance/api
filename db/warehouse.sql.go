// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: warehouse.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addWarehouseCollateral = `-- name: AddWarehouseCollateral :execrows
UPDATE warehouse SET collateral = collateral + $1 WHERE id = $2
`

type AddWarehouseCollateralParams struct {
	Collateral int64 `json:"collateral"`
	ID         int64 `json:"id"`
}

func (q *Queries) AddWarehouseCollateral(ctx context.Context, arg AddWarehouseCollateralParams) (int64, error) {
	result, err := q.db.Exec(ctx, addWarehouseCollateral, arg.Collateral, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const addWarehouseLiabiliy = `-- name: AddWarehouseLiabiliy :exec
UPDATE warehouse SET liability = liability + $1 WHERE id = $2
`

type AddWarehouseLiabiliyParams struct {
	Liability int64 `json:"liability"`
	ID        int64 `json:"id"`
}

func (q *Queries) AddWarehouseLiabiliy(ctx context.Context, arg AddWarehouseLiabiliyParams) error {
	_, err := q.db.Exec(ctx, addWarehouseLiabiliy, arg.Liability, arg.ID)
	return err
}

const getWarehouseCollateralLiabilityAndRatioLock = `-- name: GetWarehouseCollateralLiabilityAndRatioLock :one
SELECT collateral, liability, collateral_ratio FROM warehouse WHERE id = $1 FOR UPDATE
`

type GetWarehouseCollateralLiabilityAndRatioLockRow struct {
	Collateral      int64          `json:"collateral"`
	Liability       int64          `json:"liability"`
	CollateralRatio pgtype.Numeric `json:"collateral_ratio"`
}

func (q *Queries) GetWarehouseCollateralLiabilityAndRatioLock(ctx context.Context, id int64) (GetWarehouseCollateralLiabilityAndRatioLockRow, error) {
	row := q.db.QueryRow(ctx, getWarehouseCollateralLiabilityAndRatioLock, id)
	var i GetWarehouseCollateralLiabilityAndRatioLockRow
	err := row.Scan(&i.Collateral, &i.Liability, &i.CollateralRatio)
	return i, err
}

const getWarehouseUserId = `-- name: GetWarehouseUserId :one
SELECT user_id FROM warehouse WHERE id = $1
`

func (q *Queries) GetWarehouseUserId(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, getWarehouseUserId, id)
	var user_id int64
	err := row.Scan(&user_id)
	return user_id, err
}

const insertWarehouse = `-- name: InsertWarehouse :one
INSERT INTO warehouse (name, user_id, location) VALUES ($1, $2, $3) RETURNING id
`

type InsertWarehouseParams struct {
	Name     string      `json:"name"`
	UserID   int64       `json:"user_id"`
	Location interface{} `json:"location"`
}

func (q *Queries) InsertWarehouse(ctx context.Context, arg InsertWarehouseParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertWarehouse, arg.Name, arg.UserID, arg.Location)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const subtractWarehouseCollateral = `-- name: SubtractWarehouseCollateral :execrows
UPDATE warehouse SET collateral = collateral - $1 WHERE id = $2 AND collateral >= $1
`

type SubtractWarehouseCollateralParams struct {
	Collateral int64 `json:"collateral"`
	ID         int64 `json:"id"`
}

func (q *Queries) SubtractWarehouseCollateral(ctx context.Context, arg SubtractWarehouseCollateralParams) (int64, error) {
	result, err := q.db.Exec(ctx, subtractWarehouseCollateral, arg.Collateral, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const subtractWarehouseLiabiliy = `-- name: SubtractWarehouseLiabiliy :exec
UPDATE warehouse SET liability = liability - $1 WHERE id = $2
`

type SubtractWarehouseLiabiliyParams struct {
	Liability int64 `json:"liability"`
	ID        int64 `json:"id"`
}

func (q *Queries) SubtractWarehouseLiabiliy(ctx context.Context, arg SubtractWarehouseLiabiliyParams) error {
	_, err := q.db.Exec(ctx, subtractWarehouseLiabiliy, arg.Liability, arg.ID)
	return err
}

const updateWarehouseUserIdByUsername = `-- name: UpdateWarehouseUserIdByUsername :exec
UPDATE warehouse SET user_id = "user".id FROM "user" WHERE warehouse.id = $1 AND "user".username = $2
`

type UpdateWarehouseUserIdByUsernameParams struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (q *Queries) UpdateWarehouseUserIdByUsername(ctx context.Context, arg UpdateWarehouseUserIdByUsernameParams) error {
	_, err := q.db.Exec(ctx, updateWarehouseUserIdByUsername, arg.ID, arg.Username)
	return err
}
