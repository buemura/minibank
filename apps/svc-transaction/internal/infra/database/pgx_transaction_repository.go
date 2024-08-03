package database

import (
	"context"
	"errors"

	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxTransactionRepository struct {
	db *pgxpool.Pool
}

func NewPgxTransactionRepository() *PgxTransactionRepository {
	return &PgxTransactionRepository{db: Conn}
}

func (r *PgxTransactionRepository) FindById(id string) (*transaction.Transaction, error) {
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, account_id, destination_account_id, amount, status, transaction_type, created_at, updated_at
		FROM "transactions"
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	trx, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[transaction.Transaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transaction.ErrTransactionNotFound
		}
		return nil, err
	}
	return trx, nil
}

func (r *PgxTransactionRepository) FindByAccountId(accountID string) ([]*transaction.Transaction, error) {
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, account_id, destination_account_id, amount, status, transaction_type, created_at, updated_at
		FROM "transactions"
		WHERE account_id = $1`,
		accountID,
	)
	if err != nil {
		return nil, err
	}

	trxs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[transaction.Transaction])
	if err != nil {
		return nil, err
	}
	return trxs, nil
}

func (r *PgxTransactionRepository) Create(trx *transaction.Transaction) (*transaction.Transaction, error) {
	_, err := r.db.Query(
		context.Background(),
		`INSERT INTO "transactions" (id, account_id, destination_account_id, amount, status, transaction_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		trx.ID, trx.AccountID, trx.DestinationAccountID, trx.Amount, trx.Status, trx.TransactionType, trx.CreatedAt, trx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

func (r *PgxTransactionRepository) Update(trx *transaction.Transaction) (*transaction.Transaction, error) {
	_, err := r.db.Query(
		context.Background(),
		`UPDATE "transactions" SET status=$1, updated_at=$2 WHERE id=$7`,
		trx.Status, trx.UpdatedAt, trx.ID,
	)
	if err != nil {
		return nil, err
	}
	return trx, nil
}
