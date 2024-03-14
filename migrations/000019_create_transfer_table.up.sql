BEGIN;

CREATE TYPE transfer_status AS ENUM ('open', 'declined', 'approved', 'cleared');

CREATE TABLE IF NOT EXISTS transfer
(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    status transfer_status NOT NULL,
    sending_warehouse_id BIGINT NOT NULL REFERENCES warehouse(id),
    receiving_warehouse_id BIGINT NOT NULL REFERENCES warehouse(id)
);

COMMIT;
