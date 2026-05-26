CREATE TABLE IF NOT EXISTS wallets (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users(id),
    balance BIGINT DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP 
);