-- name: InsertWarehouse :exec
INSERT INTO warehouse (name, user_id, location) VALUES ($1, $2, $3);
