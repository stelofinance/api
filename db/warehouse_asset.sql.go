// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: warehouse_asset.sql

package db

import (
	"context"
)

const addWarehouseAssetQuantity = `-- name: AddWarehouseAssetQuantity :execrows
UPDATE warehouse_asset SET quantity = quantity + $1 WHERE warehouse_id = $2 AND asset_id = $3
`

type AddWarehouseAssetQuantityParams struct {
	Quantity    int64 `json:"quantity"`
	WarehouseID int64 `json:"warehouse_id"`
	AssetID     int64 `json:"asset_id"`
}

func (q *Queries) AddWarehouseAssetQuantity(ctx context.Context, arg AddWarehouseAssetQuantityParams) (int64, error) {
	result, err := q.db.Exec(ctx, addWarehouseAssetQuantity, arg.Quantity, arg.WarehouseID, arg.AssetID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const createWarehouseAsset = `-- name: CreateWarehouseAsset :exec
INSERT INTO warehouse_asset (warehouse_id, asset_id, quantity) VALUES ($1, $2, $3)
`

type CreateWarehouseAssetParams struct {
	WarehouseID int64 `json:"warehouse_id"`
	AssetID     int64 `json:"asset_id"`
	Quantity    int64 `json:"quantity"`
}

func (q *Queries) CreateWarehouseAsset(ctx context.Context, arg CreateWarehouseAssetParams) error {
	_, err := q.db.Exec(ctx, createWarehouseAsset, arg.WarehouseID, arg.AssetID, arg.Quantity)
	return err
}
