--- TABLES
CREATE TABLE accounts (
	id VARCHAR PRIMARY KEY,
	balance INT NOT NULL,
	owner_name VARCHAR NOT NULL,
	owner_document VARCHAR NOT NULL,
	status VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

--- INDEX
CREATE INDEX idx_accounts_id ON accounts (id);
CREATE INDEX idx_accounts_owner_document ON accounts (owner_document);