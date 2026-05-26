CREATE TABLE IF NOT EXISTS transfer_details (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    transaction_id INT UNIQUE NOT NULL,
    counterparty_id INT NOT NULL,
    notes TEXT,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (counterparty_id) REFERENCES users(id)
);