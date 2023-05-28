-- name: CreateTransaction :one
INSERT 
    INTO transaction (sending_wallet_id, receiving_wallet_id, created_at, memo)
    VALUES ($1, $2, $3, $4)
    RETURNING id;

-- name: GetTransactions :many
SELECT * FROM transaction 
    WHERE sending_wallet_id = $1 
    OR receiving_wallet_id = $1 
    ORDER BY created_at 
    DESC LIMIT $2 
    OFFSET $3;

-- name: GetTransactionsDetailed :many
SELECT
	tx.id,
    tx.sending_wallet_id,
    tx.receiving_wallet_id,
    sending_wallet.address AS sending_address,
    receiving_wallet.address AS receiving_address,
    sending_user.username AS sending_username,
    receiving_user.username AS receiving_username,
    tx.memo,
    tx.created_at
FROM
    transaction AS tx
JOIN 
    wallet AS sending_wallet ON tx.sending_wallet_id = sending_wallet.id
JOIN 
    wallet AS receiving_wallet ON tx.receiving_wallet_id = receiving_wallet.id
LEFT JOIN
    "user" AS sending_user ON tx.sending_wallet_id = sending_user.wallet_id
LEFT JOIN
    "user" AS receiving_user ON tx.receiving_wallet_id = receiving_user.wallet_id
WHERE
    tx.sending_wallet_id = $1
    OR tx.receiving_wallet_id = $1
ORDER BY
    tx.created_at DESC
LIMIT
    $2
OFFSET
    $3;

-- name: DeleteTransactionById :execrows
DELETE FROM transaction WHERE (sending_wallet_id = $1 OR receiving_wallet_id = $2) AND id = $3;

-- name: DeleteTransactionsById :execrows
DELETE FROM transaction 
    WHERE (sending_wallet_id = $1 OR receiving_wallet_id = $2) AND id = ANY($3::BIGINT[]);
