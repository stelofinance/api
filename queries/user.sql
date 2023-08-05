-- name: InsertUser :one
INSERT INTO "user" (username, password, created_at) VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateUserWallet :execrows
UPDATE "user" SET wallet_id = $1 WHERE id = $2;

-- name: GetUser :one
SELECT * FROM "user" WHERE username = $1;

-- name: GetUserById :one
SELECT * FROM "user" WHERE id = $1;

-- name: UpdateUserUsername :exec
UPDATE "user" SET username = $1 WHERE id = $2;

-- name: UpdateUserPassword :exec
UPDATE "user" SET password = $1 WHERE id = $2;

-- name: GetWalletByUsername :one
SELECT wallet_id FROM "user" WHERE username = $1;

-- name: GetUserSessionInfo :one
SELECT
    "user".username,
    "user".wallet_id AS primary_wallet_id,
    wallet.address AS wallet_address,
    wallet.user_id as wallet_user_id
FROM
    "user",
    wallet
WHERE
    "user".id = sqlc.arg(user_id)::bigint
    AND wallet.id = sqlc.arg(wallet_id)::bigint;

-- name: GetUserIdByUsername :one
SELECT id FROM "user" WHERE username = $1;

-- name: GetAssignedUsersByWalletId :many
SELECT "user".id, "user".username 
FROM "user"
INNER JOIN wallet_user 
    ON "user".id = wallet_user.user_id 
        AND wallet_user.wallet_id = $1;
