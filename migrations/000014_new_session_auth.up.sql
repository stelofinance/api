BEGIN;

-- Empty tables due to NOT NULL conflict
TRUNCATE user_session;
TRUNCATE wallet_session;

-- Add new key column
ALTER TABLE user_session ADD COLUMN key VARCHAR(32) NOT NULL;
ALTER TABLE wallet_session ADD COLUMN key VARCHAR(32) NOT NULL;

-- Remove used_at columns since they will lose support
ALTER TABLE user_session DROP COLUMN IF EXISTS used_at;
ALTER TABLE wallet_session DROP COLUMN IF EXISTS used_at;

COMMIT;