CREATE TABLE IF NOT EXISTS wallet_user
(
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    user_id BIGINT NOT NULL REFERENCES "user"(id)
);