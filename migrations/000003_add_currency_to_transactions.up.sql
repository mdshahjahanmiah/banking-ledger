ALTER TABLE transactions
    ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD';

DROP INDEX IF EXISTS idx_transactions_reference_id;

CREATE UNIQUE INDEX idx_transactions_reference_currency
    ON transactions (reference_id, currency);
