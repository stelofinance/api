-- name: InsertUserSession :one
INSERT INTO user_session (used_at, user_id, wallet_id) VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateUserSessionWallet :exec
UPDATE user_session SET wallet_id = $1 WHERE id = $2;

-- name: GetUserSessions :many
SELECT * FROM user_session WHERE user_id = $1;

-- name: DeleteSession :execrows
DELETE FROM user_session WHERE id = $1;

-- name: DeleteSessionsByUserId :exec
DELETE FROM user_session WHERE user_id = $1;

-- name: DeleteUserSessionByUserIdAndWalletId :exec
DELETE FROM user_session WHERE user_id = $1 AND wallet_id = $2;

-- name: UpdateUserSessionUsedAt :execrows
UPDATE user_session SET used_at = $1 WHERE id = $2;
