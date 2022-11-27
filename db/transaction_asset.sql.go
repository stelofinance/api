// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: transaction_asset.sql

package db

import (
	"context"
)

const getTransactionAssetsByTransactionIds = `-- name: GetTransactionAssetsByTransactionIds :many
SELECT id, transaction_id, asset_id, quantity FROM transaction_asset WHERE transaction_id = ANY($1::BIGINT[])
`

func (q *Queries) GetTransactionAssetsByTransactionIds(ctx context.Context, dollar_1 []int64) ([]TransactionAsset, error) {
	rows, err := q.db.Query(ctx, getTransactionAssetsByTransactionIds, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TransactionAsset
	for rows.Next() {
		var i TransactionAsset
		if err := rows.Scan(
			&i.ID,
			&i.TransactionID,
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
