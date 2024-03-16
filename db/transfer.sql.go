// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transfer.sql

package db

import (
	"context"
	"time"
)

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