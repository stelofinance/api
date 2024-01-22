-- name: CreateWarehouseAsset :exec
INSERT INTO warehouse_asset (warehouse_id, asset_id, quantity) VALUES ($1, $2, $3);

-- name: AddWarehouseAssetQuantity :execrows
UPDATE warehouse_asset SET quantity = quantity + $1 WHERE warehouse_id = $2 AND asset_id = $3;

-- name: SubtractWarehouseAssetQuantity :execrows
UPDATE warehouse_asset SET quantity = quantity - $1 WHERE warehouse_id = $2 AND quantity >= $1 AND asset_id = $3;

-- name: GetWarehouseAssets :many
SELECT a.id, a.name, a.value, wa.quantity
FROM warehouse_asset wa
JOIN asset a ON a.id = wa.asset_id
WHERE wa.warehouse_id = $1;
