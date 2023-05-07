-- name: InsertUserSession :exec
INSERT INTO user_session (key, user_id, wallet_id) VALUES ($1, $2, $3);

-- name: GetUserSession :one
SELECT id, user_id, wallet_id FROM user_session WHERE key = $1;

-- name: UpdateUserSessionWallet :exec
UPDATE user_session SET wallet_id = $1 WHERE id = $2;

-- name: GetUserSessions :many
SELECT id, user_id, wallet_id FROM user_session WHERE user_id = $1;

-- name: DeleteSession :execrows
DELETE FROM user_session WHERE id = $1;

-- name: DeleteSessionsByUserId :exec
DELETE FROM user_session WHERE user_id = $1;

-- name: DeleteUserSessionByUserIdAndWalletId :exec
DELETE FROM user_session WHERE user_id = $1 AND wallet_id = $2;
