--- TABLES
CREATE TABLE transactions (
	id VARCHAR PRIMARY KEY,
	account_id VARCHAR NOT NULL,
	destination_account_id VARCHAR,
	amount INT NOT NULL,
	status VARCHAR NOT NULL,
	transaction_type VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

--- INDEX
CREATE INDEX idx_transactions_id ON transactions (id);
CREATE INDEX idx_transactions_account_id ON transactions (account_id);