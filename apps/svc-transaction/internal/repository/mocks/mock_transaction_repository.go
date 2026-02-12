package mocks

import (
	"context"
	"time"

	"github.com/buemura/minibank/svc-transaction/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByAccountID(ctx context.Context, accountID string, page, pageSize int) ([]*domain.Transaction, int, error) {
	args := m.Called(ctx, accountID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Transaction), args.Int(1), args.Error(2)
}

func (m *MockTransactionRepository) GetByAccountIDAndDateRange(ctx context.Context, accountID string, startDate, endDate time.Time, page, pageSize int) ([]*domain.Transaction, int, error) {
	args := m.Called(ctx, accountID, startDate, endDate, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Transaction), args.Int(1), args.Error(2)
}

func (m *MockTransactionRepository) UpdateStatus(ctx context.Context, id string, status domain.TransactionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}
