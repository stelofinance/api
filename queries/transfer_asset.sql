-- name: CreateTransferAssets :batchexec
INSERT INTO transfer_asset (transfer_id, asset_id, quantity) VALUES ($1, $2, $3);

-- name: GetTransferAssetsByTransferId :many
SELECT ta.asset_id, a.name, ta.quantity
FROM transfer_asset ta
JOIN asset a ON a.id = ta.asset_id
WHERE ta.transfer_id = $1;
