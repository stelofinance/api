// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transfer.sql

package db

import (
	"context"
	"time"
)

const getTransferAssets = `-- name: GetTransferAssets :many
SELECT a.id as asset_id, a.name as asset_name, ta.quantity
FROM transfer t
JOIN transfer_asset ta ON ta.transfer_id = t.id
JOIN asset a ON a.id = ta.asset_id
WHERE t.id = $1 AND t.sending_warehouse_id = $2
`

type GetTransferAssetsParams struct {
	ID                 int64 `json:"id"`
	SendingWarehouseID int64 `json:"sending_warehouse_id"`
}

type GetTransferAssetsRow struct {
	AssetID   int64  `json:"asset_id"`
	AssetName string `json:"asset_name"`
	Quantity  int64  `json:"quantity"`
}

func (q *Queries) GetTransferAssets(ctx context.Context, arg GetTransferAssetsParams) ([]GetTransferAssetsRow, error) {
	rows, err := q.db.Query(ctx, getTransferAssets, arg.ID, arg.SendingWarehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransferAssetsRow
	for rows.Next() {
		var i GetTransferAssetsRow
		if err := rows.Scan(&i.AssetID, &i.AssetName, &i.Quantity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransferInboundRequests = `-- name: GetTransferInboundRequests :many
SELECT t.id, w.name as receiving_warehouse_name, t.receiving_warehouse_id, t.status, t.created_at
FROM transfer t
JOIN warehouse w ON t.receiving_warehouse_id = w.id
WHERE t.sending_warehouse_id = $1
`

type GetTransferInboundRequestsRow struct {
	ID                     int64          `json:"id"`
	ReceivingWarehouseName string         `json:"receiving_warehouse_name"`
	ReceivingWarehouseID   int64          `json:"receiving_warehouse_id"`
	Status                 TransferStatus `json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
}

func (q *Queries) GetTransferInboundRequests(ctx context.Context, sendingWarehouseID int64) ([]GetTransferInboundRequestsRow, error) {
	rows, err := q.db.Query(ctx, getTransferInboundRequests, sendingWarehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransferInboundRequestsRow
	for rows.Next() {
		var i GetTransferInboundRequestsRow
		if err := rows.Scan(
			&i.ID,
			&i.ReceivingWarehouseName,
			&i.ReceivingWarehouseID,
			&i.Status,
			&i.CreatedAt,
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

const getTransferOutboundRequests = `-- name: GetTransferOutboundRequests :many
SELECT t.id, w.name as sending_warehouse_name, t.status, t.created_at
FROM transfer t
JOIN warehouse w ON t.sending_warehouse_id = w.id
WHERE t.receiving_warehouse_id = $1
`

type GetTransferOutboundRequestsRow struct {
	ID                   int64          `json:"id"`
	SendingWarehouseName string         `json:"sending_warehouse_name"`
	Status               TransferStatus `json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
}

func (q *Queries) GetTransferOutboundRequests(ctx context.Context, receivingWarehouseID int64) ([]GetTransferOutboundRequestsRow, error) {
	rows, err := q.db.Query(ctx, getTransferOutboundRequests, receivingWarehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransferOutboundRequestsRow
	for rows.Next() {
		var i GetTransferOutboundRequestsRow
		if err := rows.Scan(
			&i.ID,
			&i.SendingWarehouseName,
			&i.Status,
			&i.CreatedAt,
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

const getTransferTotalLiabilityAndReceivingId = `-- name: GetTransferTotalLiabilityAndReceivingId :one
SELECT t.receiving_warehouse_id, SUM(a.value * ta.quantity) as total_liability
FROM transfer t
JOIN transfer_asset ta ON ta.transfer_id = t.id
JOIN asset a ON a.id = ta.asset_id
WHERE t.id = $1 AND t.sending_warehouse_id = $2
GROUP BY t.receiving_warehouse_id
`

type GetTransferTotalLiabilityAndReceivingIdParams struct {
	ID                 int64 `json:"id"`
	SendingWarehouseID int64 `json:"sending_warehouse_id"`
}

type GetTransferTotalLiabilityAndReceivingIdRow struct {
	ReceivingWarehouseID int64 `json:"receiving_warehouse_id"`
	TotalLiability       int64 `json:"total_liability"`
}

func (q *Queries) GetTransferTotalLiabilityAndReceivingId(ctx context.Context, arg GetTransferTotalLiabilityAndReceivingIdParams) (GetTransferTotalLiabilityAndReceivingIdRow, error) {
	row := q.db.QueryRow(ctx, getTransferTotalLiabilityAndReceivingId, arg.ID, arg.SendingWarehouseID)
	var i GetTransferTotalLiabilityAndReceivingIdRow
	err := row.Scan(&i.ReceivingWarehouseID, &i.TotalLiability)
	return i, err
}

const getTransferTotalLiabilityAndSendingId = `-- name: GetTransferTotalLiabilityAndSendingId :one
SELECT t.sending_warehouse_id, SUM(a.value * ta.quantity) as total_liability
FROM transfer t
JOIN transfer_asset ta ON ta.transfer_id = t.id
JOIN asset a ON a.id = ta.asset_id
WHERE t.id = $1 AND t.receiving_warehouse_id = $2
GROUP BY t.sending_warehouse_id
`

type GetTransferTotalLiabilityAndSendingIdParams struct {
	ID                   int64 `json:"id"`
	ReceivingWarehouseID int64 `json:"receiving_warehouse_id"`
}

type GetTransferTotalLiabilityAndSendingIdRow struct {
	SendingWarehouseID int64 `json:"sending_warehouse_id"`
	TotalLiability     int64 `json:"total_liability"`
}

func (q *Queries) GetTransferTotalLiabilityAndSendingId(ctx context.Context, arg GetTransferTotalLiabilityAndSendingIdParams) (GetTransferTotalLiabilityAndSendingIdRow, error) {
	row := q.db.QueryRow(ctx, getTransferTotalLiabilityAndSendingId, arg.ID, arg.ReceivingWarehouseID)
	var i GetTransferTotalLiabilityAndSendingIdRow
	err := row.Scan(&i.SendingWarehouseID, &i.TotalLiability)
	return i, err
}

const getTransfers = `-- name: GetTransfers :many
SELECT
	t.id,
    wr.name as receiving_warehouse_name,
    ws.name as sending_warehouse_name,
    t.receiving_warehouse_id,
    t.sending_warehouse_id,
    t.status,
    t.created_at
FROM transfer t
JOIN warehouse wr ON t.receiving_warehouse_id = wr.id
JOIN warehouse ws ON t.sending_warehouse_id = ws.id
WHERE t.sending_warehouse_id = $1 OR t.receiving_warehouse_id = $1
`

type GetTransfersRow struct {
	ID                     int64          `json:"id"`
	ReceivingWarehouseName string         `json:"receiving_warehouse_name"`
	SendingWarehouseName   string         `json:"sending_warehouse_name"`
	ReceivingWarehouseID   int64          `json:"receiving_warehouse_id"`
	SendingWarehouseID     int64          `json:"sending_warehouse_id"`
	Status                 TransferStatus `json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
}

func (q *Queries) GetTransfers(ctx context.Context, warehouseID int64) ([]GetTransfersRow, error) {
	rows, err := q.db.Query(ctx, getTransfers, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransfersRow
	for rows.Next() {
		var i GetTransfersRow
		if err := rows.Scan(
			&i.ID,
			&i.ReceivingWarehouseName,
			&i.SendingWarehouseName,
			&i.ReceivingWarehouseID,
			&i.SendingWarehouseID,
			&i.Status,
			&i.CreatedAt,
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

const insertTransfer = `-- name: InsertTransfer :one
INSERT INTO transfer (created_at, status, sending_warehouse_id, receiving_warehouse_id) VALUES ($1, $2, $3, $4) RETURNING id
`

type InsertTransferParams struct {
	CreatedAt            time.Time      `json:"created_at"`
	Status               TransferStatus `json:"status"`
	SendingWarehouseID   int64          `json:"sending_warehouse_id"`
	ReceivingWarehouseID int64          `json:"receiving_warehouse_id"`
}

func (q *Queries) InsertTransfer(ctx context.Context, arg InsertTransferParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertTransfer,
		arg.CreatedAt,
		arg.Status,
		arg.SendingWarehouseID,
		arg.ReceivingWarehouseID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateTransferStatus = `-- name: UpdateTransferStatus :execrows
UPDATE transfer
SET status = $1
WHERE
	id = $2
	AND sending_warehouse_id = $3
	AND status = $4
`

type UpdateTransferStatusParams struct {
	Status             TransferStatus `json:"status"`
	ID                 int64          `json:"id"`
	SendingWarehouseID int64          `json:"sending_warehouse_id"`
	Status_2           TransferStatus `json:"status_2"`
}

func (q *Queries) UpdateTransferStatus(ctx context.Context, arg UpdateTransferStatusParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateTransferStatus,
		arg.Status,
		arg.ID,
		arg.SendingWarehouseID,
		arg.Status_2,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
