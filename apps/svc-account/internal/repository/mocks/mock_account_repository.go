package mocks

import (
	"context"

	"github.com/buemura/minibank/svc-account/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	args := m.Called(ctx, accountNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateBalance(ctx context.Context, id string, newBalance decimal.Decimal) error {
	args := m.Called(ctx, id, newBalance)
	return args.Error(0)
}

func (m *MockAccountRepository) DebitAccount(ctx context.Context, id string, amount decimal.Decimal, txID string) (decimal.Decimal, error) {
	args := m.Called(ctx, id, amount, txID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockAccountRepository) CreditAccount(ctx context.Context, id string, amount decimal.Decimal, txID string) (decimal.Decimal, error) {
	args := m.Called(ctx, id, amount, txID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}
