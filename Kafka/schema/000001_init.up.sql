CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    balance DECIMAL NOT NULL,
    currency VARCHAR(3) NOT NULL,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(id),
    amount DECIMAL NOT NULL,
    transaction_type VARCHAR(20),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
