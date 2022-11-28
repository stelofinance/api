-- name: CreateWallet :exec
INSERT INTO wallet (address, user_id) VALUES ($1, $2);

-- name: InsertWallet :one
INSERT INTO wallet (address, user_id) VALUES ($1, $2) RETURNING id;

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

-- name: GetWalletIdByAddress :one
SELECT id FROM wallet WHERE address = $1;

-- name: UpdateWalletUserID :exec
UPDATE wallet SET user_id = $1 WHERE id = $2 AND user_id = $3;
