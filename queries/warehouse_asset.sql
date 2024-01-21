-- name: CreateWarehouseAsset :exec
INSERT INTO warehouse_asset (warehouse_id, asset_id, quantity) VALUES ($1, $2, $3);

-- name: AddWarehouseAssetQuantity :execrows
UPDATE warehouse_asset SET quantity = quantity + $1 WHERE warehouse_id = $2 AND asset_id = $3;

-- name: SubtractWarehouseAssetQuantity :execrows
UPDATE warehouse_asset SET quantity = quantity - $1 WHERE warehouse_id = $2 AND quantity >= $1 AND asset_id = $3;
