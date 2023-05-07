-- name: CreateWalletSession :exec
INSERT INTO wallet_session (key, wallet_id, name) VALUES ($1, $2, $3);

-- name: GetWalletSession :one
SELECT wallet_id FROM wallet_session WHERE key = $1;

-- name: GetWalletSessionsByWalletId :many
SELECT * FROM wallet_session WHERE wallet_id = $1;

-- name: DeleteWalletSession :execrows
DELETE FROM wallet_session WHERE id = $1;

-- name: DeleteWalletSessionsByWalletId :exec
DELETE FROM wallet_session WHERE wallet_id = $1;
