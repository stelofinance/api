ALTER TABLE transaction
    DROP CONSTRAINT transaction_receiving_wallet_id_fkey;

ALTER TABLE transaction
    DROP CONSTRAINT transaction_sending_wallet_id_fkey;
