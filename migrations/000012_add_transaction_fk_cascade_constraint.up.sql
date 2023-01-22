ALTER TABLE transaction
    ADD CONSTRAINT transaction_receiving_wallet_id_fkey
        FOREIGN KEY (receiving_wallet_id)
        REFERENCES wallet(id)
        ON DELETE CASCADE;

ALTER TABLE transaction
    ADD CONSTRAINT transaction_sending_wallet_id_fkey
        FOREIGN KEY (sending_wallet_id)
        REFERENCES wallet(id)
        ON DELETE CASCADE;