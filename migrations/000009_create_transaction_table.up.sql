CREATE TABLE IF NOT EXISTS transaction
(
    id BIGSERIAL PRIMARY KEY,
    sending_wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    receiving_wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    created_at TIMESTAMPTZ NOT NULL,
    memo VARCHAR(64) 
);