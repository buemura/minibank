package mocks

import (
	"context"

	"github.com/buemura/minibank/svc-transaction/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockIdempotencyRepository struct {
	mock.Mock
}

func (m *MockIdempotencyRepository) Create(ctx context.Context, key *domain.IdempotencyKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockIdempotencyRepository) GetByKey(ctx context.Context, key string) (*domain.IdempotencyKey, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IdempotencyKey), args.Error(1)
}

func (m *MockIdempotencyRepository) UpdateResponse(ctx context.Context, key string, status int, body []byte, txID string) error {
	args := m.Called(ctx, key, status, body, txID)
	return args.Error(0)
}
