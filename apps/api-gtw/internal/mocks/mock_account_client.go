package mocks

import (
	"context"

	accountpb "github.com/buemura/minibank/api-gtw/proto/account/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAccountServiceClient struct {
	mock.Mock
}

func (m *MockAccountServiceClient) CreateAccount(ctx context.Context, in *accountpb.CreateAccountRequest, opts ...grpc.CallOption) (*accountpb.CreateAccountResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountpb.CreateAccountResponse), args.Error(1)
}

func (m *MockAccountServiceClient) GetAccount(ctx context.Context, in *accountpb.GetAccountRequest, opts ...grpc.CallOption) (*accountpb.GetAccountResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountpb.GetAccountResponse), args.Error(1)
}

func (m *MockAccountServiceClient) GetAccountByUserId(ctx context.Context, in *accountpb.GetAccountByUserIdRequest, opts ...grpc.CallOption) (*accountpb.GetAccountByUserIdResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountpb.GetAccountByUserIdResponse), args.Error(1)
}

func (m *MockAccountServiceClient) GetAccountByNumber(ctx context.Context, in *accountpb.GetAccountByNumberRequest, opts ...grpc.CallOption) (*accountpb.GetAccountByNumberResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountpb.GetAccountByNumberResponse), args.Error(1)
}

func (m *MockAccountServiceClient) GetBalance(ctx context.Context, in *accountpb.GetBalanceRequest, opts ...grpc.CallOption) (*accountpb.GetBalanceResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountpb.GetBalanceResponse), args.Error(1)
}

func (m *MockAccountServiceClient) DebitAccount(ctx context.Context, in *accountpb.DebitAccountRequest, opts ...grpc.CallOption) (*accountpb.DebitAccountResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountpb.DebitAccountResponse), args.Error(1)
}

func (m *MockAccountServiceClient) CreditAccount(ctx context.Context, in *accountpb.CreditAccountRequest, opts ...grpc.CallOption) (*accountpb.CreditAccountResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountpb.CreditAccountResponse), args.Error(1)
}
