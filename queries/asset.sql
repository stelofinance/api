-- name: CreateAsset :exec
INSERT INTO asset (name, value) VALUES ($1, $2);

-- name: UpdateAssetValue :execrows
UPDATE asset SET value = $1 WHERE id = $2;

-- name: UpdateAssetName :execrows
UPDATE asset SET name = $1 WHERE id = $2;

-- name: DeleteAsset :execrows
DELETE FROM asset WHERE id = $1;

-- name: GetAssetsByIds :many
SELECT * FROM asset WHERE id = ANY($1::BIGINT[]);

-- name: GetAssetIdByName :one
SELECT id FROM asset WHERE name = $1 LIMIT 1;

-- name: GetAssetsIdNameByNames :many
SELECT id, name FROM asset WHERE name = ANY($1::varchar[]);
