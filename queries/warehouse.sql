-- name: InsertWarehouse :one
INSERT INTO warehouse (name, user_id, location) VALUES ($1, $2, $3) RETURNING id;

-- name: SubtractWarehouseCollateral :execrows
UPDATE warehouse SET collateral = collateral - $1 WHERE id = $2 AND collateral >= $1;

-- name: AddWarehouseCollateral :execrows
UPDATE warehouse SET collateral = collateral + $1 WHERE id = $2;

-- name: UpdateWarehouseLiability :execrows
UPDATE warehouse SET liability = $1 WHERE id = $2;

-- name: GetWarehouseUserId :one
SELECT user_id FROM warehouse WHERE id = $1;

-- name: UpdateWarehouseUserIdByUsername :exec
UPDATE warehouse SET user_id = "user".id FROM "user" WHERE warehouse.id = $1 AND "user".username = $2;

-- name: GetWarehouseCollateralLiabilityAndRatioLock :one
SELECT collateral, liability, collateral_ratio FROM warehouse WHERE id = $1 FOR UPDATE;

-- name: GetWarehouseCollateralLiabilityAndRatio :one
SELECT collateral, liability, collateral_ratio FROM warehouse WHERE id = $1;

-- name: AddWarehouseLiability :exec
UPDATE warehouse SET liability = liability + $1 WHERE id = $2;

-- name: SubtractWarehouseLiability :exec
UPDATE warehouse SET liability = liability - $1 WHERE id = $2;

-- name: GetWarehousesCollateralTotals :one
SELECT
    COALESCE(SUM(wa.quantity * a.value), 0)::BIGINT AS warehouse_assets_total,
    COALESCE(SUM(ta.quantity * a.value), 0)::BIGINT AS transferred_assets_total
FROM
    warehouse w
JOIN
	warehouse_asset wa ON wa.warehouse_id = w.id
LEFT JOIN
    transfer t ON t.receiving_warehouse_id = w.id AND t.status = 'approved' AND t.receiving_warehouse_id = $1
LEFT JOIN
    transfer_asset ta ON ta.transfer_id = t.id
JOIN
    asset a ON a.id = wa.asset_id
WHERE
    w.id = $1;
