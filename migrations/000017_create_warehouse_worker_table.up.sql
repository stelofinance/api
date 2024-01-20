CREATE TABLE IF NOT EXISTS warehouse_worker
(
    id BIGSERIAL PRIMARY KEY,
    warehouse_id BIGINT NOT NULL REFERENCES warehouse(id),
    user_id BIGINT NOT NULL REFERENCES "user"(id),
    CONSTRAINT warehouse_worker_warehouse_id_user_id_key UNIQUE (warehouse_id, user_id)
);
