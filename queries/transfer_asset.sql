-- name: CreateTransferAssets :batchexec
INSERT INTO transfer_asset (transfer_id, asset_id, quantity) VALUES ($1, $2, $3);
