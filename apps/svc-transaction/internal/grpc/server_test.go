package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/buemura/minibank/svc-transaction/internal/domain"
	grpcmocks "github.com/buemura/minibank/svc-transaction/internal/grpc/mocks"
	"github.com/buemura/minibank/svc-transaction/internal/repository/mocks"
	"github.com/buemura/minibank/svc-transaction/internal/service"
	accountpb "github.com/buemura/minibank/svc-transaction/proto/account/v1"
	pb "github.com/buemura/minibank/svc-transaction/proto/transaction/v1"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newTestServer() (*TransactionServer, *mocks.MockTransactionRepository, *mocks.MockIdempotencyRepository, *grpcmocks.MockAccountServiceClient) {
	txRepo := new(mocks.MockTransactionRepository)
	idemRepo := new(mocks.MockIdempotencyRepository)
	accountClient := new(grpcmocks.MockAccountServiceClient)
	svc := service.NewTransactionService(txRepo, idemRepo, accountClient)
	server := NewTransactionServer(svc)
	return server, txRepo, idemRepo, accountClient
}

func TestTransfer_gRPC_Success(t *testing.T) {
	server, txRepo, idemRepo, accountClient := newTestServer()

	idemRepo.On("GetByKey", mock.Anything, "key-1").Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(&accountpb.GetAccountByNumberResponse{
		Account: &accountpb.Account{Id: "dest-1", AccountNumber: "9876543210"},
	}, nil)
	accountClient.On("DebitAccount", mock.Anything, mock.Anything).Return(&accountpb.DebitAccountResponse{Success: true, NewBalance: "900"}, nil)
	accountClient.On("CreditAccount", mock.Anything, mock.Anything).Return(&accountpb.CreditAccountResponse{Success: true, NewBalance: "600"}, nil)
	txRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := server.Transfer(context.Background(), &pb.TransferRequest{
		IdempotencyKey:           "key-1",
		SourceAccountId:          "src-1",
		DestinationAccountNumber: "9876543210",
		Amount:                   "100",
		Description:              "test",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Transaction)
	assert.Equal(t, pb.TransactionType_TRANSACTION_TYPE_TRANSFER, resp.Transaction.Type)
	assert.Equal(t, pb.TransactionStatus_TRANSACTION_STATUS_COMPLETED, resp.Transaction.Status)
}

func TestTransfer_gRPC_BusinessError(t *testing.T) {
	server, _, idemRepo, _ := newTestServer()

	idemRepo.On("GetByKey", mock.Anything, "key-1").Return(&domain.IdempotencyKey{
		Key:         "key-1",
		RequestHash: "different-hash",
	}, nil)

	resp, err := server.Transfer(context.Background(), &pb.TransferRequest{
		IdempotencyKey:           "key-1",
		SourceAccountId:          "src-1",
		DestinationAccountNumber: "9876543210",
		Amount:                   "100",
		Description:              "test",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "IDEMPOTENCY_MISMATCH", resp.ErrorCode)
}

func TestTransfer_gRPC_InternalError(t *testing.T) {
	server, _, idemRepo, _ := newTestServer()

	idemRepo.On("GetByKey", mock.Anything, "key-1").Return(nil, errors.New("db error"))

	_, err := server.Transfer(context.Background(), &pb.TransferRequest{
		IdempotencyKey:           "key-1",
		SourceAccountId:          "src-1",
		DestinationAccountNumber: "9876543210",
		Amount:                   "100",
		Description:              "test",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestDeposit_gRPC_Success(t *testing.T) {
	server, txRepo, idemRepo, accountClient := newTestServer()

	idemRepo.On("GetByKey", mock.Anything, "key-1").Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("CreditAccount", mock.Anything, mock.Anything).Return(&accountpb.CreditAccountResponse{Success: true, NewBalance: "1500"}, nil)
	txRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := server.Deposit(context.Background(), &pb.DepositRequest{
		IdempotencyKey: "key-1",
		AccountId:      "acc-1",
		Amount:         "500",
		Description:    "test deposit",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Transaction)
	assert.Equal(t, pb.TransactionType_TRANSACTION_TYPE_DEPOSIT, resp.Transaction.Type)
}

func TestDeposit_gRPC_InternalError(t *testing.T) {
	server, _, idemRepo, _ := newTestServer()

	idemRepo.On("GetByKey", mock.Anything, "key-1").Return(nil, errors.New("db error"))

	_, err := server.Deposit(context.Background(), &pb.DepositRequest{
		IdempotencyKey: "key-1",
		AccountId:      "acc-1",
		Amount:         "500",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGetTransaction_gRPC_Success(t *testing.T) {
	server, txRepo, _, _ := newTestServer()

	now := time.Now()
	tx := &domain.Transaction{
		ID:          "tx-1",
		Type:        domain.TransactionTypeTransfer,
		Status:      domain.TransactionStatusCompleted,
		Amount:      decimal.RequireFromString("100"),
		Currency:    "BRL",
		ProcessedAt: &now,
	}
	txRepo.On("GetByID", mock.Anything, "tx-1").Return(tx, nil)

	resp, err := server.GetTransaction(context.Background(), &pb.GetTransactionRequest{TransactionId: "tx-1"})

	require.NoError(t, err)
	assert.Equal(t, "tx-1", resp.Transaction.Id)
	assert.Equal(t, "100", resp.Transaction.Amount)
	txRepo.AssertExpectations(t)
}

func TestGetTransaction_gRPC_NotFound(t *testing.T) {
	server, txRepo, _, _ := newTestServer()

	txRepo.On("GetByID", mock.Anything, "tx-1").Return(nil, nil)

	_, err := server.GetTransaction(context.Background(), &pb.GetTransactionRequest{TransactionId: "tx-1"})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	txRepo.AssertExpectations(t)
}

func TestGetTransactionHistory_gRPC_Success(t *testing.T) {
	server, txRepo, _, _ := newTestServer()

	txs := []*domain.Transaction{
		{ID: "tx-1", Type: domain.TransactionTypeTransfer, Amount: decimal.RequireFromString("100")},
		{ID: "tx-2", Type: domain.TransactionTypeDeposit, Amount: decimal.RequireFromString("200")},
	}
	txRepo.On("GetByAccountID", mock.Anything, "acc-1", 1, 20).Return(txs, 2, nil)

	resp, err := server.GetTransactionHistory(context.Background(), &pb.GetTransactionHistoryRequest{
		AccountId: "acc-1",
		Page:      1,
		PageSize:  20,
	})

	require.NoError(t, err)
	assert.Len(t, resp.Transactions, 2)
	assert.Equal(t, int32(2), resp.TotalCount)
	txRepo.AssertExpectations(t)
}

func TestGetStatement_gRPC_Success(t *testing.T) {
	server, txRepo, _, accountClient := newTestServer()

	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()

	txs := []*domain.Transaction{
		{ID: "tx-1", Amount: decimal.RequireFromString("100")},
	}
	txRepo.On("GetByAccountIDAndDateRange", mock.Anything, "acc-1", mock.Anything, mock.Anything, 1, 50).Return(txs, 1, nil)
	accountClient.On("GetBalance", mock.Anything, mock.Anything).Return(&accountpb.GetBalanceResponse{Balance: "1000.00"}, nil)

	resp, err := server.GetStatement(context.Background(), &pb.GetStatementRequest{
		AccountId: "acc-1",
		StartDate: timestamppb.New(start),
		EndDate:   timestamppb.New(end),
		Page:      1,
		PageSize:  50,
	})

	require.NoError(t, err)
	assert.Equal(t, "acc-1", resp.AccountId)
	assert.Equal(t, "1000.00", resp.OpeningBalance)
	assert.Equal(t, "1000.00", resp.ClosingBalance)
	assert.Len(t, resp.Transactions, 1)
	txRepo.AssertExpectations(t)
	accountClient.AssertExpectations(t)
}

func TestDomainTypeToProto(t *testing.T) {
	tests := []struct {
		input    domain.TransactionType
		expected pb.TransactionType
	}{
		{domain.TransactionTypeTransfer, pb.TransactionType_TRANSACTION_TYPE_TRANSFER},
		{domain.TransactionTypeDeposit, pb.TransactionType_TRANSACTION_TYPE_DEPOSIT},
		{domain.TransactionTypeWithdrawal, pb.TransactionType_TRANSACTION_TYPE_WITHDRAWAL},
		{domain.TransactionTypeFee, pb.TransactionType_TRANSACTION_TYPE_FEE},
		{"unknown", pb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED},
	}

	for _, tc := range tests {
		result := domainTypeToProto(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestDomainStatusToProto(t *testing.T) {
	tests := []struct {
		input    domain.TransactionStatus
		expected pb.TransactionStatus
	}{
		{domain.TransactionStatusPending, pb.TransactionStatus_TRANSACTION_STATUS_PENDING},
		{domain.TransactionStatusCompleted, pb.TransactionStatus_TRANSACTION_STATUS_COMPLETED},
		{domain.TransactionStatusFailed, pb.TransactionStatus_TRANSACTION_STATUS_FAILED},
		{domain.TransactionStatusReversed, pb.TransactionStatus_TRANSACTION_STATUS_REVERSED},
		{"unknown", pb.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED},
	}

	for _, tc := range tests {
		result := domainStatusToProto(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}
