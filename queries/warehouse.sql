-- name: InsertWarehouse :exec
INSERT INTO warehouse (name, user_id, location) VALUES ($1, $2, $3);

-- name: SubtractWarehouseCollateralQuantity :execrows
UPDATE warehouse SET collateral = collateral - $1 WHERE id = $2 AND collateral >= $1;

-- name: AddWarehouseCollateralQuantity :execrows
UPDATE warehouse SET collateral = collateral + $1 WHERE id = $2;

-- name: GetWarehouseUserId :one
SELECT user_id FROM warehouse WHERE id = $1;
