CREATE TABLE IF NOT EXISTS wallet_asset
(
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    asset_id BIGINT NOT NULL REFERENCES asset(id),
    quantity BIGINT NOT NULL
);