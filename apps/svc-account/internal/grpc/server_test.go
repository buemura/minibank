package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/buemura/minibank/svc-account/internal/domain"
	"github.com/buemura/minibank/svc-account/internal/repository/mocks"
	"github.com/buemura/minibank/svc-account/internal/service"
	pb "github.com/buemura/minibank/svc-account/proto/account/v1"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestServer() (*AccountServer, *mocks.MockAccountRepository) {
	repo := new(mocks.MockAccountRepository)
	svc := service.NewAccountService(repo)
	server := NewAccountServer(svc)
	return server, repo
}

func newActiveAccount(id, userID string) *domain.Account {
	return &domain.Account{
		ID:            id,
		UserID:        userID,
		AccountNumber: "1234567890",
		Agency:        "0001",
		AccountType:   domain.AccountTypeChecking,
		Balance:       decimal.NewFromInt(1000),
		Currency:      "BRL",
		Status:        domain.AccountStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func TestCreateAccount_gRPC_Success(t *testing.T) {
	server, repo := newTestServer()

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	resp, err := server.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		UserId:      "user-1",
		AccountType: "checking",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", resp.Account.UserId)
	assert.Equal(t, "checking", resp.Account.AccountType)
	assert.Equal(t, "0001", resp.Account.Agency)
	repo.AssertExpectations(t)
}

func TestCreateAccount_gRPC_SavingsType(t *testing.T) {
	server, repo := newTestServer()

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	resp, err := server.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		UserId:      "user-1",
		AccountType: "savings",
	})

	require.NoError(t, err)
	assert.Equal(t, "savings", resp.Account.AccountType)
	repo.AssertExpectations(t)
}

func TestCreateAccount_gRPC_InternalError(t *testing.T) {
	server, repo := newTestServer()

	repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	_, err := server.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		UserId:      "user-1",
		AccountType: "checking",
	})

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	repo.AssertExpectations(t)
}

func TestGetAccount_gRPC_Success(t *testing.T) {
	server, repo := newTestServer()

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)

	resp, err := server.GetAccount(context.Background(), &pb.GetAccountRequest{AccountId: "acc-1"})

	require.NoError(t, err)
	assert.Equal(t, "acc-1", resp.Account.Id)
	assert.Equal(t, "1000", resp.Account.Balance)
	repo.AssertExpectations(t)
}

func TestGetAccount_gRPC_NotFound(t *testing.T) {
	server, repo := newTestServer()

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, nil)

	_, err := server.GetAccount(context.Background(), &pb.GetAccountRequest{AccountId: "acc-1"})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestGetAccount_gRPC_InternalError(t *testing.T) {
	server, repo := newTestServer()

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, errors.New("db error"))

	_, err := server.GetAccount(context.Background(), &pb.GetAccountRequest{AccountId: "acc-1"})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	repo.AssertExpectations(t)
}

func TestGetAccountByUserId_gRPC_Success(t *testing.T) {
	server, repo := newTestServer()

	accounts := []*domain.Account{
		newActiveAccount("acc-1", "user-1"),
		newActiveAccount("acc-2", "user-1"),
	}
	repo.On("GetByUserID", mock.Anything, "user-1").Return(accounts, nil)

	resp, err := server.GetAccountByUserId(context.Background(), &pb.GetAccountByUserIdRequest{UserId: "user-1"})

	require.NoError(t, err)
	assert.Len(t, resp.Accounts, 2)
	repo.AssertExpectations(t)
}

func TestGetAccountByNumber_gRPC_NotFound(t *testing.T) {
	server, repo := newTestServer()

	repo.On("GetByAccountNumber", mock.Anything, "0000000000").Return(nil, nil)

	_, err := server.GetAccountByNumber(context.Background(), &pb.GetAccountByNumberRequest{AccountNumber: "0000000000"})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestGetBalance_gRPC_Success(t *testing.T) {
	server, repo := newTestServer()

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)

	resp, err := server.GetBalance(context.Background(), &pb.GetBalanceRequest{AccountId: "acc-1"})

	require.NoError(t, err)
	assert.Equal(t, "acc-1", resp.AccountId)
	assert.Equal(t, "1000", resp.Balance)
	assert.Equal(t, "BRL", resp.Currency)
	repo.AssertExpectations(t)
}

