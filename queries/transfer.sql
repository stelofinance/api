-- name: InsertTransfer :one
INSERT INTO transfer (created_at, status, sending_warehouse_id, receiving_warehouse_id) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: GetTransfers :many
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
WHERE t.sending_warehouse_id = sqlc.arg(warehouse_id) OR t.receiving_warehouse_id = sqlc.arg(warehouse_id);

-- name: GetTransferOutboundRequests :many
SELECT t.id, w.name as sending_warehouse_name, t.status, t.created_at
FROM transfer t
JOIN warehouse w ON t.sending_warehouse_id = w.id
WHERE t.receiving_warehouse_id = $1;

-- name: GetTransferInboundRequests :many
SELECT t.id, w.name as receiving_warehouse_name, t.receiving_warehouse_id, t.status, t.created_at
FROM transfer t
JOIN warehouse w ON t.receiving_warehouse_id = w.id
WHERE t.sending_warehouse_id = $1;

-- name: UpdateTransferStatus :execrows
UPDATE transfer
SET status = $1
WHERE
	id = $2
	AND sending_warehouse_id = $3
	AND status = $4;

-- name: GetTransferTotalCollateral :one
SELECT t.receiving_warehouse_id, SUM(a.value * ta.quantity) as total_collateral
FROM transfer t
JOIN transfer_asset ta ON ta.transfer_id = t.id
JOIN asset a ON a.id = ta.asset_id
WHERE t.id = $1 AND t.sending_warehouse_id = $2
GROUP BY t.receiving_warehouse_id;
