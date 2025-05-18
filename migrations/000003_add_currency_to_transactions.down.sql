DROP INDEX IF EXISTS idx_transactions_reference_currency;

ALTER TABLE transactions
DROP COLUMN IF EXISTS currency;

CREATE UNIQUE INDEX idx_transactions_reference_id
    ON transactions (reference_id);
