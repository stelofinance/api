CREATE TABLE IF NOT EXISTS asset
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE,
    value BIGINT NOT NULL
);