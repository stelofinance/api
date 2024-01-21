-- name: InsertWarehouseWorkerByUsername :exec
INSERT INTO warehouse_worker (warehouse_id, user_id)
SELECT
    $1 AS warehouse_id,
    u.id AS user_id
FROM
    "user" u
WHERE
    u.username = $2;

-- name: InsertWarehouseWorker :exec
INSERT INTO warehouse_worker (warehouse_id, user_id) VALUES ($1, $2);

-- name: ExistsWarehouseWorkerByUsername :one
SELECT EXISTS(
    SELECT 1 
    FROM warehouse_worker ww
    JOIN "user" u ON ww.user_id = u.id
    WHERE ww.warehouse_id = $1 AND u.username = $2
);

-- name: ExistsWarehouseWorker :one
SELECT EXISTS(
    SELECT 1 
    FROM warehouse_worker
    WHERE warehouse_id = $1 AND user_id = $2
);

-- name: DeleteWarehouseWorker :exec
DELETE FROM warehouse_worker WHERE id = $1 AND user_id != $2;

-- name: GetWarehouseWorkers :many
SELECT ww.id, u.username
FROM warehouse_worker ww
JOIN "user" u ON ww.user_id = u.id
WHERE ww.warehouse_id = $1;
 
