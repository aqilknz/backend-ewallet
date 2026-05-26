CREATE TABLE IF NOT EXISTS topup_details (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    transaction_id INT UNIQUE NOT NULL, 
    payment_method_id INT NOT NULL,
    discount BIGINT NOT NULL,
    tax BIGINT NOT NULL,
    sub_total BIGINT NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);