CREATE TABLE IF NOT EXISTS transaction_asset
(
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transaction(id) ON DELETE CASCADE,
    asset_id BIGINT NOT NULL REFERENCES asset(id),
    quantity BIGINT NOT NULL
);
