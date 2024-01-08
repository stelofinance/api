// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: asset.sql

package db

import (
	"context"
)

const createAsset = `-- name: CreateAsset :exec
INSERT INTO asset (name, value) VALUES ($1, $2)
`

type CreateAssetParams struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (q *Queries) CreateAsset(ctx context.Context, arg CreateAssetParams) error {
	_, err := q.db.Exec(ctx, createAsset, arg.Name, arg.Value)
	return err
}

const deleteAsset = `-- name: DeleteAsset :execrows
DELETE FROM asset WHERE id = $1
`

func (q *Queries) DeleteAsset(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAsset, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAssetsByIds = `-- name: GetAssetsByIds :many
SELECT id, name, value FROM asset WHERE id = ANY($1::BIGINT[])
`

func (q *Queries) GetAssetsByIds(ctx context.Context, dollar_1 []int64) ([]Asset, error) {
	rows, err := q.db.Query(ctx, getAssetsByIds, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Asset
	for rows.Next() {
		var i Asset
		if err := rows.Scan(&i.ID, &i.Name, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAssetsIdNameByNames = `-- name: GetAssetsIdNameByNames :many
SELECT id, name FROM asset WHERE name = ANY($1::varchar[])
`

type GetAssetsIdNameByNamesRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) GetAssetsIdNameByNames(ctx context.Context, dollar_1 []string) ([]GetAssetsIdNameByNamesRow, error) {
	rows, err := q.db.Query(ctx, getAssetsIdNameByNames, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAssetsIdNameByNamesRow
	for rows.Next() {
		var i GetAssetsIdNameByNamesRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAssetName = `-- name: UpdateAssetName :execrows
UPDATE asset SET name = $1 WHERE id = $2
`

type UpdateAssetNameParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (q *Queries) UpdateAssetName(ctx context.Context, arg UpdateAssetNameParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateAssetName, arg.Name, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updateAssetValue = `-- name: UpdateAssetValue :execrows
UPDATE asset SET value = $1 WHERE id = $2
`

type UpdateAssetValueParams struct {
	Value int64 `json:"value"`
	ID    int64 `json:"id"`
}

func (q *Queries) UpdateAssetValue(ctx context.Context, arg UpdateAssetValueParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateAssetValue, arg.Value, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
