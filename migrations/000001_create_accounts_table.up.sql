CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    balance NUMERIC NOT NULL CHECK (balance >= 0),
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'suspended', 'closed')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_currency UNIQUE (user_id, currency)
    );

CREATE INDEX idx_accounts_user_id ON accounts (user_id);
CREATE INDEX idx_accounts_status ON accounts (status);