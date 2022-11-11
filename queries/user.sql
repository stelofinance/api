-- name: InsertUser :one
INSERT INTO "user" (username, password, created_at) VALUES ($1, $2, $3) RETURNING id;

-- name: InsertWallet :one
INSERT INTO wallet (address, user_id) VALUES ($1, $2) RETURNING id;

-- name: UpdateUserWallet :execrows
UPDATE "user" SET wallet_id = $1 WHERE id = $2;

-- name: GetUser :one
SELECT * FROM "user" WHERE username = $1;

-- name: UpdateUserUsername :exec
UPDATE "user" SET username = $1 WHERE id = $2;

-- name: UpdateUserPassword :exec
UPDATE "user" SET password = $1 WHERE id = $2;

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

-- name: GetWalletByUsername :one
SELECT wallet_id FROM "user" WHERE username = $1;

-- name: GetUserIdByUsername :one
SELECT id FROM "user" WHERE username = $1;

-- name: DeleteUserSessionByUserIdAndWalletId :exec
DELETE FROM user_session WHERE user_id = $1 AND wallet_id = $2;

-- name: UpdateUserSessionUsedAt :execrows
UPDATE user_session SET used_at = $1 WHERE id = $2;
