-- name: CreateTransactionAssets :batchexec
INSERT INTO transaction_asset (transaction_id, asset_id, quantity) VALUES ($1, $2, $3);

-- name: GetTransactionAssetsByTransactionIds :many
SELECT * FROM transaction_asset WHERE transaction_id = ANY($1::BIGINT[]);
