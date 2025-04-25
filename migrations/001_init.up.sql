CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    wallet_id UUID REFERENCES wallets(id),
    operation_type TEXT NOT NULL CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);