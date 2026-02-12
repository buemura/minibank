package mocks

import (
	"context"

	"github.com/buemura/minibank/svc-auth/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByTokenHash(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}
