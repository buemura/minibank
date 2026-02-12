package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/buemura/minibank/svc-transaction/internal/domain"
	grpcmocks "github.com/buemura/minibank/svc-transaction/internal/grpc/mocks"
	"github.com/buemura/minibank/svc-transaction/internal/repository/mocks"
	accountpb "github.com/buemura/minibank/svc-transaction/proto/account/v1"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestService() (*TransactionService, *mocks.MockTransactionRepository, *mocks.MockIdempotencyRepository, *grpcmocks.MockAccountServiceClient) {
	txRepo := new(mocks.MockTransactionRepository)
	idemRepo := new(mocks.MockIdempotencyRepository)
	accountClient := new(grpcmocks.MockAccountServiceClient)
	svc := NewTransactionService(txRepo, idemRepo, accountClient)
	return svc, txRepo, idemRepo, accountClient
}

func defaultTransferInput() TransferInput {
	return TransferInput{
		IdempotencyKey:           "idem-key-1",
		SourceAccountID:          "src-acc-1",
		DestinationAccountNumber: "9876543210",
		Amount:                   "100.00",
		Description:              "Test transfer",
	}
}

func defaultDepositInput() DepositInput {
	return DepositInput{
		IdempotencyKey: "idem-key-2",
		AccountID:      "acc-1",
		Amount:         "500.00",
		Description:    "Test deposit",
	}
}

func defaultWithdrawInput() WithdrawInput {
	return WithdrawInput{
		IdempotencyKey: "idem-key-3",
		AccountID:      "acc-1",
		Amount:         "200.00",
		Description:    "Test withdrawal",
	}
}

// --- Transfer Tests ---

func TestTransfer_Success(t *testing.T) {
	svc, txRepo, idemRepo, accountClient := newTestService()
	input := defaultTransferInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(&accountpb.GetAccountByNumberResponse{
		Account: &accountpb.Account{Id: "dest-acc-1", AccountNumber: "9876543210"},
	}, nil)
	accountClient.On("GetAccount", mock.Anything, mock.Anything).Return(&accountpb.GetAccountResponse{
		Account: &accountpb.Account{Id: "src-acc-1", AccountNumber: "1234567890"},
	}, nil)
	accountClient.On("DebitAccount", mock.Anything, mock.Anything).Return(&accountpb.DebitAccountResponse{
		Success:    true,
		NewBalance: "900.00",
	}, nil)
	accountClient.On("CreditAccount", mock.Anything, mock.Anything).Return(&accountpb.CreditAccountResponse{
		Success:    true,
		NewBalance: "600.00",
	}, nil)

	txRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Transaction)
	assert.Equal(t, domain.TransactionTypeTransfer, result.Transaction.Type)
	assert.Equal(t, domain.TransactionStatusCompleted, result.Transaction.Status)
	assert.Equal(t, "1234567890", result.Transaction.SourceAccountNumber)
	assert.True(t, result.Transaction.Amount.Equal(decimal.RequireFromString("100.00")))
	txRepo.AssertExpectations(t)
	idemRepo.AssertExpectations(t)
	accountClient.AssertExpectations(t)
}

func TestTransfer_InvalidAmount_Zero(t *testing.T) {
	svc, _, _, _ := newTestService()
	input := defaultTransferInput()
	input.Amount = "0"

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "INVALID_AMOUNT", result.ErrorCode)
}

func TestTransfer_InvalidAmount_Negative(t *testing.T) {
	svc, _, _, _ := newTestService()
	input := defaultTransferInput()
	input.Amount = "-50"

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "INVALID_AMOUNT", result.ErrorCode)
}

func TestTransfer_InvalidAmount_NonNumeric(t *testing.T) {
	svc, _, _, _ := newTestService()
	input := defaultTransferInput()
	input.Amount = "not-a-number"

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "INVALID_AMOUNT", result.ErrorCode)
}