func TestGetBalance_gRPC_NotFound(t *testing.T) {
	server, repo := newTestServer()

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, nil)

	_, err := server.GetBalance(context.Background(), &pb.GetBalanceRequest{AccountId: "acc-1"})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestDebitAccount_gRPC_Success(t *testing.T) {
	server, repo := newTestServer()

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("DebitAccount", mock.Anything, "acc-1", decimal.NewFromInt(100), "tx-1").
		Return(decimal.NewFromInt(900), nil)

	resp, err := server.DebitAccount(context.Background(), &pb.DebitAccountRequest{
		AccountId:     "acc-1",
		Amount:        "100",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "900", resp.NewBalance)
	repo.AssertExpectations(t)
}

func TestDebitAccount_gRPC_InvalidAmountFormat(t *testing.T) {
	server, _ := newTestServer()

	resp, err := server.DebitAccount(context.Background(), &pb.DebitAccountRequest{
		AccountId:     "acc-1",
		Amount:        "not-a-number",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "invalid amount format", resp.ErrorMessage)
}

func TestDebitAccount_gRPC_InsufficientFunds(t *testing.T) {
	server, repo := newTestServer()

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("DebitAccount", mock.Anything, "acc-1", mock.Anything, "tx-1").
		Return(decimal.Zero, service.ErrInsufficientFunds)

	resp, err := server.DebitAccount(context.Background(), &pb.DebitAccountRequest{
		AccountId:     "acc-1",
		Amount:        "5000",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "insufficient funds", resp.ErrorMessage)
	repo.AssertExpectations(t)
}

func TestDebitAccount_gRPC_AccountNotFound(t *testing.T) {
	server, repo := newTestServer()

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, nil)

	resp, err := server.DebitAccount(context.Background(), &pb.DebitAccountRequest{
		AccountId:     "acc-1",
		Amount:        "100",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "account not found", resp.ErrorMessage)
	repo.AssertExpectations(t)
}

func TestDebitAccount_gRPC_AccountBlocked(t *testing.T) {
	server, repo := newTestServer()

	account := newActiveAccount("acc-1", "user-1")
	account.Status = domain.AccountStatusBlocked
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)

	resp, err := server.DebitAccount(context.Background(), &pb.DebitAccountRequest{
		AccountId:     "acc-1",
		Amount:        "100",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "account is blocked", resp.ErrorMessage)
	repo.AssertExpectations(t)
}

func TestCreditAccount_gRPC_Success(t *testing.T) {
	server, repo := newTestServer()

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("CreditAccount", mock.Anything, "acc-1", decimal.NewFromInt(200), "tx-1").
		Return(decimal.NewFromInt(1200), nil)

	resp, err := server.CreditAccount(context.Background(), &pb.CreditAccountRequest{
		AccountId:     "acc-1",
		Amount:        "200",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "1200", resp.NewBalance)
	repo.AssertExpectations(t)
}

func TestCreditAccount_gRPC_InvalidAmount(t *testing.T) {
	server, _ := newTestServer()

	resp, err := server.CreditAccount(context.Background(), &pb.CreditAccountRequest{
		AccountId:     "acc-1",
		Amount:        "abc",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
}

func TestCreditAccount_gRPC_Failure(t *testing.T) {
	server, repo := newTestServer()

	account := newActiveAccount("acc-1", "user-1")
	account.Status = domain.AccountStatusBlocked
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)

	resp, err := server.CreditAccount(context.Background(), &pb.CreditAccountRequest{
		AccountId:     "acc-1",
		Amount:        "200",
		TransactionId: "tx-1",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	repo.AssertExpectations(t)
}
