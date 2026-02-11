package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/buemura/minibank/transaction-service/internal/domain"
	"github.com/buemura/minibank/transaction-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) repository.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			idempotency_key, type, status, source_account_id, source_account_number,
			destination_account_id, destination_account_number, amount, currency,
			description, reference_id, processed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		nilIfEmpty(tx.IdempotencyKey),
		tx.Type,
		tx.Status,
		nilIfEmpty(tx.SourceAccountID),
		nilIfEmpty(tx.SourceAccountNumber),
		nilIfEmpty(tx.DestinationAccountID),
		nilIfEmpty(tx.DestinationAccountNumber),
		tx.Amount,
		tx.Currency,
		nilIfEmpty(tx.Description),
		nilIfEmpty(tx.ReferenceID),
		tx.ProcessedAt,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)

	return err
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `
		SELECT id, idempotency_key, type, status, source_account_id, source_account_number,
			   destination_account_id, destination_account_number, amount, currency,
			   description, reference_id, processed_at, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	tx := &domain.Transaction{}
	var idempotencyKey, sourceAccID, sourceAccNum, destAccID, destAccNum, desc, refID *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&tx.ID, &idempotencyKey, &tx.Type, &tx.Status, &sourceAccID, &sourceAccNum,
		&destAccID, &destAccNum, &tx.Amount, &tx.Currency,
		&desc, &refID, &tx.ProcessedAt, &tx.CreatedAt, &tx.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tx.IdempotencyKey = deref(idempotencyKey)
	tx.SourceAccountID = deref(sourceAccID)
	tx.SourceAccountNumber = deref(sourceAccNum)
	tx.DestinationAccountID = deref(destAccID)
	tx.DestinationAccountNumber = deref(destAccNum)
	tx.Description = deref(desc)
	tx.ReferenceID = deref(refID)

	return tx, nil
}

func (r *transactionRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	query := `
		SELECT id, idempotency_key, type, status, source_account_id, source_account_number,
			   destination_account_id, destination_account_number, amount, currency,
			   description, reference_id, processed_at, created_at, updated_at
		FROM transactions
		WHERE idempotency_key = $1
	`

	tx := &domain.Transaction{}
	var idempotencyKey, sourceAccID, sourceAccNum, destAccID, destAccNum, desc, refID *string

	err := r.db.QueryRow(ctx, query, key).Scan(
		&tx.ID, &idempotencyKey, &tx.Type, &tx.Status, &sourceAccID, &sourceAccNum,
		&destAccID, &destAccNum, &tx.Amount, &tx.Currency,
		&desc, &refID, &tx.ProcessedAt, &tx.CreatedAt, &tx.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tx.IdempotencyKey = deref(idempotencyKey)
	tx.SourceAccountID = deref(sourceAccID)
	tx.SourceAccountNumber = deref(sourceAccNum)
	tx.DestinationAccountID = deref(destAccID)
	tx.DestinationAccountNumber = deref(destAccNum)
	tx.Description = deref(desc)
	tx.ReferenceID = deref(refID)

	return tx, nil
}

func (r *transactionRepository) GetByAccountID(ctx context.Context, accountID string, page, pageSize int) ([]*domain.Transaction, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `
		SELECT COUNT(*)
		FROM transactions
		WHERE source_account_id = $1 OR destination_account_id = $1
	`
	var totalCount int
	if err := r.db.QueryRow(ctx, countQuery, accountID).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, idempotency_key, type, status, source_account_id, source_account_number,
			   destination_account_id, destination_account_number, amount, currency,
			   description, reference_id, processed_at, created_at, updated_at
		FROM transactions
		WHERE source_account_id = $1 OR destination_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, accountID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanTransactions(rows, totalCount)
}

func (r *transactionRepository) GetByAccountIDAndDateRange(ctx context.Context, accountID string, startDate, endDate time.Time, page, pageSize int) ([]*domain.Transaction, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `
		SELECT COUNT(*)
		FROM transactions
		WHERE (source_account_id = $1 OR destination_account_id = $1)
		  AND created_at >= $2 AND created_at <= $3
	`
	var totalCount int
	if err := r.db.QueryRow(ctx, countQuery, accountID, startDate, endDate).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, idempotency_key, type, status, source_account_id, source_account_number,
			   destination_account_id, destination_account_number, amount, currency,
			   description, reference_id, processed_at, created_at, updated_at
		FROM transactions
		WHERE (source_account_id = $1 OR destination_account_id = $1)
		  AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(ctx, query, accountID, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanTransactions(rows, totalCount)
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, id string, status domain.TransactionStatus) error {
	query := `UPDATE transactions SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *transactionRepository) scanTransactions(rows pgx.Rows, totalCount int) ([]*domain.Transaction, int, error) {
	var transactions []*domain.Transaction

	for rows.Next() {
		tx := &domain.Transaction{}
		var idempotencyKey, sourceAccID, sourceAccNum, destAccID, destAccNum, desc, refID *string

		err := rows.Scan(
			&tx.ID, &idempotencyKey, &tx.Type, &tx.Status, &sourceAccID, &sourceAccNum,
			&destAccID, &destAccNum, &tx.Amount, &tx.Currency,
			&desc, &refID, &tx.ProcessedAt, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		tx.IdempotencyKey = deref(idempotencyKey)
		tx.SourceAccountID = deref(sourceAccID)
		tx.SourceAccountNumber = deref(sourceAccNum)
		tx.DestinationAccountID = deref(destAccID)
		tx.DestinationAccountNumber = deref(destAccNum)
		tx.Description = deref(desc)
		tx.ReferenceID = deref(refID)

		transactions = append(transactions, tx)
	}

	return transactions, totalCount, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
