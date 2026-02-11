package repository

import (
	"context"

	"github.com/buemura/minibank/account-service/internal/domain"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Account, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
	UpdateBalance(ctx context.Context, id string, newBalance decimal.Decimal) error
	DebitAccount(ctx context.Context, id string, amount decimal.Decimal, txID string) (decimal.Decimal, error)
	CreditAccount(ctx context.Context, id string, amount decimal.Decimal, txID string) (decimal.Decimal, error)
}
