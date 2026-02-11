package postgres

import (
	"context"
	"errors"

	"github.com/buemura/minibank/svc-transaction/internal/domain"
	"github.com/buemura/minibank/svc-transaction/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type idempotencyRepository struct {
	db *pgxpool.Pool
}

func NewIdempotencyRepository(db *pgxpool.Pool) repository.IdempotencyRepository {
	return &idempotencyRepository{db: db}
}

func (r *idempotencyRepository) Create(ctx context.Context, key *domain.IdempotencyKey) error {
	query := `
		INSERT INTO idempotency_keys (key, request_hash, response_status, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		key.Key,
		key.RequestHash,
		key.ResponseStatus,
		key.ExpiresAt,
	).Scan(&key.ID, &key.CreatedAt)

	return err
}

func (r *idempotencyRepository) GetByKey(ctx context.Context, key string) (*domain.IdempotencyKey, error) {
	query := `
		SELECT id, key, request_hash, response_status, response_body, transaction_id, expires_at, created_at
		FROM idempotency_keys
		WHERE key = $1
	`

	iKey := &domain.IdempotencyKey{}
	var txID *string

	err := r.db.QueryRow(ctx, query, key).Scan(
		&iKey.ID,
		&iKey.Key,
		&iKey.RequestHash,
		&iKey.ResponseStatus,
		&iKey.ResponseBody,
		&txID,
		&iKey.ExpiresAt,
		&iKey.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if txID != nil {
		iKey.TransactionID = *txID
	}

	return iKey, nil
}

func (r *idempotencyRepository) UpdateResponse(ctx context.Context, key string, status int, body []byte, txID string) error {
	query := `
		UPDATE idempotency_keys
		SET response_status = $1, response_body = $2, transaction_id = $3
		WHERE key = $4
	`

	var txIDPtr *string
	if txID != "" {
		txIDPtr = &txID
	}

	_, err := r.db.Exec(ctx, query, status, body, txIDPtr, key)
	return err
}