func TestTransfer_IdempotencyKey_SameRequest(t *testing.T) {
	svc, txRepo, idemRepo, _ := newTestService()
	input := defaultTransferInput()

	requestHash := calculateTransferRequestHash(input)
	existingTx := &domain.Transaction{
		ID:     "tx-1",
		Status: domain.TransactionStatusCompleted,
		Amount: decimal.RequireFromString("100.00"),
	}

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(&domain.IdempotencyKey{
		Key:           input.IdempotencyKey,
		RequestHash:   requestHash,
		TransactionID: "tx-1",
	}, nil)
	txRepo.On("GetByID", mock.Anything, "tx-1").Return(existingTx, nil)

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "tx-1", result.Transaction.ID)
	idemRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestTransfer_IdempotencyKey_Mismatch(t *testing.T) {
	svc, _, idemRepo, _ := newTestService()
	input := defaultTransferInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(&domain.IdempotencyKey{
		Key:         input.IdempotencyKey,
		RequestHash: "different-hash",
	}, nil)

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "IDEMPOTENCY_MISMATCH", result.ErrorCode)
	idemRepo.AssertExpectations(t)
}

func TestTransfer_DestinationNotFound(t *testing.T) {
	svc, _, idemRepo, accountClient := newTestService()
	input := defaultTransferInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "INVALID_DESTINATION", result.ErrorCode)
	idemRepo.AssertExpectations(t)
	accountClient.AssertExpectations(t)
}

func TestTransfer_DebitFails(t *testing.T) {
	svc, _, idemRepo, accountClient := newTestService()
	input := defaultTransferInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(&accountpb.GetAccountByNumberResponse{
		Account: &accountpb.Account{Id: "dest-acc-1"},
	}, nil)
	accountClient.On("GetAccount", mock.Anything, mock.Anything).Return(&accountpb.GetAccountResponse{
		Account: &accountpb.Account{Id: "src-acc-1", AccountNumber: "1234567890"},
	}, nil)
	accountClient.On("DebitAccount", mock.Anything, mock.Anything).Return(&accountpb.DebitAccountResponse{
		Success:      false,
		ErrorMessage: "insufficient funds",
	}, nil)

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "INSUFFICIENT_FUNDS", result.ErrorCode)
	accountClient.AssertExpectations(t)
}

func TestTransfer_CreditFails_Rollback(t *testing.T) {
	svc, _, idemRepo, accountClient := newTestService()
	input := defaultTransferInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(&accountpb.GetAccountByNumberResponse{
		Account: &accountpb.Account{Id: "dest-acc-1"},
	}, nil)
	accountClient.On("GetAccount", mock.Anything, mock.Anything).Return(&accountpb.GetAccountResponse{
		Account: &accountpb.Account{Id: "src-acc-1", AccountNumber: "1234567890"},
	}, nil)
	accountClient.On("DebitAccount", mock.Anything, mock.Anything).Return(&accountpb.DebitAccountResponse{
		Success:    true,
		NewBalance: "900.00",
	}, nil)
	// First CreditAccount (to destination) fails
	accountClient.On("CreditAccount", mock.Anything, mock.MatchedBy(func(req *accountpb.CreditAccountRequest) bool {
		return req.AccountId == "dest-acc-1"
	})).Return(&accountpb.CreditAccountResponse{Success: false}, nil)
	// Rollback CreditAccount (back to source)
	accountClient.On("CreditAccount", mock.Anything, mock.MatchedBy(func(req *accountpb.CreditAccountRequest) bool {
		return req.AccountId == "src-acc-1"
	})).Return(&accountpb.CreditAccountResponse{Success: true}, nil)

	result, err := svc.Transfer(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "TRANSFER_FAILED", result.ErrorCode)
	accountClient.AssertExpectations(t)
}

func TestTransfer_IdempotencyRepoError(t *testing.T) {
	svc, _, idemRepo, _ := newTestService()
	input := defaultTransferInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, errors.New("db error"))

	_, err := svc.Transfer(context.Background(), input)

	assert.EqualError(t, err, "db error")
	idemRepo.AssertExpectations(t)
}

// --- Deposit Tests ---

