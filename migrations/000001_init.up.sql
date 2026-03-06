-- Create users table
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    balance DECIMAL(18, 2) NOT NULL DEFAULT 100.00,
    CONSTRAINT chk_users_balance_non_negative CHECK (balance >= 0)
);

-- Insert users
INSERT INTO users (id, balance) VALUES ('123e4567-e89b-12d3-a456-426614174000', 100.00);
INSERT INTO users (id, balance) VALUES ('123e4567-e89b-12d3-a456-426614174001', 100.00);

-- Create withdrawals table
CREATE TABLE IF NOT EXISTS withdrawals(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(18, 2) NOT NULL CONSTRAINT chk_withdrawals_amount_positive CHECK (amount > 0),
    destination UUID NOT NULL REFERENCES users(id),
    idempotency_key VARCHAR(128) NOT NULL,
    payload_hash TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status = 'pending'),
    UNIQUE (user_id, idempotency_key)
);
