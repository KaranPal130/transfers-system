-- accounts table
CREATE TABLE accounts (
    account_id BIGINT PRIMARY KEY,
    balance DECIMAL(20, 5) NOT NULL
);

-- transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    source_account_id BIGINT REFERENCES accounts(account_id),
    destination_account_id BIGINT REFERENCES accounts(account_id),
    amount DECIMAL(20, 5) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_accounts_account_id ON accounts(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_source_account_id ON transactions(source_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_destination_account_id ON transactions(destination_account_id);