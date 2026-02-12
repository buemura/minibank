package mocks

import (
	"context"

	transactionpb "github.com/buemura/minibank/api-gtw/proto/transaction/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockTransactionServiceClient struct {
	mock.Mock
}

func (m *MockTransactionServiceClient) Transfer(ctx context.Context, in *transactionpb.TransferRequest, opts ...grpc.CallOption) (*transactionpb.TransferResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionpb.TransferResponse), args.Error(1)
}

func (m *MockTransactionServiceClient) Deposit(ctx context.Context, in *transactionpb.DepositRequest, opts ...grpc.CallOption) (*transactionpb.DepositResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionpb.DepositResponse), args.Error(1)
}

func (m *MockTransactionServiceClient) Withdraw(ctx context.Context, in *transactionpb.WithdrawRequest, opts ...grpc.CallOption) (*transactionpb.WithdrawResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionpb.WithdrawResponse), args.Error(1)
}

func (m *MockTransactionServiceClient) GetTransaction(ctx context.Context, in *transactionpb.GetTransactionRequest, opts ...grpc.CallOption) (*transactionpb.GetTransactionResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionpb.GetTransactionResponse), args.Error(1)
}

func (m *MockTransactionServiceClient) GetTransactionHistory(ctx context.Context, in *transactionpb.GetTransactionHistoryRequest, opts ...grpc.CallOption) (*transactionpb.GetTransactionHistoryResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionpb.GetTransactionHistoryResponse), args.Error(1)
}

func (m *MockTransactionServiceClient) GetStatement(ctx context.Context, in *transactionpb.GetStatementRequest, opts ...grpc.CallOption) (*transactionpb.GetStatementResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionpb.GetStatementResponse), args.Error(1)
}
