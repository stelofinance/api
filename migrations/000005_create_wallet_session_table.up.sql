CREATE TABLE IF NOT EXISTS wallet_session
(
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    name VARCHAR(32),
    used_at TIMESTAMPTZ NOT NULL
);