func TestDeposit_Success(t *testing.T) {
	svc, txRepo, idemRepo, accountClient := newTestService()
	input := defaultDepositInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("CreditAccount", mock.Anything, mock.Anything).Return(&accountpb.CreditAccountResponse{
		Success:    true,
		NewBalance: "1500.00",
	}, nil)

	txRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	result, err := svc.Deposit(context.Background(), input)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Transaction)
	assert.Equal(t, domain.TransactionTypeDeposit, result.Transaction.Type)
	assert.Equal(t, domain.TransactionStatusCompleted, result.Transaction.Status)
	txRepo.AssertExpectations(t)
	idemRepo.AssertExpectations(t)
	accountClient.AssertExpectations(t)
}

func TestDeposit_InvalidAmount(t *testing.T) {
	svc, _, _, _ := newTestService()
	input := defaultDepositInput()
	input.Amount = "0"

	result, err := svc.Deposit(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "INVALID_AMOUNT", result.ErrorCode)
}

func TestDeposit_IdempotencyMismatch(t *testing.T) {
	svc, _, idemRepo, _ := newTestService()
	input := defaultDepositInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(&domain.IdempotencyKey{
		Key:         input.IdempotencyKey,
		RequestHash: "different-hash",
	}, nil)

	result, err := svc.Deposit(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "IDEMPOTENCY_MISMATCH", result.ErrorCode)
	idemRepo.AssertExpectations(t)
}

func TestDeposit_IdempotencyReplay(t *testing.T) {
	svc, txRepo, idemRepo, _ := newTestService()
	input := defaultDepositInput()

	requestHash := calculateDepositRequestHash(input)
	existingTx := &domain.Transaction{
		ID:     "tx-1",
		Status: domain.TransactionStatusCompleted,
	}

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(&domain.IdempotencyKey{
		Key:           input.IdempotencyKey,
		RequestHash:   requestHash,
		TransactionID: "tx-1",
	}, nil)
	txRepo.On("GetByID", mock.Anything, "tx-1").Return(existingTx, nil)

	result, err := svc.Deposit(context.Background(), input)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "tx-1", result.Transaction.ID)
	idemRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestDeposit_CreditFails_gRPCError(t *testing.T) {
	svc, _, idemRepo, accountClient := newTestService()
	input := defaultDepositInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("CreditAccount", mock.Anything, mock.Anything).Return(nil, errors.New("gRPC error"))

	result, err := svc.Deposit(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "DEPOSIT_FAILED", result.ErrorCode)
	accountClient.AssertExpectations(t)
}

func TestDeposit_CreditFails_NotSuccess(t *testing.T) {
	svc, _, idemRepo, accountClient := newTestService()
	input := defaultDepositInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("CreditAccount", mock.Anything, mock.Anything).Return(&accountpb.CreditAccountResponse{
		Success: false,
	}, nil)

	result, err := svc.Deposit(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "DEPOSIT_FAILED", result.ErrorCode)
	accountClient.AssertExpectations(t)
}

// --- GetTransaction Tests ---

func TestGetTransaction_Success(t *testing.T) {
	svc, txRepo, _, _ := newTestService()

	tx := &domain.Transaction{
		ID:     "tx-1",
		Type:   domain.TransactionTypeTransfer,
		Status: domain.TransactionStatusCompleted,
		Amount: decimal.RequireFromString("100.00"),
	}
	txRepo.On("GetByID", mock.Anything, "tx-1").Return(tx, nil)

	result, err := svc.GetTransaction(context.Background(), "tx-1")

	require.NoError(t, err)
	assert.Equal(t, "tx-1", result.ID)
	txRepo.AssertExpectations(t)
}

func TestGetTransaction_NotFound(t *testing.T) {
	svc, txRepo, _, _ := newTestService()

	txRepo.On("GetByID", mock.Anything, "tx-1").Return(nil, nil)

	_, err := svc.GetTransaction(context.Background(), "tx-1")

	assert.ErrorIs(t, err, ErrTransactionNotFound)
	txRepo.AssertExpectations(t)
}

// --- GetTransactionHistory Tests ---

