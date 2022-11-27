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
