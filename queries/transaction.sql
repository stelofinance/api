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
    DESC LIMIT $2;

-- name: DeleteTransactionById :execrows
DELETE FROM transaction WHERE (sending_wallet_id = $1 OR receiving_wallet_id = $2) AND id = $3;

-- name: DeleteTransactionsById :execrows
DELETE FROM transaction 
    WHERE (sending_wallet_id = $1 OR receiving_wallet_id = $2) AND id = ANY($3::BIGINT[]);
