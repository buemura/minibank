package database

import (
	"context"
	"errors"

	"github.com/buemura/minibank/svc-account/internal/core/domain/account"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxAccountRepository struct {
	db *pgxpool.Pool
}

func NewPgxAccountRepository() *PgxAccountRepository {
	return &PgxAccountRepository{
		db: Conn,
	}
}

func (r *PgxAccountRepository) FindById(id string) (*account.Account, error) {
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, balance, owner_name, owner_document, status, created_at, updated_at
		FROM "accounts" 
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	acc, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[account.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, err
	}
	return acc, nil
}

func (r *PgxAccountRepository) FindByOwnerDocument(document string) (*account.Account, error) {
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, balance, owner_name, owner_document, status, created_at, updated_at
		FROM "accounts" 
		WHERE owner_document = $1`,
		document,
	)
	if err != nil {
		return nil, err
	}

	acc, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[account.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return acc, nil
}

func (r *PgxAccountRepository) Create(acc *account.Account) (*account.Account, error) {
	_, err := r.db.Query(
		context.Background(),
		`INSERT INTO "accounts" (id, balance, owner_name, owner_document, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		acc.ID, acc.Balance, acc.OwnerName, acc.OwnerDocument, acc.Status, acc.CreatedAt, acc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (r *PgxAccountRepository) Update(acc *account.Account) (*account.Account, error) {
	_, err := r.db.Query(
		context.Background(),
		`UPDATE "accounts" SET balance=$1, owner_name=$2, owner_document=$3, status=$4, updated_at=$5 WHERE id=$6`,
		acc.Balance, acc.OwnerName, acc.OwnerDocument, acc.Status, acc.UpdatedAt, acc.ID,
	)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
