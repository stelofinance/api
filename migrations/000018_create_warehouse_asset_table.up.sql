CREATE TABLE IF NOT EXISTS warehouse_asset
(
    id BIGSERIAL PRIMARY KEY,
    warehouse_id BIGINT NOT NULL REFERENCES warehouse(id),
    asset_id BIGINT NOT NULL REFERENCES asset(id),
    quantity BIGINT NOT NULL
);