func TestGetTransactionHistory_Success(t *testing.T) {
	svc, txRepo, _, _ := newTestService()

	txs := []*domain.Transaction{
		{ID: "tx-1", Amount: decimal.RequireFromString("100")},
		{ID: "tx-2", Amount: decimal.RequireFromString("200")},
	}
	txRepo.On("GetByAccountID", mock.Anything, "acc-1", 1, 20).Return(txs, 2, nil)

	result, total, err := svc.GetTransactionHistory(context.Background(), "acc-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
	txRepo.AssertExpectations(t)
}

func TestGetTransactionHistory_DefaultPagination(t *testing.T) {
	svc, txRepo, _, _ := newTestService()

	txRepo.On("GetByAccountID", mock.Anything, "acc-1", 1, 20).Return([]*domain.Transaction{}, 0, nil)

	_, _, err := svc.GetTransactionHistory(context.Background(), "acc-1", 0, 0)

	require.NoError(t, err)
	txRepo.AssertExpectations(t)
}

func TestGetTransactionHistory_MaxPageSize(t *testing.T) {
	svc, txRepo, _, _ := newTestService()

	txRepo.On("GetByAccountID", mock.Anything, "acc-1", 1, 20).Return([]*domain.Transaction{}, 0, nil)

	_, _, err := svc.GetTransactionHistory(context.Background(), "acc-1", 1, 200)

	require.NoError(t, err)
	txRepo.AssertExpectations(t)
}

// --- GetStatement Tests ---

func TestGetStatement_Success(t *testing.T) {
	svc, txRepo, _, accountClient := newTestService()

	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()

	txs := []*domain.Transaction{
		{ID: "tx-1", Amount: decimal.RequireFromString("100")},
	}
	txRepo.On("GetByAccountIDAndDateRange", mock.Anything, "acc-1", start, end, 1, 50).Return(txs, 1, nil)
	accountClient.On("GetBalance", mock.Anything, mock.Anything).Return(&accountpb.GetBalanceResponse{
		Balance:  "1000.00",
		Currency: "BRL",
	}, nil)

	transactions, total, opening, closing, err := svc.GetStatement(context.Background(), "acc-1", start, end, 1, 50)

	require.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, 1, total)
	assert.Equal(t, "1000.00", opening)
	assert.Equal(t, "1000.00", closing)
	txRepo.AssertExpectations(t)
	accountClient.AssertExpectations(t)
}

func TestGetStatement_DefaultPagination(t *testing.T) {
	svc, txRepo, _, accountClient := newTestService()

	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()

	txRepo.On("GetByAccountIDAndDateRange", mock.Anything, "acc-1", start, end, 1, 50).Return([]*domain.Transaction{}, 0, nil)
	accountClient.On("GetBalance", mock.Anything, mock.Anything).Return(&accountpb.GetBalanceResponse{
		Balance: "0.00",
	}, nil)

	_, _, _, _, err := svc.GetStatement(context.Background(), "acc-1", start, end, 0, 0)

	require.NoError(t, err)
	txRepo.AssertExpectations(t)
}

func TestGetStatement_RepoError(t *testing.T) {
	svc, txRepo, _, _ := newTestService()

	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()

	txRepo.On("GetByAccountIDAndDateRange", mock.Anything, "acc-1", start, end, 1, 50).Return(nil, 0, errors.New("db error"))

	_, _, _, _, err := svc.GetStatement(context.Background(), "acc-1", start, end, 1, 50)

	assert.EqualError(t, err, "db error")
	txRepo.AssertExpectations(t)
}

func TestGetStatement_BalanceError(t *testing.T) {
	svc, txRepo, _, accountClient := newTestService()

	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()

	txRepo.On("GetByAccountIDAndDateRange", mock.Anything, "acc-1", start, end, 1, 50).Return([]*domain.Transaction{}, 0, nil)
	accountClient.On("GetBalance", mock.Anything, mock.Anything).Return(nil, errors.New("gRPC error"))

	_, _, _, _, err := svc.GetStatement(context.Background(), "acc-1", start, end, 1, 50)

	assert.EqualError(t, err, "gRPC error")
	txRepo.AssertExpectations(t)
	accountClient.AssertExpectations(t)
}

// Helper functions to calculate request hashes (same logic as the service)

