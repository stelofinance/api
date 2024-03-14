-- name: InsertTransfer :one
INSERT INTO transfer (created_at, status, sending_warehouse_id, receiving_warehouse_id) VALUES ($1, $2, $3, $4) RETURNING id;
