package database

import (
	"database/sql"

	"github.com/buemura/minibank/svc-account/internal/core/domain/account"
)

type SqlAccountRepository struct {
	db *sql.DB
}

func NewSqlAccountRepository() *SqlAccountRepository {
	return &SqlAccountRepository{
		db: Conn,
	}
}

func (r *SqlAccountRepository) FindById(id string) (*account.Account, error) {
	var acc account.Account
	err := r.db.QueryRow(
		`SELECT id, balance, owner_name, owner_document, status, created_at, updated_at
		FROM "accounts" 
		WHERE id = $1`,
		id,
	).Scan(&acc.ID, &acc.Balance, &acc.OwnerName, &acc.OwnerDocument, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func (r *SqlAccountRepository) FindByOwnerDocument(document string) (*account.Account, error) {
	var acc account.Account
	err := r.db.QueryRow(
		`SELECT id, balance, owner_name, owner_document, status, created_at, updated_at
		FROM "accounts" 
		WHERE owner_document = $1`,
		document,
	).Scan(&acc.ID, &acc.Balance, &acc.OwnerName, &acc.OwnerDocument, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// if acc == nil {
	// 	return nil, account.ErrAccountNotFound
	// }
	return &acc, nil
}

func (r *SqlAccountRepository) Create(acc *account.Account) (*account.Account, error) {
	_, err := r.db.Exec(
		`INSERT INTO "accounts" (id, balance, owner_name, owner_document, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		acc.ID, acc.Balance, acc.OwnerName, acc.OwnerDocument, acc.Status, acc.CreatedAt, acc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (r *SqlAccountRepository) Update(acc *account.Account) (*account.Account, error) {
	_, err := r.db.Exec(
		`UPDATE "accounts" SET balance=$1, owner_name=$2, owner_document=$3, status=$4, updated_at=$5 WHERE id=$6`,
		acc.Balance, acc.OwnerName, acc.OwnerDocument, acc.Status, acc.UpdatedAt, acc.ID,
	)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