func calculateTransferRequestHash(input TransferInput) string {
	data := map[string]string{
		"source_account_id":          input.SourceAccountID,
		"destination_account_number": input.DestinationAccountNumber,
		"amount":                     input.Amount,
		"description":                input.Description,
	}
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

func calculateDepositRequestHash(input DepositInput) string {
	data := map[string]string{
		"account_id":  input.AccountID,
		"amount":      input.Amount,
		"description": input.Description,
	}
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

// --- Withdraw Tests ---

func TestWithdraw_Success(t *testing.T) {
	svc, txRepo, idemRepo, accountClient := newTestService()
	input := defaultWithdrawInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("DebitAccount", mock.Anything, mock.Anything).Return(&accountpb.DebitAccountResponse{
		Success:    true,
		NewBalance: "800.00",
	}, nil)

	txRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	result, err := svc.Withdraw(context.Background(), input)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Transaction)
	assert.Equal(t, domain.TransactionTypeWithdrawal, result.Transaction.Type)
	assert.Equal(t, domain.TransactionStatusCompleted, result.Transaction.Status)
	assert.True(t, result.Transaction.Amount.Equal(decimal.RequireFromString("200.00")))
	assert.Equal(t, input.AccountID, result.Transaction.SourceAccountID)
	txRepo.AssertExpectations(t)
	idemRepo.AssertExpectations(t)
	accountClient.AssertExpectations(t)
}

func TestWithdraw_InvalidAmount(t *testing.T) {
	svc, _, _, _ := newTestService()
	input := defaultWithdrawInput()
	input.Amount = "0"

	result, err := svc.Withdraw(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "INVALID_AMOUNT", result.ErrorCode)
}

func TestWithdraw_IdempotencyMismatch(t *testing.T) {
	svc, _, idemRepo, _ := newTestService()
	input := defaultWithdrawInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(&domain.IdempotencyKey{
		Key:         input.IdempotencyKey,
		RequestHash: "different-hash",
	}, nil)

	result, err := svc.Withdraw(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "IDEMPOTENCY_MISMATCH", result.ErrorCode)
	idemRepo.AssertExpectations(t)
}

func TestWithdraw_IdempotencyReplay(t *testing.T) {
	svc, txRepo, idemRepo, _ := newTestService()
	input := defaultWithdrawInput()

	requestHash := calculateWithdrawRequestHash(input)
	existingTx := &domain.Transaction{
		ID:     "tx-1",
		Status: domain.TransactionStatusCompleted,
	}

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(&domain.IdempotencyKey{
		Key:           input.IdempotencyKey,
		RequestHash:   requestHash,
		TransactionID: "tx-1",
	}, nil)
	txRepo.On("GetByID", mock.Anything, "tx-1").Return(existingTx, nil)

	result, err := svc.Withdraw(context.Background(), input)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "tx-1", result.Transaction.ID)
	idemRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestWithdraw_DebitFails_gRPCError(t *testing.T) {
	svc, _, idemRepo, accountClient := newTestService()
	input := defaultWithdrawInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("DebitAccount", mock.Anything, mock.Anything).Return(nil, errors.New("gRPC error"))

	result, err := svc.Withdraw(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "WITHDRAWAL_FAILED", result.ErrorCode)
	accountClient.AssertExpectations(t)
}

func TestWithdraw_DebitFails_NotSuccess(t *testing.T) {
	svc, _, idemRepo, accountClient := newTestService()
	input := defaultWithdrawInput()

	idemRepo.On("GetByKey", mock.Anything, input.IdempotencyKey).Return(nil, nil)
	idemRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	idemRepo.On("UpdateResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	accountClient.On("DebitAccount", mock.Anything, mock.Anything).Return(&accountpb.DebitAccountResponse{
		Success:      false,
		ErrorMessage: "insufficient funds",
	}, nil)

	result, err := svc.Withdraw(context.Background(), input)

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "WITHDRAWAL_FAILED", result.ErrorCode)
	assert.Equal(t, "insufficient funds", result.ErrorMessage)
	accountClient.AssertExpectations(t)
}

func calculateWithdrawRequestHash(input WithdrawInput) string {
	data := map[string]string{
		"account_id":  input.AccountID,
		"amount":      input.Amount,
		"description": input.Description,
	}
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}
