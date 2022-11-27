-- name: CreateWalletAsset :exec
INSERT INTO wallet_asset (wallet_id, asset_id, quantity) VALUES ($1, $2, $3);

-- name: DeleteWalletAsset :execrows
DELETE FROM wallet_asset WHERE wallet_id = $1 AND asset_id = $2;

-- name: GetWalletAssets :many
SELECT * FROM wallet_asset WHERE wallet_id = $1;

-- name: SubtractWalletAssetQuantity :execrows
UPDATE wallet_asset SET quantity = quantity - $1 WHERE wallet_id = $2 AND quantity >= $1 AND asset_id = $3;

-- name: AddWalletAssetQuantity :execrows
UPDATE wallet_asset SET quantity = quantity + $1 WHERE wallet_id = $2 AND asset_id = $3;

