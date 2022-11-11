-- name: CreateWallet :exec
INSERT INTO wallet (address, user_id) VALUES ($1, $2);

-- name: CreateWalletAsset :exec
INSERT INTO wallet_asset (wallet_id, asset_id, quantity) VALUES ($1, $2, $3);

-- name: DeleteWalletAsset :execrows
DELETE FROM wallet_asset WHERE wallet_id = $1 AND asset_id = $2;

-- name: CountWalletsByIdAndUserId :one
SELECT count(*) FROM wallet WHERE id = $1 AND user_id = $2;

-- name: GetWalletsByUserId :many
SELECT * FROM wallet WHERE user_id = $1;

-- name: GetAssignedWalletsByUserId :many
SELECT wallet.id, wallet.address, wallet.user_id 
FROM wallet 
INNER JOIN wallet_user 
    ON wallet.id = wallet_user.wallet_id 
        AND wallet_user.user_id = $1;

-- name: CountAssignedWallet :one
SELECT count(*) FROM wallet_user WHERE user_id = $1 AND wallet_id = $2;

-- name: GetWalletAssets :many
SELECT * FROM wallet_asset WHERE wallet_id = $1;

-- name: SubtractWalletAssetQuantity :execrows
UPDATE wallet_asset SET quantity = quantity - $1 WHERE wallet_id = $2 AND quantity >= $1 AND asset_id = $3;

-- name: GetWalletIdByAddress :one
SELECT id FROM wallet WHERE address = $1;

-- name: AddWalletAssetQuantity :execrows
UPDATE wallet_asset SET quantity = quantity + $1 WHERE wallet_id = $2 AND asset_id = $3;

-- name: CreateWalletUser :exec
INSERT INTO wallet_user (wallet_id, user_id) VALUES ($1, $2);

-- name: DeleteWalletUser :exec
DELETE FROM wallet_user WHERE wallet_id = $1 AND user_id = $2;

-- name: GetAssignedUsersByWalletId :many
SELECT "user".id, "user".username 
FROM "user"
INNER JOIN wallet_user 
    ON "user".id = wallet_user.user_id 
        AND wallet_user.wallet_id = $1;

-- name: CreateWalletSession :one
INSERT INTO wallet_session (wallet_id, name, used_at) VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateWalletSessionUsedAt :execrows
UPDATE wallet_session SET used_at = $1 WHERE id = $2;

-- name: GetWalletSessionsByWalletId :many
SELECT * FROM wallet_session WHERE wallet_id = $1;

-- name: DeleteWalletSession :execrows
DELETE FROM wallet_session WHERE id = $1;

-- name: DeleteWalletSessionsByWalletId :exec
DELETE FROM wallet_session WHERE wallet_id = $1;
