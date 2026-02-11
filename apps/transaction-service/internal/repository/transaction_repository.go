package repository

import (
	"context"
	"time"

	"github.com/buemura/minibank/transaction-service/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error)
	GetByAccountID(ctx context.Context, accountID string, page, pageSize int) ([]*domain.Transaction, int, error)
	GetByAccountIDAndDateRange(ctx context.Context, accountID string, startDate, endDate time.Time, page, pageSize int) ([]*domain.Transaction, int, error)
	UpdateStatus(ctx context.Context, id string, status domain.TransactionStatus) error
}

type IdempotencyRepository interface {
	Create(ctx context.Context, key *domain.IdempotencyKey) error
	GetByKey(ctx context.Context, key string) (*domain.IdempotencyKey, error)
	UpdateResponse(ctx context.Context, key string, status int, body []byte, txID string) error
}
