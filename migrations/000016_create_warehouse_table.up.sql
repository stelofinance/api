CREATE TABLE IF NOT EXISTS warehouse
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE COLLATE "en-US-u-ks-level2",
    user_id BIGINT NOT NULL REFERENCES "user"(id),
    location GEOMETRY(POINT) NOT NULL,
    liability BIGINT NOT NULL DEFAULT 0 CHECK (liability >= 0),
    collateral BIGINT NOT NULL DEFAULT 0 CHECK (collateral >= 0),
    collateral_ratio DECIMAL(4, 3) NOT NULL DEFAULT 2 CHECK (collateral_ratio >= 0 AND collateral_ratio <= 2),
    CONSTRAINT check_warehouse_liability_collateral_ratio CHECK (collateral >= liability * collateral_ratio)
);
