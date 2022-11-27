-- name: CountAssignedWallet :one
SELECT count(*) FROM wallet_user WHERE user_id = $1 AND wallet_id = $2;

-- name: CreateWalletUser :exec
INSERT INTO wallet_user (wallet_id, user_id) VALUES ($1, $2);

-- name: DeleteWalletUser :exec
DELETE FROM wallet_user WHERE wallet_id = $1 AND user_id = $2;

