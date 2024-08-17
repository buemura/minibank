package database

import (
	"database/sql"

	"github.com/buemura/minibank/svc-transaction/internal/core/domain/common"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
)

type SqlTransactionRepository struct {
	db *sql.DB
}

func NewSqlTransactionRepository() *SqlTransactionRepository {
	return &SqlTransactionRepository{db: Conn}
}

func (r *SqlTransactionRepository) FindById(id string) (*transaction.Transaction, error) {
	var trx *transaction.Transaction
	err := r.db.QueryRow(
		`SELECT id, account_id, destination_account_id, amount, status, transaction_type, created_at, updated_at
		FROM "transactions"
		WHERE id = $1`,
		id,
	).Scan(trx)
	if err != nil {
		return nil, err
	}

	if trx == nil {
		return nil, transaction.ErrTransactionNotFound
	}

	return trx, nil
}

func (r *SqlTransactionRepository) FindByAccountId(in *transaction.GetTransactionListIn) (*transaction.GetTransactionListOut, error) {
	limit := in.Items
	offset := (in.Page - 1) * in.Items

	rows, err := r.db.Query(
		`SELECT id, account_id, destination_account_id, amount, status, transaction_type, created_at, updated_at
		FROM "transactions"
		WHERE account_id = $1 OR destination_account_id = $1
		LIMIT $2 OFFSET $3`,
		in.AccountID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trxs []*transaction.Transaction

	for rows.Next() {
		var trx *transaction.Transaction
		if err := rows.Scan(trx); err != nil {
			return nil, err
		}
		trxs = append(trxs, trx)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var totalCount int
	err = r.db.QueryRow(`SELECT count(id) as total_count FROM "transactions"`).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	return &transaction.GetTransactionListOut{
		Data: trxs,
		Meta: common.PaginationMetaOut{
			Page:       in.Page,
			Items:      in.Items,
			TotalPages: int(totalCount/in.Items) + 1,
			TotalItems: totalCount,
		},
	}, nil
}

func (r *SqlTransactionRepository) Create(trx *transaction.Transaction) (*transaction.Transaction, error) {
	_, err := r.db.Exec(
		`INSERT INTO "transactions" (id, account_id, destination_account_id, amount, status, transaction_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		trx.ID, trx.AccountID, trx.DestinationAccountID, trx.Amount, trx.Status, trx.TransactionType, trx.CreatedAt, trx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

func (r *SqlTransactionRepository) Update(trx *transaction.Transaction) (*transaction.Transaction, error) {
	_, err := r.db.Exec(
		`UPDATE "transactions" SET status=$1, updated_at=$2 WHERE id=$3`,
		trx.Status, trx.UpdatedAt, trx.ID,
	)
	if err != nil {
		return nil, err
	}
	return trx, nil
}
