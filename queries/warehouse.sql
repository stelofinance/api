-- name: InsertWarehouse :one
INSERT INTO warehouse (name, user_id, location) VALUES ($1, $2, $3) RETURNING id;

-- name: SubtractWarehouseCollateral :execrows
UPDATE warehouse SET collateral = collateral - $1 WHERE id = $2 AND collateral >= $1;

-- name: AddWarehouseCollateral :execrows
UPDATE warehouse SET collateral = collateral + $1 WHERE id = $2;

-- name: GetWarehouseUserId :one
SELECT user_id FROM warehouse WHERE id = $1;

-- name: UpdateWarehouseUserIdByUsername :exec
UPDATE warehouse SET user_id = "user".id FROM "user" WHERE warehouse.id = $1 AND "user".username = $2;

-- name: GetWarehouseCollateralLiabilityAndRatioLock :one
SELECT collateral, liability, collateral_ratio FROM warehouse WHERE id = $1 FOR UPDATE;

-- name: AddWarehouseLiabiliy :exec
UPDATE warehouse SET liability = liability + $1 WHERE id = $2;
