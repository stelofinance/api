CREATE TABLE IF NOT EXISTS user_session
(
    id BIGSERIAL PRIMARY KEY,
    used_at TIMESTAMPTZ NOT NULL,
    user_id BIGINT NOT NULL REFERENCES "user"(id),
    wallet_id BIGINT NOT NULL REFERENCES wallet(id)
);