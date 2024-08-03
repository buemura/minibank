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

--- SEED
INSERT INTO accounts (id, balance, owner_name, owner_document, status, created_at, updated_at)
VALUES 
	('clzc1pqd0000008mnfmkq9r50', 34847, 'Bruno Uemura', '2354448788', 'ACTIVE', NOW(), NOW()),
    ('cjd7n3v6g0001gq9g7j2m3pbk', 1000, 'John Doe', '12345678901', 'ACTIVE', NOW(), NOW()),
    ('cjd7n3v6g0002gq9g7j2m3pbk', 1500, 'Jane Smith', '98765432100', 'INACTIVE', NOW(), NOW()),
    ('cjd7n3v6g0003gq9g7j2m3pbk', 2000, 'Alice Johnson', '13579246802', 'BLOCKED', NOW(), NOW()),
    ('cjd7n3v6g0004gq9g7j2m3pbk', 2500, 'Bob Brown', '24680135709', 'ACTIVE', NOW(), NOW()),
    ('cjd7n3v6g0005gq9g7j2m3pbk', 3000, 'Charlie Davis', '31415926535', 'INACTIVE', NOW(), NOW());
