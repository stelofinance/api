// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: wallet_asset.sql

package db

import (
	"context"
)

const addWalletAssetQuantity = `-- name: AddWalletAssetQuantity :execrows
UPDATE wallet_asset SET quantity = quantity + $1 WHERE wallet_id = $2 AND asset_id = $3
`

type AddWalletAssetQuantityParams struct {
	Quantity int64 `json:"quantity"`
	WalletID int64 `json:"wallet_id"`
	AssetID  int64 `json:"asset_id"`
}

func (q *Queries) AddWalletAssetQuantity(ctx context.Context, arg AddWalletAssetQuantityParams) (int64, error) {
	result, err := q.db.Exec(ctx, addWalletAssetQuantity, arg.Quantity, arg.WalletID, arg.AssetID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const createWalletAsset = `-- name: CreateWalletAsset :exec
INSERT INTO wallet_asset (wallet_id, asset_id, quantity) VALUES ($1, $2, $3)
`

type CreateWalletAssetParams struct {
	WalletID int64 `json:"wallet_id"`
	AssetID  int64 `json:"asset_id"`
	Quantity int64 `json:"quantity"`
}

func (q *Queries) CreateWalletAsset(ctx context.Context, arg CreateWalletAssetParams) error {
	_, err := q.db.Exec(ctx, createWalletAsset, arg.WalletID, arg.AssetID, arg.Quantity)
	return err
}

const deleteWalletAsset = `-- name: DeleteWalletAsset :execrows
DELETE FROM wallet_asset WHERE wallet_id = $1 AND asset_id = $2
`

type DeleteWalletAssetParams struct {
	WalletID int64 `json:"wallet_id"`
	AssetID  int64 `json:"asset_id"`
}

func (q *Queries) DeleteWalletAsset(ctx context.Context, arg DeleteWalletAssetParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteWalletAsset, arg.WalletID, arg.AssetID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteWalletAssets = `-- name: DeleteWalletAssets :execrows
DELETE FROM wallet_asset WHERE wallet_id = $1
`

func (q *Queries) DeleteWalletAssets(ctx context.Context, walletID int64) (int64, error) {
	result, err := q.db.Exec(ctx, deleteWalletAssets, walletID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getWalletAssets = `-- name: GetWalletAssets :many
SELECT id, wallet_id, asset_id, quantity FROM wallet_asset WHERE wallet_id = $1
`

func (q *Queries) GetWalletAssets(ctx context.Context, walletID int64) ([]WalletAsset, error) {
	rows, err := q.db.Query(ctx, getWalletAssets, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WalletAsset
	for rows.Next() {
		var i WalletAsset
		if err := rows.Scan(
			&i.ID,
			&i.WalletID,
			&i.AssetID,
			&i.Quantity,
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

const subtractWalletAssetQuantity = `-- name: SubtractWalletAssetQuantity :execrows
UPDATE wallet_asset SET quantity = quantity - $1 WHERE wallet_id = $2 AND quantity >= $1 AND asset_id = $3
`

type SubtractWalletAssetQuantityParams struct {
	Quantity int64 `json:"quantity"`
	WalletID int64 `json:"wallet_id"`
	AssetID  int64 `json:"asset_id"`
}

func (q *Queries) SubtractWalletAssetQuantity(ctx context.Context, arg SubtractWalletAssetQuantityParams) (int64, error) {
	result, err := q.db.Exec(ctx, subtractWalletAssetQuantity, arg.Quantity, arg.WalletID, arg.AssetID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
