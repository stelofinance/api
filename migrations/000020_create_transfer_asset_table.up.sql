CREATE TABLE IF NOT EXISTS transfer_asset
(
    id BIGSERIAL PRIMARY KEY,
    transfer_id BIGINT NOT NULL REFERENCES transfer(id),
    asset_id BIGINT NOT NULL REFERENCES asset(id),
    quantity BIGINT NOT NULL
);
